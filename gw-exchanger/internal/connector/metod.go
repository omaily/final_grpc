package connector

import (
	"context"

	"github.com/omaily/final_grpc/gw-exchanger/internal/storage"
	exchange "github.com/omaily/final_grpc/gw-exchanger/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IServiceExchange interface {
	ExchangeRates(
		ctx context.Context,
	) (rates []storage.Exchange, err error)

	ExchangeCurrency(
		ctx context.Context,
		from_currency string,
		to_currency string,
	) (rateOdds float64, err error)
}

type grpcConnector struct {
	exchange.UnimplementedExchangeServiceServer
	sex IServiceExchange
}

func (con *grpcConnector) GetExchangeRates(
	ctx context.Context,
	in *exchange.Empty,
) (*exchange.RatesResponse, error) {
	exchangeRates, _ := con.sex.ExchangeRates(ctx)
	mExchangeRates := make(map[string]float64)
	for _, ex := range exchangeRates {
		mExchangeRates[ex.Note] = ex.Rate
	}

	return &exchange.RatesResponse{
		Rates: mExchangeRates,
	}, nil
}

func (con *grpcConnector) GetExchangeCurrency(
	ctx context.Context,
	in *exchange.CurrencyRequest,
) (*exchange.CurrencyResponse, error) {
	if in.FromCurrency == "" {
		return nil, status.Error(codes.InvalidArgument, "from_currency is required")
	}

	if in.ToCurrency == "" {
		return nil, status.Error(codes.InvalidArgument, "to_Currency is required")
	}

	odds, _ := con.sex.ExchangeCurrency(ctx, in.FromCurrency, in.ToCurrency)

	return &exchange.CurrencyResponse{
		FromCurrency: in.FromCurrency,
		ToCurrency:   in.ToCurrency,
		Rate:         odds,
	}, nil

}
