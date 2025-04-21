package main

import (
	"ChatService/crud/internal/app"
	"ChatService/crud/internal/clients/sso"
	"ChatService/crud/internal/config"
	Logger "ChatService/crud/internal/lib/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cnf := config.MustLoad()

	logger := Logger.SetupLogger(cnf.Env)
	logger.Info("Starting logger")

	ssoClient, err := sso.New(
		context.Background(),
		logger,
		cnf.ClientsConfig.SSO.Addr,
		cnf.ClientsConfig.SSO.Timeout,
		cnf.ClientsConfig.SSO.RetriesCount,
	)
	if err != nil {
		logger.Error("failed to initialize SSO client", err)
		os.Exit(1)
	}

	application := app.New(logger, cnf.StoragePath, cnf.AppSecret, cnf.GRPCServer.Port, ssoClient)
	logger.Info("Starting application")

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	logger.Info("Gracefully stopped")
}
