package storage

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

func (db *Instance) ExchangeRates(ctx context.Context) ([]Exchange, error) {
	slog.Info("storage no cursor, method bulk select")

	query := `select note, rate from cyrrency`
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		slog.Error("error while executing query:", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Exchange])
}

func (db *Instance) ExchangeRate(
	ctx context.Context, from_currency string, to_currency string,
) (float64, error) {
	slog.Info("storage no cursor, method single select")

	var fromRate, toRate float64
	query := `select note, rate FROM cyrrency where note = $1`

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

// func (s *Instance) GetAmount(ctx context.Context, uuid string) (int, error) {
// 	query := `select amount from account where uuid = $1`
// 	var cash int
// 	row := s.pool.QueryRow(ctx, query, uuid)
// 	err := row.Scan(&cash)
// 	if err != nil {
// 		slog.Error("Error Fetching Book Details: %w", slog.String("err", err.Error()))
// 		return 0, err
// 	}

// 	return cash, nil
// }

// func (s *Instance) DepositPay(ctx context.Context, uuid string, amount int) error {
// 	cash, err := s.GetAmount(ctx, uuid)
// 	if err != nil {
// 		return err
// 	}

// 	cash += amount

// 	tx, err := s.pool.Begin(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if err != nil {
// 			tx.Rollback(ctx)
// 		} else {
// 			tx.Commit(ctx)
// 		}
// 	}()

// 	args := pgx.NamedArgs{
// 		"cash": cash,
// 		"uuid": uuid,
// 	}

// 	_, err = tx.Exec(ctx, "UPDATE account SET amount = @cash where uuid = @uuid", args)
// 	if err != nil {
// 		slog.Error("UPDATE fall: %w", slog.String("err", err.Error()))
// 		return err
// 	}

// 	slog.Info("user depositPay", slog.String("user", uuid))
// 	return nil
// }

// func (s *Instance) WithdrawPay(ctx context.Context, uuid string, amount int) error {

// 	cash, err := s.GetAmount(ctx, uuid)
// 	if err != nil {
// 		return err
// 	}

// 	if amount > cash {
// 		return errors.New("there are insufficient funds in your account")
// 	}

// 	cash -= amount

// 	tx, err := s.pool.Begin(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if err != nil {
// 			tx.Rollback(ctx)
// 		} else {
// 			tx.Commit(ctx)
// 		}
// 	}()

// 	args := pgx.NamedArgs{
// 		"cash": cash,
// 		"uuid": uuid,
// 	}

// 	_, err = tx.Exec(ctx, "UPDATE account SET amount = @cash where uuid = @uuid", args)
// 	if err != nil {
// 		slog.Error("UPDATE fall: %w", slog.String("err", err.Error()))
// 		return err
// 	}

// 	slog.Info("user withdrawPay", slog.String("user", uuid))
// 	return nil
// }
