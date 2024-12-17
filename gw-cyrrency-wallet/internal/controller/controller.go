package controller

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/omaily/final_grpc/gw-cyrrency-wallet/config"
)

type Http struct {
	conf *config.HTTPServer
}

func New(conf config.HTTPServer) *Http {
	return &Http{
		conf: &conf,
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
