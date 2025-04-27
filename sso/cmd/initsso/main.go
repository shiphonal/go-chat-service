package main

import (
	"ChatService/sso/internal/app"
	client "ChatService/sso/internal/clients"
	"ChatService/sso/internal/config"
	"ChatService/sso/internal/lib/logger"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// init config
	cfg := config.MustLoad()

	// init logger
	log := logger.SetupLogger(cfg.Env)
	log.Info("Starting sso")

	clientFabric := client.ClientMustLoad(cfg, log)

	go func() {
		log.Info("Starting HTTP server on :8080")
		if err := clientFabric.HttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server error", err)
		}
	}()

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Ожидаем сигнала для завершения работы
	<-stop

	application.GRPCServer.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := clientFabric.HttpServer.Shutdown(ctx); err != nil {
		log.Error("HTTP shutdown error", err)
	}

	application.GRPCServer.Stop()
	log.Info("Stopping sso")
}
