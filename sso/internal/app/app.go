package app

import grpcapp "ChatService/sso/internal/app/grpc"

type App struct {
	gRPCServer *grpcapp.App
}

func New() *App {
	// TODO: init storage
	// TODO: init auth service
	// TODO: init profile service
	// TODO init grpc
	return &App{}
}
