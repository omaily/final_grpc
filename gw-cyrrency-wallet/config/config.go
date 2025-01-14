package config

import (
	"log/slog"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"local"`
	HTTPServer  `yaml:"http_server"`
	Storage     `yaml:"storage"`
	GRPCServer  `yaml:"grpc_remote_server"`
	RedisServer `yaml:"redis_remote_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost"`
	Port        string        `yaml:"port" env-default:":8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type GRPCServer struct {
	Address string `yaml:"address" env-default:"localhost"`
	Port    int    `yaml:"port" env-default:":4000"`
}

type RedisServer struct {
	Address string `yaml:"address" env-default:"localhost"`
	Port    int    `yaml:"port" env-default:":6379"`
}

type Storage struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	Role     string `yaml:"user" env-default:"postgres"`
	Pass     string `yaml:"pass" env-default:""`
	Database string `yaml:"database" env-default:"postgres"`
}

var cfg Config
var once sync.Once

func MustLoad() *Config {
	once.Do(func() {
		slog.Info("read application configuration")
		cfg = Config{}
		if err := cleanenv.ReadConfig("config/config.yml", &cfg); err != nil {
			slog.Error("cannot read config: ", slog.String("error", err.Error()))
		}
	})
	return &cfg
}
