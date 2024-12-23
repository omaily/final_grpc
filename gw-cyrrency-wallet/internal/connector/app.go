package connector

import (
	"log/slog"

	"google.golang.org/grpc"

	pb "github.com/omaily/final_grpc/gw-cyrrency-wallet/pkg/proto"
)

type GrpcClient struct {
	conf   string
	Conn   *grpc.ClientConn
	Client *pb.ExchangeServiceClient
}

func New(address string) *GrpcClient {
	listen, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		slog.Error("Не могу подключиться: %v", slog.String("error", err.Error()))
	}

	client := pb.NewExchangeServiceClient(listen)
	return &GrpcClient{
		conf:   address,
		Conn:   listen,
		Client: &client,
	}
}

func (c *GrpcClient) Stop() {
	c.Conn.Close()
	slog.Info("...down grpc connector")
}
