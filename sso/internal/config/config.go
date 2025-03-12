package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        *GRPCConfig   `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	// if config path is empty
	if configPath == "" {
		panic("CONFIG_PATH environment variable not set")
	}
	// if config exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("CONFIG_PATH does not exist")
	}

	var cnf Config
	if err := cleanenv.ReadConfig(configPath, &cnf); err != nil {
		panic("Failed with reading config" + err.Error())
	}
	return &cnf
}
