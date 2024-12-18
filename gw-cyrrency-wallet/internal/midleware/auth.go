package midleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MiddlewareOne() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Print("Executing middlewareOne")
	}
}

func middlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing middleware  Two")
	})
}
