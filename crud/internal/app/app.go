package app

import (
	crudApp "ChatService/crud/internal/app/crud"
	grpcApp "ChatService/crud/internal/app/grpc"
	"ChatService/crud/internal/storage/postgres"
	"log/slog"
)

type App struct {
	GRPCServer *grpcApp.App
}

func New(log *slog.Logger, port int) *App {

	storagePostgres, err := postgres.New()
	if err != nil {
		panic(err)
	}
	crudService := crudApp.New(log, storagePostgres)
	grpcSever := grpcApp.New(log, crudService, port)
	return &App{
		GRPCServer: grpcSever,
	}
}
