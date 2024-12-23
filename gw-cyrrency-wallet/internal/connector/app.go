package connector

import (
	"log/slog"

	"google.golang.org/grpc"

	pb "github.com/omaily/final_grpc/gw-cyrrency-wallet/pkg/gen"
)

type GrpcClient struct {
	conf   string
	Conn   *grpc.ClientConn
	Client *pb.ExchangeServiceClient
}

func New(address string) *GrpcClient {
	lis, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		slog.Error("Не могу подключиться: %v", slog.String("error", err.Error()))
	}

	client := pb.NewExchangeServiceClient(lis)
	return &GrpcClient{
		conf:   address,
		Conn:   lis,
		Client: &client,
	}
}

func (c *GrpcClient) Stop() {
	c.Conn.Close()
	slog.Info("...down grpc connector")
}
