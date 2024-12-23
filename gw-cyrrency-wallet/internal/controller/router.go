package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Http) router() http.Handler {
	router := gin.Default()

	s.public(router) //Публичные маршруты
	s.secure(router)

	router.Run()
	return router
}

func (s *Http) public(router *gin.Engine) {
	router.POST("/api/v1/register", register(s.storage))
	router.POST("/api/v1/login", login(s.storage))
}

func (s *Http) secure(router *gin.Engine) {
	gr := router.Group("")
	// gr.Use(midleware.MiddlewareOne())

	gr.GET("/api/v1/balance", balance)
	gr.POST("/api/v1/wallet/deposit", deposit)
	gr.POST("/api/v1/wallet/withdraw", withdraw)

	gr.GET("/api/v1/exchange/rates", rates(s.clientGrpc))
	gr.POST("/api/v1/exchange", exchange(s.clientGrpc))
}
