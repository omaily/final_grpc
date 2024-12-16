package storage

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/omaily/final_grpc/gw-exchanger/config"
)

type Connector struct {
	conf   *config.Storage
	dbPath string
	pool   *pgxpool.Pool
}

var pgOnce sync.Once

func New(conf config.Storage) *Connector {
	return &Connector{
		dbPath: fmt.Sprintf("postgres://%s:%s@%s:%d/%s", conf.Role, conf.Pass, conf.Host, conf.Port, conf.Database),
		conf:   &conf,
	}
}

func (c *Connector) Start(ctx context.Context) error {
	pgOnce.Do(func() {
		pool, err := pgxpool.New(ctx, c.dbPath)
		if err != nil {
			slog.Error("unable to create connection pool", slog.String("err", err.Error()))
			return
		}
		c.pool = pool
	})
	return nil
}

func (c *Connector) Stop(_ context.Context) error {
	c.pool.Close()
	return nil
}
