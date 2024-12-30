package connector

import (
	"fmt"
	"log/slog"

	"google.golang.org/grpc"

	"github.com/omaily/final_grpc/gw-cyrrency-wallet/config"
	pb "github.com/omaily/final_grpc/gw-cyrrency-wallet/pkg/proto"
)

type GrpcClient struct {
	conf   config.GRPCServer
	Conn   *grpc.ClientConn
	Client *pb.ExchangeServiceClient
}

func New(conf config.GRPCServer) *GrpcClient {
	address := fmt.Sprintf("%s:%d", conf.Address, conf.Port)
	listen, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		slog.Error("Не могу подключиться: %v", slog.String("error", err.Error()))
	}

	client := pb.NewExchangeServiceClient(listen)
	return &GrpcClient{
		conf:   conf,
		Conn:   listen,
		Client: &client,
	}
}

func (c *GrpcClient) Stop() {
	c.Conn.Close()
	slog.Info("...down grpc connector")
}
