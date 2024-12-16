package storage

import (
	"context"
)

type Instance interface {
	GetAmount(ctx context.Context, uuid string) (int, error)
	DepositPay(ctx context.Context, uuid string, amount int) error
	WithdrawPay(ctx context.Context, uuid string, amount int) error

	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
