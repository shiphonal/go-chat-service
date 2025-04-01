package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string        `json:"env" env-default:"local"`
	StoragePath string        `json:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `json:"token_ttl" env-required:"true"`
	GRPCServer  *GRPCConfig   `json:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	const op = "config.MustLoad"
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("CONFIG_PATH environment variable not set")
	}
	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		panic("failed to read config: " + err.Error())
	}
	return &config
}
