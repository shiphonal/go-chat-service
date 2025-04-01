package app

import (
	grpcApp "ChatService/crud/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcApp.App
}

func New(log *slog.Logger, tokenTTL time.Duration, storagePath string, port int) *App {
	// TODO: init crud service
	grpcSever := grpcApp.New(log, port)
	return &App{
		GRPCServer: grpcSever,
	}
}
