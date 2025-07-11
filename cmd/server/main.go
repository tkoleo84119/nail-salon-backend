package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/tkoleo84119/nail-salon-backend/internal/handler"
)

// init loads environment variables from a local .env file when not running in Release mode.
func init() {
	if gin.Mode() != gin.ReleaseMode {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	router := gin.Default()

	router.GET("/health", handler.Health)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
