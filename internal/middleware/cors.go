package middleware

import (
	"slices"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
)

// CORSMiddleware creates a CORS middleware using gin-contrib/cors with the given configuration
func CORSMiddleware(corsConfig config.CORSConfig) gin.HandlerFunc {
	corsConf := cors.Config{
		AllowOrigins:     corsConfig.AllowedOrigins,
		AllowMethods:     corsConfig.AllowedMethods,
		AllowHeaders:     corsConfig.AllowedHeaders,
		ExposeHeaders:    corsConfig.ExposedHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
		MaxAge:           time.Duration(corsConfig.MaxAge) * time.Second,
	}

	// Handle wildcard origins
	if slices.Contains(corsConfig.AllowedOrigins, "*") {
		corsConf.AllowAllOrigins = true
		corsConf.AllowOrigins = nil
	}

	return cors.New(corsConf)
}
