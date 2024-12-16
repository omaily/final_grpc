package controller

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/omaily/final_grpc/gw-exchanger/config"
	h "github.com/omaily/final_grpc/gw-exchanger/internal/delivery"
	"github.com/omaily/final_grpc/gw-exchanger/internal/storage"
)

type Http struct {
	conf    *config.HTTPServer
	storage *storage.Connector
}

func New(conf config.HTTPServer, instance *storage.Connector) *Http {
	return &Http{
		conf:    &conf,
		storage: instance,
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

func (s *Http) router() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("POST /api/v1/wallet", h.PostWallet(s.storage))
	router.HandleFunc("GET /api/v1/wallets/{uuid}", h.GetWallet(s.storage))

	return router
}