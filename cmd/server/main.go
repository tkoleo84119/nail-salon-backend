package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	"github.com/tkoleo84119/nail-salon-backend/internal/handler"
	staffHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/infra/db"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	staffService "github.com/tkoleo84119/nail-salon-backend/internal/service/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
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
	// load config
	cfg := config.Load()

	// connect to database
	database, cleanup, err := db.New(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer cleanup()

	// initialize snowflake
	if err := utils.InitSnowflake(cfg.Server.SnowflakeNodeId); err != nil {
		log.Fatalf("Failed to initialize snowflake: %v", err)
	}

	// initialize sqlc queries
	queries := dbgen.New(database.PgxPool)

	// initialize services
	staffLoginService := staffService.NewLoginService(queries, cfg.JWT)
	staffCreateService := staffService.NewCreateStaffService(queries, database.PgxPool)

	// initialize handlers
	staffLoginHandler := staffHandler.NewLoginHandler(staffLoginService)
	staffCreateHandler := staffHandler.NewCreateStaffHandler(staffCreateService)

	router := gin.Default()

	// routes
	router.GET("/health", handler.Health)
	api := router.Group("/api")
	{
		staff := api.Group("/staff")
		{
			staff.POST("/login", staffLoginHandler.Login)
			staff.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), staffCreateHandler.CreateStaff)
		}
	}

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
