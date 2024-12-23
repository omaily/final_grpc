package storage

import (
	"context"
	"fmt"
	"log/slog"
)

type Exchange struct {
	currency string
	rate     int
}

func (ex *Exchange) GetRates(ctx context.Context) (map[string]float64, error) {
	slog.Info("metod Exchanges")
	rate := map[string]float64{}
	return rate, nil
}

func (ex *Exchange) GetRate(
	ctx context.Context, from_currency string, to_currency string,
) (int, error) {
	slog.Info("metod Exchange")
	return 11, fmt.Errorf("empty")
}
