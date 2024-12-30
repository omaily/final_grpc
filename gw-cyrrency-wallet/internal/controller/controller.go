package controller

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/omaily/final_grpc/gw-cyrrency-wallet/config"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/connector"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/storage"
)

type Http struct {
	conf       *config.HTTPServer
	storage    *storage.Instance //убрать
	clientGrpc *connector.GrpcClient
}

func New(conf config.HTTPServer, storage *storage.Instance, grpc *connector.GrpcClient) *Http {
	return &Http{
		conf:       &conf,
		storage:    storage,
		clientGrpc: grpc,
	}
}

func (s *Http) Start(cxt context.Context) error {

	srv := &http.Server{
		Addr:         s.conf.Port,
		Handler:      s.router(),
		ReadTimeout:  s.conf.Timeout * time.Second,
		WriteTimeout: s.conf.Timeout * 2 * time.Second,
		IdleTimeout:  s.conf.IdleTimeout * time.Second,
	}

	go func() {
		slog.Info("starting server to", slog.String("addres", s.conf.Address))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("не маслает: ", slog.String("err", err.Error()))
		}
	}()

	return nil
}

func (s *Http) Stop() {
	slog.Info("...down http_server")
}
