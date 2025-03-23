package app

import (
	grpcapp "ChatService/sso/internal/app/grpc"
	"ChatService/sso/internal/services/auth"
	"ChatService/sso/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, port int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	// TODO: init profile service

	grpcApp := grpcapp.New(log, authService, port)

	return &App{
		GRPCServer: grpcApp,
	}
}
