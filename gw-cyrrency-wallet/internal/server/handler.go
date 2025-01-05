package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	connGrpc "github.com/omaily/final_grpc/gw-cyrrency-wallet/connection/grpc"
	connRedis "github.com/omaily/final_grpc/gw-cyrrency-wallet/connection/redis"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/connection/storage"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/auth"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/midleware"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/pkg/model"
)

func parsAuthenticat(c *gin.Context) *model.Account {
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

func register(st storage.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := slog.With(
			slog.String("HandlerFunc", "register"),
		)

		user := parsAuthenticat(c)
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

		logger.Info(fmt.Sprintf("create account %x", uuid))
		c.JSON(http.StatusOK, gin.H{"status": "User registered successfully"})
	}
}

func login(st storage.Repository, rd connRedis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := slog.With(
			slog.String("HandlerFunc", "login"),
		)

		user := parsAuthenticat(c)
		if user == nil {
			logger.Error("invalid request parameters")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request parameters"})
			return
		}

		userId, passwordCript, err := st.FindAccount(c.Request.Context(), user)
		if err != nil {
			logger.Error("err", slog.String("error", err.Error()))
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Invalid username or password"})
			return
		}

		if !user.CheckPassword(passwordCript) {
			logger.Error("wrong password")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "passwords do not match"})
			return
		}

		accessToken, refreshToken, err := auth.GeneratePairToken(rd, fmt.Sprintf("%x", userId))
		if err != nil {
			logger.Error("error creating token", slog.String("err", err.Error()))
			c.JSON(http.StatusUnauthorized, gin.H{"err": err})
			return
		}

		http.SetCookie(c.Writer, accessToken)
		http.SetCookie(c.Writer, refreshToken)
		c.JSON(http.StatusOK, gin.H{"token": midleware.Bearer(accessToken.Value)})
	}
}

func balance(st storage.Repository, userId *string) gin.HandlerFunc {
	return func(c *gin.Context) {
		balance, err := st.CheckBalance(c.Request.Context(), *userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": balance})
	}
}

func deposit(st storage.Repository, userId *string) gin.HandlerFunc {
	return func(c *gin.Context) {

		var json midleware.Transfer
		if err := c.ShouldBindJSON(&json); err != nil {
			slog.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount or currency"})
			return
		}

		deposit := model.Transfer(json)
		if err := deposit.Validate(); err != nil {
			slog.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		balance, err := st.PutMoney(c.Request.Context(), *userId, deposit)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "account topped up successfully", "new_balance": balance})
	}
}

func withdraw(st storage.Repository, userId *string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json midleware.Transfer
		if err := c.ShouldBindJSON(&json); err != nil {
			slog.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount or currency"})
			return
		}

		withdraw := model.Transfer(json)
		if err := withdraw.Validate(); err != nil {
			slog.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		balance, err := st.TakeMoney(c.Request.Context(), *userId, withdraw)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Insufficient funds or invalid amount", "new_balance": balance})
	}
}

func rates(grpclient connGrpc.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		rates, err := grpclient.ExchangeRates(context.Background())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err})
		}

		c.JSON(http.StatusOK, gin.H{"rates": rates})
	}
}

func exchange(grpclient connGrpc.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json midleware.Exchange
		if err := c.ShouldBindJSON(&json); err != nil {
			slog.Error(fmt.Sprintf("json fields are incorrect: %s", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "json fields are incorrect"})
			return
		}

		ex := model.Exchange(json)
		odds, err := grpclient.ExchangeCurency(context.Background(), ex)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err})
		}

		c.JSON(http.StatusOK, gin.H{"message": "this is handler POST EXCHANGE", "coofisient": odds})
	}
}
