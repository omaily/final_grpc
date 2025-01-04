package main

import (
	"context"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/config"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/app"
)

func init() {
	err := godotenv.Load("wallet.env")
	if err != nil {
		slog.Error("Error loading .env file:" + err.Error())
	}
}

func main() {
	conf := config.MustLoad()

	app, err := app.New(context.Background(), conf)
	if err != nil {
		slog.Error("could not initialize server: " + err.Error())
		return
	}

	if err := app.Run(); err != nil {
		slog.Error(err.Error())
	}
}
