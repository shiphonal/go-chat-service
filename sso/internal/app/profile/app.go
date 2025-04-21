package profileApp

import (
	"ChatService/sso/internal/services/profile"
	"log/slog"
	"time"
)

func New(log *slog.Logger,
	userRefactor profile.UserRefactor, userAdmin profile.UserAdmin, tokenTTL time.Duration) *profile.Profile {
	return &profile.Profile{
		Log:          log,
		UserRefactor: userRefactor,
		UserAdmin:    userAdmin,
		TokenTTL:     tokenTTL,
	}
}
