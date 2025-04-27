package main

import (
	"ChatService/crud/internal/app"
	client "ChatService/crud/internal/clients/app"
	"ChatService/crud/internal/config"

	Logger "ChatService/crud/internal/lib/logger"
	"context"
	"errors"

	"net/http"
	"os"
	"os/signal"

	"syscall"
	"time"
)

func main() {
	cnf := config.MustLoad()

	logger := Logger.SetupLogger(cnf.Env)
	logger.Info("Starting logger")

	clientFabric := client.ClientMustLoad(cnf, logger)

	go func() {
		logger.Info("Starting HTTP server on :8080")
		if err := clientFabric.HttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server error", err)
		}
	}()

	application := app.New(logger, cnf.Storage.Path, cnf.AppSecret, cnf.GRPC.Server.Port, clientFabric.CRUD)
	logger.Info("Starting application")

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := clientFabric.HttpServer.Shutdown(ctx); err != nil {
		logger.Error("HTTP shutdown error", err)
	}

	logger.Info("Gracefully stopped")
}
