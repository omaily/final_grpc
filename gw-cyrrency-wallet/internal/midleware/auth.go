package midleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/auth"
)

func MiddlewareOne() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{"message": "request does not contain an access token"})
			return
		}

		err := auth.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err})
			return
		}

		log.Print(token)
	}
}
