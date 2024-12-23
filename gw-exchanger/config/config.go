package config

import (
	"log/slog"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-default:"local"`
	ServerGrpc `yaml:"server_grpc"`
	Storage    `yaml:"storage"`
}

type Storage struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	Role     string `yaml:"user" env-default:"postgres"`
	Pass     string `yaml:"pass" env-default:""`
	Database string `yaml:"database" env-default:"postgres"`
}

type ServerGrpc struct {
	Address string `yaml:"address" env-default:"localhost:8080"`
	Port    string `yaml:"port" env-default:":4000"`
}

var cfg Config
var once sync.Once

func MustLoad() *Config {
	once.Do(func() {
		slog.Info("read application configuration")
		cfg = Config{}
		if err := cleanenv.ReadConfig("config/config.yml", &cfg); err != nil {
			slog.Error("cannot read config: %s", slog.String("error", err.Error()))
		}
	})
	return &cfg
}
