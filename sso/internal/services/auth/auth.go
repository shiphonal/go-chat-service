package auth

import (
	"ChatService/sso/internal/domain/models"
	"ChatService/sso/internal/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log          slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, username, email string, passHash []byte) (int64, error)
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (models.User, error)
	GetApp(ctx context.Context, id int) (models.App, error)
}

type AppProvider interface {
	GetApp(ctx context.Context, id int) (models.App, error)
}

func New(log slog.Logger,
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
	a.log.Debug("start")
	user, err := a.userProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found")
			return "", fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		a.log.Warn("failed to get user")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Warn("invalid password")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	app, err := a.appProvider.GetApp(ctx, appID)
	if err != nil {
		a.log.Warn("failed to get app")
		return "", fmt.Errorf("%s: %w", op, err)
	}
	a.log.Debug("user logged in successfully")

	return app.Secret, nil
	// TODO: get token jwt and return it
}

func (a *Auth) Logout(ctx context.Context) {
	const op = "auth.Logout"
	log := a.log.With(slog.String("op", op))
	log.Debug("start")
}

func (a *Auth) Register(ctx context.Context, username, email, password string) (int64, error) {
	const op = "auth.Register"
	a.log.With(slog.String("op", op))
	a.log.Debug("start")

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
