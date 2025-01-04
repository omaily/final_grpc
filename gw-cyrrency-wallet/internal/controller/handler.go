package controller

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/auth"

	connectGrpc "github.com/omaily/final_grpc/gw-cyrrency-wallet/connection/grpc"
	connectRedis "github.com/omaily/final_grpc/gw-cyrrency-wallet/connection/redis"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/midleware"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/internal/storage"

	"github.com/omaily/final_grpc/gw-cyrrency-wallet/pkg/model"
	pb "github.com/omaily/final_grpc/gw-cyrrency-wallet/pkg/proto"
)

func test(red connectRedis.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		sender := "{" + c.Request.Host + ", Agent:" + c.Request.UserAgent() + "}"

		err := red.Set(context.Background(), "greeting", sender, 0)
		if err != nil {
			fmt.Println("Failed to set key:", err)
			return
		}

		val, err := red.Get(context.Background(), "greeting")
		if err != nil {
			fmt.Println("Failed to get key:", err)
			return
		}

		fmt.Println("Value for key 'greeting':", val)
	}
}

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

func register(st *storage.Instance) gin.HandlerFunc {
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

func login(st *storage.Instance) gin.HandlerFunc {
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

		fmt.Printf("%x\n", userId)
		accessToken, refreshToken, err := auth.GeneratePairToken(fmt.Sprintf("%x", userId))
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

func balance(st *storage.Instance, userId *string) gin.HandlerFunc {
	return func(c *gin.Context) {
		balance, err := st.CheckBalance(c.Request.Context(), *userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": balance})
	}
}

func deposit(st *storage.Instance, userId *string) gin.HandlerFunc {
	return func(c *gin.Context) {

		var json midleware.Deposit
		if err := c.ShouldBindJSON(&json); err != nil {
			slog.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount or currency"})
			return
		}

		deposit := model.Deposit(json)
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

func withdraw(st *storage.Instance, userId *string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json midleware.Deposit
		if err := c.ShouldBindJSON(&json); err != nil {
			slog.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount or currency"})
			return
		}

		deposit := model.Deposit(json)
		if err := deposit.Validate(); err != nil {
			slog.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		balance, err := st.GetMoney(c.Request.Context(), *userId, deposit)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Insufficient funds or invalid amount", "new_balance": balance})
	}
}

func rates(clientGrpc *connectGrpc.GrpcClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, err := (*clientGrpc.Client).GetExchangeRates(context.Background(), &pb.Empty{})
		if err != nil {
			slog.Error("не удалось создать: отправить Rates: ", slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}
		fmt.Println("return map:", r)
		c.JSON(http.StatusOK, gin.H{"message": "this is handler POST RATES"})
	}
}

func exchange(clientGrpc *connectGrpc.GrpcClient) gin.HandlerFunc {
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
