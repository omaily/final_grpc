package storage

import (
	"context"

	"github.com/omaily/final_grpc/gw-cyrrency-wallet/pkg/model"
)

type Repository interface {
	Start(ctx context.Context) error
	Stop()

	CreateAccount(ctx context.Context, user *model.Account) (userId [16]byte, err error)
	FindAccount(ctx context.Context, user *model.Account) (userId [16]byte, pass string, err error)

	CheckBalance(ctx context.Context, strUUID string) (map[string]float64, error)
	PutMoney(ctx context.Context, strUUID string, deposit model.Transfer) (map[string]float64, error)
	TakeMoney(ctx context.Context, strUUID string, deposit model.Transfer) (map[string]float64, error)
	ChangeMoney(ctx context.Context, strUUID string, ex model.Exchange, odds float64) error
}
