package storage

import (
	"context"
	"fmt"
	"log/slog"
)

type Account struct {
	username string
	password string
	email    string
}

func (a *Account) String() string {
	out := fmt.Sprintf("username = %v\t", a.username)
	out += fmt.Sprintf("password = %v\n", a.password)
	return out
}

type Exchange struct {
	currency string
	rate     int
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
