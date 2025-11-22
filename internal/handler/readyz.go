package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tkoleo84119/nail-salon-backend/internal/infra/db"
)

func Readyz(database *db.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		// check if the database is ready
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := database.PgxPool.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "error",
				"error":  "database is not ready",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}
