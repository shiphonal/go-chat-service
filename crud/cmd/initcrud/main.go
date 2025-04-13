package main

import (
	"ChatService/crud/internal/app"
	"ChatService/crud/internal/config"
	Logger "ChatService/crud/internal/lib/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cnf := config.MustLoad()

	logger := Logger.SetupLogger(cnf.Env)
	logger.Info("Starting logger")

	application := app.New(logger, cnf.StoragePath, cnf.GRPCServer.Port)
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
