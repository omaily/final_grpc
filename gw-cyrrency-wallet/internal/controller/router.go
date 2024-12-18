package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/midleware"
)

func (s *Http) router() http.Handler {
	router := gin.Default()

	public(router) // Публичные маршруты
	secure(router)

	router.Run()
	return router
}

func public(router *gin.Engine) {
	router.POST("/api/v1/register", register)
	router.POST("/api/v1/login", login)
}

func secure(router *gin.Engine) {
	gr := router.Group("")
	gr.Use(midleware.MiddlewareOne())

	gr.GET("/api/v1/balance", balance)
	gr.GET("/api/v1/exchange/rates", rates)
	gr.POST("/api/v1/exchange", exchange)
	gr.POST("/api/v1/wallet/deposit", deposit)
	gr.POST("/api/v1/wallet/withdraw", withdraw)
}
