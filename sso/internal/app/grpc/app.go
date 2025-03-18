package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int) *App {
	// TODO: register auth and profile services
	return &App{
		log:        log,
		gRPCServer: grpc.NewServer(),
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "sso.app.Run"

	log := a.log.With(slog.String("op", op))
	// Прослушиваем порт по tcp
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("Starting gRPC server")
	// Создаём gRPS сервер
	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "sso.app.Stop"
	log := a.log.With(slog.String("op", op))
	a.gRPCServer.GracefulStop()
	log.Info("Stopped gRPC server")
}
