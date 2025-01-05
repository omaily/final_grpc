package storage

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

func (db *Instance) ExchangeRates(ctx context.Context) ([]Exchange, error) {
	query := `select note, rate from cyrrency`
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		slog.Error("error while executing query:", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Exchange])
}

func (db *Instance) ExchangeCurrency(
	ctx context.Context, from_currency string, to_currency string,
) (float64, error) {
	var fromRate, toRate float64
	query := `select rate FROM cyrrency where note = $1`

	err := db.pool.QueryRow(ctx, query, from_currency).Scan(&fromRate)
	if err != nil {
		slog.Error("error while executing query:", slog.String("error", err.Error()))
		return 0, err
	}

	err = db.pool.QueryRow(ctx, query, to_currency).Scan(&toRate)
	if err != nil {
		slog.Error("error while executing query:", slog.String("error", err.Error()))
		return 0, err
	}

	return fromRate / toRate, nil
}
