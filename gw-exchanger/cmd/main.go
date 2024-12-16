package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/omaily/final_grpc/gw-exchanger/config"
	"github.com/omaily/final_grpc/gw-exchanger/internal/app"
)

func main() {
	fmt.Println("hello my friend")

	conf := config.MustLoad()

	app, err := app.New(context.Background(), conf)
	if err != nil {
		slog.Error("could not initialize server: %w", err)
		return
	}

	app.Run()
}
