package main

import (
	"ChatService/sso/internal/config"
	"ChatService/sso/internal/lib/logger"
)

func main() {
	// init config
	cfg := config.MustLoad()

	// init logger
	log := logger.SetupLogger(cfg.StoragePath)
	log.Info("Starting sso")

	// TODO: init app
}
