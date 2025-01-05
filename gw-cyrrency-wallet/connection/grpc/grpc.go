package grpc

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"

	"github.com/omaily/final_grpc/gw-cyrrency-wallet/config"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/pkg/model"
	pb "github.com/omaily/final_grpc/gw-cyrrency-wallet/pkg/proto"
)

type Client interface {
	Start(ctx context.Context) error
	Stop()

	ExchangeRates(ctx context.Context) (map[string]float64, error)
	ExchangeCurency(ctx context.Context, ex model.Exchange) (float64, error)
}

type grpcClient struct {
	conf   config.GRPCServer
	conn   *grpc.ClientConn
	client *pb.ExchangeServiceClient
}

func New(conf config.GRPCServer) Client {
	return &grpcClient{conf: conf}
}

func (g *grpcClient) Start(_ context.Context) error {
	address := fmt.Sprintf("%s:%d", g.conf.Address, g.conf.Port)

	listen, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	g.conn = listen

	client := pb.NewExchangeServiceClient(listen)
	g.client = &client

	return nil
}

func (g *grpcClient) Stop() {
	g.conn.Close()
}

func (g *grpcClient) ExchangeRates(ctx context.Context) (map[string]float64, error) {
	currency, err := (*g.client).GetExchangeRates(context.Background(), &pb.Empty{})
	if err != nil {
		slog.Error("failed to call grpc-method: GetExchangeRates: ", slog.String("error", err.Error()))
		return nil, err
	}
	return currency.Rates, nil
}

func (g *grpcClient) ExchangeCurency(ctx context.Context, ex model.Exchange) (float64, error) {
	exResult, err := (*g.client).GetExchangeCurrency(context.Background(), &pb.CurrencyRequest{
		FromCurrency: ex.FromCurrency,
		ToCurrency:   ex.ToCurrency,
	})
	if err != nil {
		slog.Error(err.Error())
		return 0, fmt.Errorf("failed to call grpc-method: GetExchangeCurency")
	}

	fmt.Println("from:", exResult.FromCurrency, ", to:", exResult.ToCurrency, ", rate:", exResult.Rate)
	return exResult.Rate, nil
}
