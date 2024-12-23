package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

func (s *Instance) GetAmount(ctx context.Context, uuid string) (int, error) {
	query := `select amount from account where uuid = $1`
	var cash int
	row := s.pool.QueryRow(ctx, query, uuid)
	err := row.Scan(&cash)
	if err != nil {
		slog.Error("Error Fetching Book Details: %w", slog.String("err", err.Error()))
		return 0, err
	}

	return cash, nil
}

func (s *Instance) DepositPay(ctx context.Context, uuid string, amount int) error {
	cash, err := s.GetAmount(ctx, uuid)
	if err != nil {
		return err
	}

	cash += amount

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	args := pgx.NamedArgs{
		"cash": cash,
		"uuid": uuid,
	}

	_, err = tx.Exec(ctx, "UPDATE account SET amount = @cash where uuid = @uuid", args)
	if err != nil {
		slog.Error("UPDATE fall: %w", slog.String("err", err.Error()))
		return err
	}

	slog.Info("user depositPay", slog.String("user", uuid))
	return nil
}

func (s *Instance) WithdrawPay(ctx context.Context, uuid string, amount int) error {

	cash, err := s.GetAmount(ctx, uuid)
	if err != nil {
		return err
	}

	if amount > cash {
		return errors.New("there are insufficient funds in your account")
	}

	cash -= amount

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	args := pgx.NamedArgs{
		"cash": cash,
		"uuid": uuid,
	}

	_, err = tx.Exec(ctx, "UPDATE account SET amount = @cash where uuid = @uuid", args)
	if err != nil {
		slog.Error("UPDATE fall: %w", slog.String("err", err.Error()))
		return err
	}

	slog.Info("user withdrawPay", slog.String("user", uuid))
	return nil
}

func (ex *Exchange) Exchanges(ctx context.Context) (map[string]float64, error) {
	slog.Info("metod Exchanges")
	rate := map[string]float64{}
	return rate, nil
}

func (ex *Exchange) Exchange(
	ctx context.Context, from_currency string, to_currency string,
) (int, error) {
	slog.Info("metod Exchange")
	return 11, fmt.Errorf("empty")
}
