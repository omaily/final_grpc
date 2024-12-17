package controller

import (
	"net/http"

	h "github.com/omaily/final_grpc/gw-exchanger/internal/delivery"
)

func (s *Http) router() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /api/v1/register", h.Register)
	router.HandleFunc("POST /api/v1/login", h.Login)
	router.HandleFunc("GET  /api/v1/balance", h.Balance)
	router.HandleFunc("POST /api/v1/wallet/deposit", h.Deposit)
	router.HandleFunc("POST /api/v1/wallet/withdraw", h.Withdraw)
	router.HandleFunc("GET  /api/v1/exchange/rates", h.Rates)
	router.HandleFunc("POST /api/v1/exchange", h.Exchange)

	return router
}
