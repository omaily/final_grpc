package server

import (
	"context"
	"net/http"
	"time"

	"github.com/omaily/final_grpc/gw-cyrrency-wallet/config"
	connGrpc "github.com/omaily/final_grpc/gw-cyrrency-wallet/connection/grpc"
	connRedis "github.com/omaily/final_grpc/gw-cyrrency-wallet/connection/redis"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/connection/storage"
)

type Http struct {
	conf        *config.HTTPServer
	conn        *http.Server
	storage     storage.Repository
	clientGrpc  connGrpc.Client
	clientRedis connRedis.Client
}

func New(conf config.HTTPServer, storage storage.Repository, grpc connGrpc.Client, redis connRedis.Client) *Http {
	return &Http{
		conf:        &conf,
		storage:     storage,
		clientRedis: redis,
		clientGrpc:  grpc,
	}
}

func (s *Http) Start(ctx context.Context) error {
	srv := &http.Server{
		Addr:         s.conf.Port,
		Handler:      s.router(),
		ReadTimeout:  s.conf.Timeout * time.Second,
		WriteTimeout: s.conf.Timeout * 2 * time.Second,
		IdleTimeout:  s.conf.IdleTimeout * time.Second,
	}

	okCh, errCh := make(chan struct{}), make(chan error)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		okCh <- struct{}{}
	}()

	select {
	case <-okCh:
		return nil
	case err := <-errCh:
		return err
	}
}

func (s *Http) Stop(ctx context.Context) error {
	okCh, errCh := make(chan struct{}), make(chan error)
	go func() {
		if err := s.conn.Shutdown(ctx); err != nil {
			errCh <- err
		}
		okCh <- struct{}{}
	}()

	select {
	case <-okCh:
		return nil
	case err := <-errCh:
		return err
	}
}
