package config

import (
	"log/slog"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
	Storage     `yaml:"storage"`
}

type Storage struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	Role     string `yaml:"user" env-default:"postgres"`
	Pass     string `yaml:"pass" env-default:""`
	Database string `yaml:"database" env-default:"postgres"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Port        string        `yaml:"port" env-default:":4000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

var cfg Config
var once sync.Once

func MustLoad() *Config {
	once.Do(func() {
		slog.Info("read application configuration")
		cfg = Config{}
		if err := cleanenv.ReadConfig("config/config.yml", &cfg); err != nil {
			slog.Error("cannot read config: %s", err)
		}
	})
	return &cfg
}
