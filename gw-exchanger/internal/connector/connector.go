package connector

import (
	"github.com/omaily/final_grpc/gw-exchanger/internal/storage"
	exchange "github.com/omaily/final_grpc/gw-exchanger/pkg/gen"
	"google.golang.org/grpc"
)

type ServerGrpc struct {
	*grpc.Server
}

func New() *ServerGrpc {
	server := grpc.NewServer()
	st := &storage.Exchange{}

	exchange.RegisterExchangeServiceServer(server, &grpcConnector{sex: st})

	return &ServerGrpc{server}
}
