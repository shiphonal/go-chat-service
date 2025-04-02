package jwt

import (
	"ChatService/sso/internal/domain/models"
	"context"
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

func ValidateToken(ctx context.Context, tokenString string, app models.App) (bool, error) {
	// Decoding jwt
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, fmt.Errorf("unexpected signing method")
		}
		return []byte(app.Secret), nil
	})
	if err != nil || !token.Valid {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Извлечение времени истечения
		exp, ok := claims["exp"].(float64)
		if !ok {
			return false, fmt.Errorf("invalid expiration time")
		}

		// Преобразование времени истечения в time.Time
		expTime := time.Unix(int64(exp), 0)

		// Проверка, истек ли токен
		if time.Now().After(expTime) {
			return false, fmt.Errorf("token has expired")
		}

		// Проверка контекста на истечение времени
		select {
		case <-ctx.Done():
			return false, fmt.Errorf("context deadline exceeded")
		default:
			// Контекст еще не истек
		}
		return true, nil
	}
	return false, fmt.Errorf("invalid token")
}
