package jwt

import (
	"ChatService/sso/internal/domain/models"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = user.ID
	claims["email"] = user.Email
	claims["appID"] = app.ID
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

type TokenInfo struct {
	Error  error
	UserID int64
}

var (
	ErrTokenExpired = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
)

func ValidateToken(tokenString string, secret string) TokenInfo {
	// Decoding jwt
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return TokenInfo{Error: err}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Извлечение времени истечения
		exp, ok := claims["exp"].(float64)
		if !ok {
			return TokenInfo{Error: ErrInvalidToken}
		}

		// Преобразование времени истечения в time.Time
		expTime := time.Unix(int64(exp), 0)

		// Проверка, истек ли токен
		if time.Now().After(expTime) {
			return TokenInfo{Error: ErrTokenExpired}
		}

		userId, ok := claims["userID"]
		if !ok {
			return TokenInfo{Error: ErrInvalidToken}
		}
		return TokenInfo{Error: nil, UserID: int64(userId.(float64))}
	}
	return TokenInfo{Error: ErrInvalidToken}
}
