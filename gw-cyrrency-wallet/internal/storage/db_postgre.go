package storage

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/config"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/pkg/model"
)

type Instance struct {
	conf *config.Storage
	pool *pgxpool.Pool
	User *model.Account
}

var pgOnce sync.Once

func New(conf config.Storage) *Instance {
	return &Instance{
		conf: &conf,
	}
}

func (c *Instance) Start(ctx context.Context) error {
	pgOnce.Do(func() {
		dbPath := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", c.conf.Role, c.conf.Pass, c.conf.Host, c.conf.Port, c.conf.Database)
		pool, err := pgxpool.New(ctx, dbPath)
		if err != nil {
			slog.Error("unable to create connection pool", slog.String("err", err.Error()))
			return
		}
		c.pool = pool
	})
	return nil
}

func (c *Instance) Stop() {
	c.pool.Close()
}
