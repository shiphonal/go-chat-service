package main

import (
	"ChatService/sso/internal/app"
	"ChatService/sso/internal/config"
	"ChatService/sso/internal/lib/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// init config
	cfg := config.MustLoad()

	// init logger
	log := logger.SetupLogger(cfg.Env)
	log.Info("Starting sso")

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Ожидаем сигнала для завершения работы
	<-stop

	application.GRPCServer.Stop()
	log.Info("Stopping sso")
}
