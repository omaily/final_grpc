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
	connGrpc "github.com/omaily/final_grpc/gw-cyrrency-wallet/connection/grpc"
	connRedis "github.com/omaily/final_grpc/gw-cyrrency-wallet/connection/redis"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/server"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/storage"
)

type appContext struct {
	conf     *config.Config
	instance *server.Http
	cmps     []remoteServers
}

type remoteServers struct {
	Name    string
	Service ClientApplications
}

type ClientApplications interface {
	Start(ctx context.Context) error
	Stop()
}

func New(ctx context.Context, conf *config.Config) (*appContext, error) {
	if conf == nil {
		return nil, errors.New("configuration files are not initialized")
	}

	http := &conf.HTTPServer
	if http.Address == "" || http.Port == "" {
		return nil, errors.New("configuration address http_server cannot be blank")
	}

	grpc := &conf.GRPCServer
	if grpc.Address == "" {
		return nil, errors.New("configuration address grpc_server cannot be blank")
	}

	return &appContext{
		conf: conf,
	}, nil
}

func (a *appContext) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.start(ctx); err != nil {
		return err
	}

	chShutdown := make(chan os.Signal, 1)
	signal.Notify(chShutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-chShutdown

	slog.Info("stopping server due to syscall or collapse")
	ctx, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	return a.stop(ctx)
}

func (a *appContext) start(ctx context.Context) error {
	storage := storage.New(a.conf.Storage)
	clientGrpc := connGrpc.New(a.conf.GRPCServer)
	clientRedis := connRedis.New(a.conf.RedisServer)

	a.cmps = []remoteServers{
		{"db-postgre", storage},
		{"GRPC server", clientGrpc},
		{"Redis server", clientGrpc},
	}

	for _, cmp := range a.cmps {
		if err := cmp.Service.Start(ctx); err != nil {
			return fmt.Errorf("error connecting to %s: %w", cmp.Name, err)
		}
		slog.Info(fmt.Sprintf("%v started", cmp.Name))
	}

	//start http service
	srv := server.New(a.conf.HTTPServer, storage, clientGrpc, clientRedis)
	if err := srv.Start(ctx); err != nil {
		return fmt.Errorf("could not initialize controller: %s", err)
	}
	a.instance = srv

	return nil
}

func (a *appContext) stop(ctx context.Context) error {
	slog.Info("process shutting down server... ")
	a.instance.Stop(ctx)

	for _, client := range a.cmps {
		slog.Info("process shutting down " + client.Name + "...")
		client.Service.Stop()
	}

	return nil
}
