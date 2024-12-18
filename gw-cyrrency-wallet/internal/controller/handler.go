package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func register(c *gin.Context) {
	status := true

	if status {
		c.JSON(http.StatusOK, gin.H{"message": "this is handler POST REGISTER"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "error"})
	}
}

func login(c *gin.Context) {
	status := true

	if status {
		c.JSON(http.StatusOK, gin.H{"message": "this is handler POST LOGIN"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "error"})
	}
}

func balance(c *gin.Context) {
	status := true

	if status {
		c.JSON(http.StatusOK, gin.H{"message": "this is handler POST BALANCE"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "error"})
	}
}

func deposit(c *gin.Context) {
	status := true

	if status {
		c.JSON(http.StatusOK, gin.H{"message": "this is handler POST DEPOSiT"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "error"})
	}
}

func withdraw(c *gin.Context) {
	status := true

	if status {
		c.JSON(http.StatusOK, gin.H{"message": "this is handler POST WITHDRAW"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "error"})
	}
}

func rates(c *gin.Context) {
	status := true

	if status {
		c.JSON(http.StatusOK, gin.H{"message": "this is handler POST RATES"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "error"})
	}
}

func exchange(c *gin.Context) {
	status := true

	if status {
		c.JSON(http.StatusOK, gin.H{"message": "this is handler POST EXCHANGE"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "error"})
	}
}
