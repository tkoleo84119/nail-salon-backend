package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/handler"
	authHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/auth"
	staffHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/staff"
	storeAccessHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/store-access"
	stylistHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/stylist"
	"github.com/tkoleo84119/nail-salon-backend/internal/infra/db"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	staffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	authService "github.com/tkoleo84119/nail-salon-backend/internal/service/auth"
	staffService "github.com/tkoleo84119/nail-salon-backend/internal/service/staff"
	storeAccessService "github.com/tkoleo84119/nail-salon-backend/internal/service/store-access"
	stylistService "github.com/tkoleo84119/nail-salon-backend/internal/service/stylist"
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

	// initialize error manager
	errorManager := errorCodes.GetManager()
	if err := errorManager.LoadFromFile("internal/errors/errors.yaml"); err != nil {
		log.Fatalf("Failed to load error definitions: %v", err)
	}

	// initialize sqlc queries
	queries := dbgen.New(database.PgxPool)

	// initialize repositories
	stylistRepository := sqlx.NewStylistRepository(database.Sqlx)

	// initialize services
	authLoginService := authService.NewLoginService(queries, cfg.JWT)
	staffCreateService := staffService.NewCreateStaffService(queries, database.PgxPool)
	staffUpdateService := staffService.NewUpdateStaffService(queries, database.Sqlx)
	staffStoreAccessService := storeAccessService.NewCreateStoreAccessService(queries)
	staffDeleteStoreAccessService := storeAccessService.NewDeleteStoreAccessService(queries)
	stylistCreateService := stylistService.NewCreateStylistService(queries)
	stylistUpdateService := stylistService.NewUpdateStylistService(queries, stylistRepository)

	// initialize handlers
	authLoginHandler := authHandler.NewLoginHandler(authLoginService)
	staffCreateHandler := staffHandler.NewCreateStaffHandler(staffCreateService)
	staffUpdateHandler := staffHandler.NewUpdateStaffHandler(staffUpdateService)
	staffStoreAccessHandler := storeAccessHandler.NewCreateStoreAccessHandler(staffStoreAccessService)
	staffDeleteStoreAccessHandler := storeAccessHandler.NewDeleteStoreAccessHandler(staffDeleteStoreAccessService)
	stylistCreateHandler := stylistHandler.NewCreateStylistHandler(stylistCreateService)
	stylistUpdateHandler := stylistHandler.NewUpdateStylistHandler(stylistUpdateService)

	router := gin.Default()

	// routes
	router.GET("/health", handler.Health)
	api := router.Group("/api")
	{
		staff := api.Group("/staff")
		{
			staff.POST("/login", authLoginHandler.Login)
			staff.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), staffCreateHandler.CreateStaff)
			staff.PATCH("/:id", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), staffUpdateHandler.UpdateStaff)
			staff.POST("/:id/store-access", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), staffStoreAccessHandler.CreateStoreAccess)
			staff.DELETE("/:id/store-access", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), staffDeleteStoreAccessHandler.DeleteStoreAccess)
		}

		stylists := api.Group("/stylists")
		{
			stylists.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireRoles(staffModel.RoleAdmin, staffModel.RoleManager, staffModel.RoleStylist), stylistCreateHandler.CreateStylist)
			stylists.PATCH("/me", middleware.JWTAuth(*cfg, queries), middleware.RequireRoles(staffModel.RoleAdmin, staffModel.RoleManager, staffModel.RoleStylist), stylistUpdateHandler.UpdateStylist)
		}
	}

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
