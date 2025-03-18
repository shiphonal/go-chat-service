package main

import (
	"ChatService/sso/internal/config"
	"ChatService/sso/internal/lib/logger"
	"fmt"
)

func main() {
	// init config
	cfg := config.MustLoad()
	fmt.Printf("%#v\n", cfg)

	// init logger
	log := logger.SetupLogger(cfg.Env)
	log.Info("Starting sso")

	// TODO: init app
}
