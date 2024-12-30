package controller

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/connector"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/midleware"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/storage"

	pb "github.com/omaily/final_grpc/gw-cyrrency-wallet/pkg/proto"
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

		st.User.CreateAccount(json.Username, json.Password)
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

func rates(clientGrpc *connector.GrpcClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, err := (*clientGrpc.Client).GetExchangeRates(context.Background(), &pb.Empty{})
		if err != nil {
			slog.Error("Не удалось создать: отправить Rates: ", slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}
		fmt.Println("return map:", r)
		c.JSON(http.StatusOK, gin.H{"message": "this is handler POST RATES"})
	}
}

func exchange(clientGrpc *connector.GrpcClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, err := (*clientGrpc.Client).GetExchangeRate(context.Background(), &pb.CurrencyRequest{
			FromCurrency: "usd",
			ToCurrency:   "rub",
		})
		if err != nil {
			slog.Error("Не удалось создать: отправить Rate", slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}

		slog.Info("from: %v, to: %v, rate: %f", r.FromCurrency, r.ToCurrency, r.Rate)
		c.JSON(http.StatusOK, gin.H{"message": "this is handler POST EXCHANGE"})
	}
}
