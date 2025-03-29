package authApp

import (
	"ChatService/sso/internal/services/auth"
	"log/slog"
	"time"
)

func New(log *slog.Logger,
	userSaver auth.UserSaver, userProvider auth.UserProvider, appProvider auth.AppProvider, tokenTTL time.Duration) *auth.Auth {
	return &auth.Auth{
		Log:          log,
		UserSaver:    userSaver,
		UserProvider: userProvider,
		AppProvider:  appProvider,
		TokenTTL:     tokenTTL,
	}
}
