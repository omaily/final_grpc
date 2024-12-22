package connector

import (
	"google.golang.org/grpc"

	exchange "github.com/omaily/final_grpc/gw-cyrrency-wallet/gen"

	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/storage"
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
