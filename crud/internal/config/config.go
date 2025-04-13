package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env           string         `yaml:"env" env-default:"local"`
	AppSecret     string         `yaml:"app_secret" env-required:"true" env:"APP_SECRET"`
	StoragePath   string         `yaml:"storage_path" env-required:"true" env:"STORAGE_PATH"`
	GRPCServer    *GRPCConfig    `yaml:"grpc"`
	ClientsConfig *ClientsConfig `yaml:"clients"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Client struct {
	Addr         string        `yaml:"addr"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retries_count"`
}

type ClientsConfig struct {
	SSO Client `yaml:"sso"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH_FOR_CRUD")
	if configPath == "" {
		panic("CONFIG_PATH environment variable not set")
	}
	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		panic("failed to read config: " + err.Error())
	}
	return &config
}
