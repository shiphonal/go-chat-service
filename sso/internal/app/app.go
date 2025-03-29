package app

import (
	authapp "ChatService/sso/internal/app/auth"
	grpcapp "ChatService/sso/internal/app/grpc"
	profileapp "ChatService/sso/internal/app/profile"
	"ChatService/sso/internal/services/auth"
	"ChatService/sso/internal/services/profile"
	"ChatService/sso/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
	AUTH       *auth.Auth
	PROFILE    *profile.Profile
}

func New(log *slog.Logger, port int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}
	authService := authapp.New(log, storage, storage, storage, tokenTTL)

	profileService := profileapp.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, profileService, port)

	return &App{
		GRPCServer: grpcApp,
		AUTH:       authService,
		PROFILE:    profileService,
	}
}
