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
	"github.com/tkoleo84119/nail-salon-backend/internal/infra/redis"
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

	database, dbCleanup, err := db.New(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbCleanup()

	redisClient, redisCleanup, err := redis.NewClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatalf("Failed to initialize redis: %v", err)
	}
	defer redisCleanup()

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

	container, err := app.NewContainer(cfg, database, redisClient)
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

	router := app.SetupRoutes(container)

	// start line reminder job
	if err := container.GetJobs().LineReminderJob.Start(); err != nil {
		log.Fatalf("Failed to start line reminder job: %v", err)
	}
	defer container.GetJobs().LineReminderJob.Stop()

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
