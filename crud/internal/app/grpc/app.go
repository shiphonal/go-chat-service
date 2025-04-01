package grpc

import (
	"google.golang.org/grpc"
	"log/slog"
)

type App struct {
	logger     *slog.Logger
	port       int
	GrpcServer *grpc.Server
}

func New(log *slog.Logger /* //TODO: init crud service */, port int) *App {
	gRPCServer := grpc.NewServer()
	return &App{
		logger:     log,
		port:       port,
		GrpcServer: gRPCServer,
	}
}
