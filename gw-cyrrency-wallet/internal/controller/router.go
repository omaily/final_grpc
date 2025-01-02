package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/midleware"
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
	var userId string
	gr := router.Group("/api/v1").Use(midleware.AuthCookie(&userId))
	{
		gr.GET("/balance", balance(s.storage, &userId))
		gr.POST("/wallet/deposit", deposit(s.storage, &userId))
		gr.POST("/wallet/withdraw", withdraw)

		gr.GET("/exchange/rates", rates(s.clientGrpc))
		gr.POST("/exchange", exchange(s.clientGrpc))
	}
}
