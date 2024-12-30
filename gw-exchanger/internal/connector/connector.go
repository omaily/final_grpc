package connector

import (
	exchange "github.com/omaily/final_grpc/gw-exchanger/pkg/proto"

	"github.com/omaily/final_grpc/gw-exchanger/internal/storage"
	"google.golang.org/grpc"
)

type ServerGrpc struct {
	*grpc.Server
}

func New(st *storage.Instance) *ServerGrpc {
	server := grpc.NewServer()
	exchange.RegisterExchangeServiceServer(server, &grpcConnector{sex: st})
	return &ServerGrpc{server}
}
