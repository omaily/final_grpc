package midleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/auth"
)

func AuthCookie(userId *string) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Request.Cookie("access_token")
		if err != nil {
			slog.Error(err.Error())
			if err == http.ErrNoCookie {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "request does not contain cookie"})
				return
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "request does not contain an access token"})
			return
		}

		*userId, err = auth.ValidateToken(accessToken.Value)
		if err != nil {
			slog.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "access_token not valid"})
			return
		}
	}
}

func AuthHeader(userId *string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is missing"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authToken := strings.Split(authHeader, " ")
		if len(authToken) != 2 || authToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var err error
		*userId, err = auth.ValidateToken(authToken[1])
		if err != nil {
			slog.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "access_token not valid"})
			return
		}
	}
}
