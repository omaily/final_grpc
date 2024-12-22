package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/midleware"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/storage"
)

func (s *Http) router() http.Handler {
	router := gin.Default()

	public(router, s.storage) // Публичные маршруты
	secure(router, s.storage)

	router.Run()
	return router
}

func public(router *gin.Engine, st *storage.Instance) {
	router.POST("/api/v1/register", register(st))
	router.POST("/api/v1/login", login(st))
}

func secure(router *gin.Engine, st *storage.Instance) {
	gr := router.Group("")
	gr.Use(midleware.MiddlewareOne())

	gr.GET("/api/v1/balance", balance)
	gr.GET("/api/v1/exchange/rates", rates)
	gr.POST("/api/v1/exchange", exchange)
	gr.POST("/api/v1/wallet/deposit", deposit)
	gr.POST("/api/v1/wallet/withdraw", withdraw)
}
