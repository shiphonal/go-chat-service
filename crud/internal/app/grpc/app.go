package grpc

import (
	"ChatService/crud/internal/grpc/crud"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	logger     *slog.Logger
	port       int
	GrpcServer *grpc.Server
}

func New(log *slog.Logger, crudService crud.CRUD, port int) *App {
	gRPCServer := grpc.NewServer()
	crud.RegisterServer(gRPCServer, crudService)
	return &App{
		logger:     log,
		port:       port,
		GrpcServer: gRPCServer,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpc.App.Run"
	a.logger.With(slog.String("op", op))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if err = a.GrpcServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	a.logger.Info("grpc server started", slog.String("address", lis.Addr().String()))
	return nil
}

func (a *App) Stop() {
	const op = "grpc.App.Stop"
	a.logger.With(slog.String("op", op)).
		Info("stopping grpc server", slog.Int("port", a.port))
	a.GrpcServer.GracefulStop()
}
