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
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

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

func New(log *slog.Logger,
	userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email, password string, appID int) (string, error) {
	const op = "auth.Login"
	// Get User
	user, err := a.userProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found")
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		a.log.Warn("failed to get user")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Valid Password
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Warn("invalid password")
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	// Get App
	app, err := a.appProvider.GetApp(ctx, appID)
	if err != nil {
		a.log.Warn(err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}
	a.log.Debug("user logged in successfully")

	// Create Token
	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to create token")
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}

func (a *Auth) Logout(ctx context.Context, token string, userID int64) (bool, error) {
	const op = "auth.Logout"
	log := a.log.With(slog.String("op", op))
	log.Debug("start")
	return true, nil
}

func (a *Auth) Register(ctx context.Context, username, email, password string) (int64, error) {
	const op = "auth.Register"
	a.log.With(slog.String("op", op))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error("failed to generate password hash")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, username, email, passHash)
	if err != nil {
		a.log.Error("failed to save user")
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"
	a.log.With(slog.String("op", op))

	success, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		a.log.Warn("failed to get user")
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return success, nil
}

func (a *Auth) IsModerator(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsModerator"
	a.log.With(slog.String("op", op))
	success, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		a.log.Warn("failed to get user")
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return success, nil
}
