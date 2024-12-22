package controller

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/midleware"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/storage"
)

func register(st *storage.Instance) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json midleware.Login
		if err := c.ShouldBindJSON(&json); err != nil {
			slog.Info("faled to decode json")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if json.Username == "" || json.Password == "" || json.Email == "" {
			slog.Info("Empty required fields are missing")
			c.JSON(http.StatusUnauthorized, gin.H{"status": "empty required fields are missing"})
			return
		}

		st.CreateAccount(json.Username, json.Password)
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	}
}

func login(st *storage.Instance) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json midleware.Login
		if err := c.ShouldBindJSON(&json); err != nil {
			slog.Info("faled to decode json")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if json.Username == "" || json.Password == "" {
			slog.Info("Empty required fields are missing")
			c.JSON(http.StatusUnauthorized, gin.H{"status": "empty required fields are missing"})
			return
		}

		if isExists := st.FindAccount(json.Username, json.Password); !isExists {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Invalid username or password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": "JWT_TOKEN"})
	}
}

func balance(c *gin.Context) {
	status := true

	if status {
		c.JSON(http.StatusOK, gin.H{"message": "this is handler POST BALANCE"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
	}
}

func deposit(c *gin.Context) {
	status := true

	if status {
		c.JSON(http.StatusOK, gin.H{"message": "this is handler POST DEPOSiT"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
	}
}

func withdraw(c *gin.Context) {
	status := true

	if status {
		c.JSON(http.StatusOK, gin.H{"message": "this is handler POST WITHDRAW"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
	}
}

func rates(c *gin.Context) {
	status := true

	if status {
		c.JSON(http.StatusOK, gin.H{"message": "this is handler POST RATES"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
	}
}

func exchange(c *gin.Context) {
	status := true

	if status {
		c.JSON(http.StatusOK, gin.H{"message": "this is handler POST EXCHANGE"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
	}
}
