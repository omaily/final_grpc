package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/omaily/final_grpc/gw-cyrrency-wallet/config"
	pb "github.com/omaily/final_grpc/gw-cyrrency-wallet/pkg/proto"
)

type GrpcClient struct {
	conf   config.GRPCServer
	conn   *grpc.ClientConn
	Client *pb.ExchangeServiceClient
}

func New(conf config.GRPCServer) *GrpcClient {
	return &GrpcClient{conf: conf}
}

func (c *GrpcClient) Start(_ context.Context) error {
	address := fmt.Sprintf("%s:%d", c.conf.Address, c.conf.Port)
	listen, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	c.conn = listen

	client := pb.NewExchangeServiceClient(listen)
	c.Client = &client

	return nil
}

func (c *GrpcClient) Stop() {
	c.conn.Close()
}
