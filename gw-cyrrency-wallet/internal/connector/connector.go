package connector

import (
	"context"

	exchange "github.com/omaily/final_grpc/gw-cyrrency-wallet/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IServiceExchange interface {
	Exchanges(
		ctx context.Context,
	) (rates map[string]float64, err error)

	Exchange(
		ctx context.Context,
		from_currency string,
		to_currency string,
	) (rate int, err error)
}

type grpcConnector struct {
	exchange.UnimplementedExchangeServiceServer
	sex IServiceExchange
}

func (con *grpcConnector) GetExchangeRates(
	ctx context.Context,
	in *exchange.Empty,
) (*exchange.RatesResponse, error) {

	// заглушка
	exchangeRate, _ := con.sex.Exchanges(ctx)
	_ = exchangeRate

	exchangeRate = map[string]float64{
		"rub":  1,
		"euro": 112,
		"usd":  100,
	}

	return &exchange.RatesResponse{
		Rates: exchangeRate,
	}, nil
}

func (con *grpcConnector) GetExchangeRate(
	ctx context.Context,
	in *exchange.CurrencyRequest,
) (*exchange.CurrencyResponse, error) {
	if in.FromCurrency == "" {
		return nil, status.Error(codes.InvalidArgument, "from_currency is required")
	}

	if in.ToCurrency == "" {
		return nil, status.Error(codes.InvalidArgument, "to_Currency is required")
	}

	// заглушка
	exchangeRate, _ := con.sex.Exchange(ctx, in.FromCurrency, in.ToCurrency)
	_ = exchangeRate

	return &exchange.CurrencyResponse{
		FromCurrency: "from Rub",
		ToCurrency:   "convert Bat",
		Rate:         float64(99),
	}, nil

}
