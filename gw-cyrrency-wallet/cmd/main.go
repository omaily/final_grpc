package main

import (
	"context"
	"log/slog"

	"github.com/joho/godotenv"

	"github.com/omaily/final_grpc/gw-cyrrency-wallet/config"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/app"
)

func init() {
	err := godotenv.Load("file.env")
	if err != nil {
		slog.Error("Error loading .env file: %s", slog.String("error", err.Error()))
	}
}

func main() {
	conf := config.MustLoad()

	app, err := app.New(context.Background(), conf)
	if err != nil {
		slog.Error("could not initialize server: %w", slog.String("error", err.Error()))
		return
	}

	app.Run()
}
