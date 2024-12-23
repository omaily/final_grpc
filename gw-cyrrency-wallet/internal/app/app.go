package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/omaily/final_grpc/gw-cyrrency-wallet/config"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/controller"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/storage"
)

type App struct {
	conf       *config.Config
	storage    *storage.Instance
	serverHttp *controller.Http
}

func New(ctx context.Context, conf *config.Config) (*App, error) {
	if conf == nil {
		return nil, errors.New("configuration files are not initialized")
	}

	http := &conf.HTTPServer
	if http.Address == "" || http.Port == "" {
		return nil, errors.New("configuration address cannot be blank")
	}

	return &App{
		conf:    conf,
		storage: storage.NewConnector(),
	}, nil
}

func (a *App) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.start(ctx); err != nil {
		slog.Error("could not initialize server: %s", slog.String("error", err.Error()))
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
	//start http service
	listener := controller.New(a.conf.HTTPServer)
	if err := listener.Start(ctx); err != nil {
		return fmt.Errorf("could not initialize controller: %s", err)
	}
	a.serverHttp = listener

	return nil
}

func (a *App) stop(_ context.Context) error {
	slog.Info("process shutting down service...")
	a.serverHttp.Stop()

	return nil
}
