package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	"github.com/tkoleo84119/nail-salon-backend/internal/app"
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/infra/db"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

func init() {
	if gin.Mode() != gin.ReleaseMode {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
}

func main() {
	cfg := config.Load()

	database, cleanup, err := db.New(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer cleanup()

	if err := utils.InitSnowflake(cfg.Server.SnowflakeNodeId); err != nil {
		log.Fatalf("Failed to initialize snowflake: %v", err)
	}

	errorManager := errorCodes.GetManager()
	if err := errorManager.LoadFromFile("internal/errors/errors.json"); err != nil {
		log.Fatalf("Failed to load error definitions: %v", err)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("noBlank", utils.NoBlank)
		v.RegisterValidation("taiwanlandline", utils.ValidateTaiwanLandline)
		v.RegisterValidation("taiwanmobile", utils.ValidateTaiwanMobile)
		v.RegisterValidation("taiwanphone", utils.ValidateTaiwanPhone)
	}

	container := app.NewContainer(cfg, database)
	router := app.SetupRoutes(container)

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
