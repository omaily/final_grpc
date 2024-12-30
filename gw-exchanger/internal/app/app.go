package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/omaily/final_grpc/gw-exchanger/config"
	"github.com/omaily/final_grpc/gw-exchanger/internal/connector"
	"github.com/omaily/final_grpc/gw-exchanger/internal/storage"
)

type App struct {
	conf       *config.Config
	storage    *storage.Instance
	serverGrpc *connector.ServerGrpc
}

func New(ctx context.Context, conf *config.Config) (*App, error) {
	if conf == nil {
		return nil, errors.New("configuration files are not initialized")
	}

	return &App{
		conf: conf,
	}, nil
}

func (a *App) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.start(ctx); err != nil {
		slog.Error("could not initialize server: %s", slog.Any("error", err))
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	slog.Info("stopping server due to syscall or collapse")
	ctx, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	return a.stop(ctx)
}

func (a *App) start(ctx context.Context) error {
	storage := storage.New(a.conf.Storage)
	if err := storage.Start(ctx); err != nil {
		return fmt.Errorf("could not initialize storage: %s", err)
	}
	a.storage = storage
	a.serverGrpc = connector.New(storage)

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen on port 8081: %v", err)
	}
	log.Printf("gRPC-сервер прослушивает %v", lis.Addr())

	if err := a.serverGrpc.Serve(lis); err != nil {
		log.Fatalf("не удалось обслужить: %v", err)
	}

	return nil
}

func (a *App) stop(_ context.Context) error {
	slog.Info("process shutting down service...")
	return nil
}
