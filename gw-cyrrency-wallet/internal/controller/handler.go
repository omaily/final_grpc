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

	"github.com/omaily/final_grpc/gw-cyrrency-wallet/pkg/model"
	pb "github.com/omaily/final_grpc/gw-cyrrency-wallet/pkg/proto"
)

func parseAutorizate(c *gin.Context) *model.Account {
	logger := slog.With(
		slog.String("HandlerFunc", "parseAutorizate"),
	)

	var json midleware.Login
	if err := c.ShouldBindJSON(&json); err != nil {
		logger.Error("faled to decode json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil
	}

	if json.Username == "" || json.Password == "" {
		logger.Error("empty required fields are missing")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "empty required fields are missing"})
		return nil
	}
	acc := model.Account(json)
	return &acc
}

func register(st *storage.Instance) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := slog.With(
			slog.String("HandlerFunc", "register"),
		)

		user := parseAutorizate(c)
		if user == nil {
			logger.Error("invalid request parameters")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request parameters"})
			return
		}

		uuid, err := st.CreateAccount(c.Request.Context(), user)
		if err != nil {
			logger.Error("create account", slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		logger.Info(fmt.Sprintf("create account %x", *uuid))
		c.JSON(http.StatusOK, gin.H{"status": "User registered successfully"})
	}
}

func login(st *storage.Instance) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := slog.With(
			slog.String("HandlerFunc", "login"),
		)

		user := parseAutorizate(c)
		if user == nil {
			logger.Error("invalid request parameters")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request parameters"})
			return
		}

		passwordCript, err := st.FindAccount(c.Request.Context(), user)
		if err != nil {
			logger.Error("err", slog.String("error", err.Error()))
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Invalid username or password"})
			return
		}

		if !user.CheckPassword(passwordCript) {
			logger.Error("wrong password")
			c.JSON(http.StatusOK, gin.H{"error": "passwords do not match"})
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
