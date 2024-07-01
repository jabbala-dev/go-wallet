package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jabbala-dev/go-wallet/services"
)

func GenerateKeyPair(c *gin.Context) {
	privateKey, address, err := services.GenerateKeyPair()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"private_key": privateKey, "address": address})
}

func GetAddress(c *gin.Context) {
	address, err := services.GetAddress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"address": address})
}

func SignMessage(c *gin.Context) {
	var request struct {
		Message string `json:"message"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	signature, err := services.SignMessage(request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"signature": signature})
}

func VerifyMessage(c *gin.Context) {
	var request struct {
		Message   string `json:"message"`
		Signature string `json:"signature"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	isValid, err := services.VerifyMessage(request.Message, request.Signature)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": isValid})
}
