package auth

import (
	"ChatService/sso/internal/domain/models"
	"ChatService/sso/internal/lib/jwt"
	"ChatService/sso/internal/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	Log          *slog.Logger
	UserSaver    UserSaver
	UserProvider UserProvider
	AppProvider  AppProvider
	TokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, username, email string, passHash []byte) (int64, error)
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, id int64) (bool, error)
	IsModerator(ctx context.Context, id int64) (bool, error)
}

type AppProvider interface {
	GetApp(ctx context.Context, id int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (a *Auth) Login(ctx context.Context, email, password string, appID int) (string, error) {
	const op = "auth.Login"
	a.Log.With(slog.String("op", op))
	// Get User
	user, err := a.UserProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.Log.Warn("user not found")
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		a.Log.Warn("failed to get user")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Valid Password
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.Log.Warn("invalid password")
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	// Get App
	app, err := a.AppProvider.GetApp(ctx, appID)
	if err != nil {
		a.Log.Warn(err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}
	a.Log.Debug("user logged in successfully")

	// Create Token
	token, err := jwt.NewToken(user, app, a.TokenTTL)
	if err != nil {
		a.Log.Error("failed to create token")
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}

func (a *Auth) Logout(ctx context.Context, token string) (bool, error) {
	const op = "auth.Logout"
	a.Log.With(slog.String("op", op))
	a.Log.Info("user logged out successfully")
	return true, nil
}

func (a *Auth) Register(ctx context.Context, username, email, password string) (int64, error) {
	const op = "auth.Register"
	a.Log.With(slog.String("op", op))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.Log.Error("failed to generate password hash")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.UserSaver.SaveUser(ctx, username, email, passHash)
	if err != nil {
		a.Log.Error("failed to save user")
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	a.Log.Debug("user register in successfully")
	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"
	a.Log.With(slog.String("op", op))

	success, err := a.UserProvider.IsAdmin(ctx, userID)
	if err != nil {
		a.Log.Warn("failed to get user")
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return success, nil
}

func (a *Auth) IsModerator(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsModerator"
	a.Log.With(slog.String("op", op))
	success, err := a.UserProvider.IsModerator(ctx, userID)
	if err != nil {
		a.Log.Warn("failed to get user")
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return success, nil
}
