package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jabbala/go-wallet/handlers"
)

func main() {
	r := gin.Default()

	// Define routes
	r.GET("/generate", handlers.GenerateKeyPair)
	r.GET("/address", handlers.GetAddress)
	r.POST("/sign", handlers.SignMessage)
	r.POST("/verify", handlers.VerifyMessage)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
