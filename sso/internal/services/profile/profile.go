package profile

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

type Profile struct {
	Log          *slog.Logger
	UserRefactor UserRefactor
	UserAdmin    UserAdmin
	TokenTTL     time.Duration
}

type UserRefactor interface {
	GetUserById(ctx context.Context, id int64) (models.User, error)
	UpdatePassword(ctx context.Context, password []byte, id int64) (bool, error)
	UpdateName(ctx context.Context, id int64, newName string) (bool, error)
}

type UserAdmin interface {
	UpdateRole(ctx context.Context, id int64, newRole int32) (bool, error)
}

var (
	ErrUserNotAdmin       = errors.New("user is not admin")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (p *Profile) ChangePassword(ctx context.Context, oldPassword string, password string, id int64) (bool, error) {
	const op = "services.profile.ChangePassword"
	p.Log.With(slog.String("op", op))

	user, err := p.UserRefactor.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			p.Log.Warn("user not found")
			return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		p.Log.Warn("failed to get user")
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(oldPassword)); err != nil {
		p.Log.Warn("invalid password")
		return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		p.Log.Error("failed to generate password hash")
		return false, fmt.Errorf("%s: %w", op, err)
	}

	success, err := p.UserRefactor.UpdatePassword(ctx, passHash, id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			p.Log.Warn("User not found")
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		p.Log.Error("Failed with refactor user")
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return success, nil
}

func (p *Profile) ChangeName(ctx context.Context, id int64, newName string) (bool, error) {
	const op = "services.profile.ChangeName"
	p.Log.With(slog.String("op", op))

	success, err := p.UserRefactor.UpdateName(ctx, id, newName)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			p.Log.Warn("User not found")
			return false, fmt.Errorf("%s, %w", op, storage.ErrUserNotFound)
		}
		p.Log.Error("Failed with refactor user")
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return success, nil
}

func (p *Profile) ChangeRole(ctx context.Context, password string, idAdmin int64, id int64, newRole int32) (bool, error) {
	const op = "services.profile.ChangeRole"
	p.Log.With(slog.String("op", op))

	user, err := p.UserRefactor.GetUserById(ctx, idAdmin)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			p.Log.Warn("user not found")
			return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		p.Log.Warn("failed to get user")
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		p.Log.Warn("invalid password")
		return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	success, err := p.UserAdmin.UpdateRole(ctx, id, newRole)
	if err != nil {
		p.Log.Error("Failed with refactor user")
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return success, nil
}
