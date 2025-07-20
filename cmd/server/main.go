package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/handler"
	authHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/auth"
	customerHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/customer"
	scheduleHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/schedule"
	serviceHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/service"
	staffHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/staff"
	storeHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/store"
	stylistHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/stylist"
	timeSlotTemplateHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/infra/db"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	staffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	authService "github.com/tkoleo84119/nail-salon-backend/internal/service/auth"
	customerService "github.com/tkoleo84119/nail-salon-backend/internal/service/customer"
	scheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/schedule"
	serviceService "github.com/tkoleo84119/nail-salon-backend/internal/service/service"
	staffService "github.com/tkoleo84119/nail-salon-backend/internal/service/staff"
	storeService "github.com/tkoleo84119/nail-salon-backend/internal/service/store"
	stylistService "github.com/tkoleo84119/nail-salon-backend/internal/service/stylist"
	timeSlotTemplateService "github.com/tkoleo84119/nail-salon-backend/internal/service/time-slot-template"
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

	// register custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("taiwanlandline", utils.ValidateTaiwanLandline)
	}

	// initialize sqlc queries
	queries := dbgen.New(database.PgxPool)

	// initialize repositories
	stylistRepository := sqlx.NewStylistRepository(database.Sqlx)
	staffUserRepository := sqlx.NewStaffUserRepository(database.Sqlx)
	storeRepository := sqlx.NewStoreRepository(database.Sqlx)
	serviceRepository := sqlx.NewServiceRepository(database.Sqlx)

	// initialize services
	authLoginService := authService.NewLoginService(queries, cfg.JWT)
	customerLineLoginService := customerService.NewLineLoginService(queries, cfg.Line, cfg.JWT)
	staffCreateService := staffService.NewCreateStaffService(queries, database.PgxPool)
	staffUpdateService := staffService.NewUpdateStaffService(queries, database.Sqlx)
	staffUpdateMeService := staffService.NewUpdateMyStaffService(queries, staffUserRepository)
	staffStoreAccessService := staffService.NewCreateStoreAccessService(queries)
	staffDeleteStoreAccessService := staffService.NewDeleteStoreAccessBulkService(queries)
	storeCreateService := storeService.NewCreateStoreService(queries, database.PgxPool)
	storeUpdateService := storeService.NewUpdateStoreService(queries, storeRepository)
	stylistCreateService := stylistService.NewCreateMyStylistService(queries)
	stylistUpdateService := stylistService.NewUpdateMyStylistService(queries, stylistRepository)
	scheduleCreateBulkService := scheduleService.NewCreateSchedulesBulkService(queries, database.PgxPool)
	scheduleDeleteBulkService := scheduleService.NewDeleteSchedulesBulkService(queries, database.PgxPool)
	scheduleCreateTimeSlotService := scheduleService.NewCreateTimeSlotService(queries)
	timeSlotRepository := sqlx.NewTimeSlotRepository(database.Sqlx)
	scheduleUpdateTimeSlotService := scheduleService.NewUpdateTimeSlotService(queries, timeSlotRepository)
	scheduleDeleteTimeSlotService := scheduleService.NewDeleteTimeSlotService(queries)
	timeSlotTemplateCreateService := timeSlotTemplateService.NewCreateTimeSlotTemplateService(queries, database.PgxPool)
	timeSlotTemplateRepository := sqlx.NewTimeSlotTemplateRepository(database.Sqlx)
	timeSlotTemplateUpdateService := timeSlotTemplateService.NewUpdateTimeSlotTemplateService(queries, timeSlotTemplateRepository)
	timeSlotTemplateDeleteService := timeSlotTemplateService.NewDeleteTimeSlotTemplateService(queries)
	timeSlotTemplateCreateItemService := timeSlotTemplateService.NewCreateTimeSlotTemplateItemService(queries)
	timeSlotTemplateUpdateItemService := timeSlotTemplateService.NewUpdateTimeSlotTemplateItemService(queries)
	timeSlotTemplateDeleteItemService := timeSlotTemplateService.NewDeleteTimeSlotTemplateItemService(queries)
	serviceCreateService := serviceService.NewCreateServiceService(queries)
	serviceUpdateService := serviceService.NewUpdateServiceService(queries, serviceRepository)

	// initialize handlers
	authLoginHandler := authHandler.NewLoginHandler(authLoginService)
	customerLineLoginHandler := customerHandler.NewLineLoginHandler(customerLineLoginService)
	staffCreateHandler := staffHandler.NewCreateStaffHandler(staffCreateService)
	staffUpdateHandler := staffHandler.NewUpdateStaffHandler(staffUpdateService)
	staffUpdateMeHandler := staffHandler.NewUpdateMyStaffHandler(staffUpdateMeService)
	staffStoreAccessHandler := staffHandler.NewCreateStoreAccessHandler(staffStoreAccessService)
	staffDeleteStoreAccessHandler := staffHandler.NewDeleteStoreAccessBulkHandler(staffDeleteStoreAccessService)
	storeCreateHandler := storeHandler.NewCreateStoreHandler(storeCreateService)
	storeUpdateHandler := storeHandler.NewUpdateStoreHandler(storeUpdateService)
	stylistCreateHandler := stylistHandler.NewCreateMyStylistHandler(stylistCreateService)
	stylistUpdateHandler := stylistHandler.NewUpdateMyStylistHandler(stylistUpdateService)
	scheduleCreateBulkHandler := scheduleHandler.NewCreateSchedulesBulkHandler(scheduleCreateBulkService)
	scheduleDeleteBulkHandler := scheduleHandler.NewDeleteSchedulesBulkHandler(scheduleDeleteBulkService)
	scheduleCreateTimeSlotHandler := scheduleHandler.NewCreateTimeSlotHandler(scheduleCreateTimeSlotService)
	scheduleUpdateTimeSlotHandler := scheduleHandler.NewUpdateTimeSlotHandler(scheduleUpdateTimeSlotService)
	scheduleDeleteTimeSlotHandler := scheduleHandler.NewDeleteTimeSlotHandler(scheduleDeleteTimeSlotService)
	timeSlotTemplateCreateHandler := timeSlotTemplateHandler.NewCreateTimeSlotTemplateHandler(timeSlotTemplateCreateService)
	timeSlotTemplateUpdateHandler := timeSlotTemplateHandler.NewUpdateTimeSlotTemplateHandler(timeSlotTemplateUpdateService)
	timeSlotTemplateDeleteHandler := timeSlotTemplateHandler.NewDeleteTimeSlotTemplateHandler(timeSlotTemplateDeleteService)
	timeSlotTemplateCreateItemHandler := timeSlotTemplateHandler.NewCreateTimeSlotTemplateItemHandler(timeSlotTemplateCreateItemService)
	timeSlotTemplateUpdateItemHandler := timeSlotTemplateHandler.NewUpdateTimeSlotTemplateItemHandler(timeSlotTemplateUpdateItemService)
	timeSlotTemplateDeleteItemHandler := timeSlotTemplateHandler.NewDeleteTimeSlotTemplateItemHandler(timeSlotTemplateDeleteItemService)
	serviceCreateHandler := serviceHandler.NewCreateServiceHandler(serviceCreateService)
	serviceUpdateHandler := serviceHandler.NewUpdateServiceHandler(serviceUpdateService)

	router := gin.Default()

	// routes
	router.GET("/health", handler.Health)
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			customer := auth.Group("/customer")
			{
				customer.POST("/line/login", customerLineLoginHandler.LineLogin)
			}
		}
		staff := api.Group("/staff")
		{
			staff.POST("/login", authLoginHandler.Login)
			staff.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), staffCreateHandler.CreateStaff)
			staff.PATCH("/:id", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), staffUpdateHandler.UpdateStaff)
			staff.PATCH("/me", middleware.JWTAuth(*cfg, queries), staffUpdateMeHandler.UpdateMyStaff)
			staff.POST("/:id/store-access", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), staffStoreAccessHandler.CreateStoreAccess)
			staff.DELETE("/:id/store-access", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), staffDeleteStoreAccessHandler.DeleteStoreAccessBulk)
		}

		stylists := api.Group("/stylists")
		{
			stylists.POST("/me", middleware.JWTAuth(*cfg, queries), middleware.RequireRoles(staffModel.RoleAdmin, staffModel.RoleManager, staffModel.RoleStylist), stylistCreateHandler.CreateMyStylist)
			stylists.PATCH("/me", middleware.JWTAuth(*cfg, queries), middleware.RequireRoles(staffModel.RoleAdmin, staffModel.RoleManager, staffModel.RoleStylist), stylistUpdateHandler.UpdateMyStylist)
		}

		schedules := api.Group("/schedules")
		{
			schedules.POST("/bulk", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), scheduleCreateBulkHandler.CreateSchedulesBulk)
			schedules.DELETE("/bulk", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), scheduleDeleteBulkHandler.DeleteSchedulesBulk)
			schedules.POST("/:scheduleId/time-slots", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), scheduleCreateTimeSlotHandler.CreateTimeSlot)
			schedules.PATCH("/:scheduleId/time-slots/:timeSlotId", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), scheduleUpdateTimeSlotHandler.UpdateTimeSlot)
			schedules.DELETE("/:scheduleId/time-slots/:timeSlotId", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), scheduleDeleteTimeSlotHandler.DeleteTimeSlot)
		}

		stores := api.Group("/stores")
		{
			stores.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), storeCreateHandler.CreateStore)
			stores.PATCH("/:storeId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), storeUpdateHandler.UpdateStore)
		}

		timeSlotTemplates := api.Group("/time-slot-templates")
		{
			timeSlotTemplates.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), timeSlotTemplateCreateHandler.CreateTimeSlotTemplate)
			timeSlotTemplates.PATCH("/:templateId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), timeSlotTemplateUpdateHandler.UpdateTimeSlotTemplate)
			timeSlotTemplates.DELETE("/:templateId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), timeSlotTemplateDeleteHandler.DeleteTimeSlotTemplate)
			timeSlotTemplates.POST("/:templateId/items", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), timeSlotTemplateCreateItemHandler.CreateTimeSlotTemplateItem)
			timeSlotTemplates.PATCH("/:templateId/items/:itemId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), timeSlotTemplateUpdateItemHandler.UpdateTimeSlotTemplateItem)
			timeSlotTemplates.DELETE("/:templateId/items/:itemId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), timeSlotTemplateDeleteItemHandler.DeleteTimeSlotTemplateItem)
		}

		services := api.Group("/services")
		{
			services.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), serviceCreateHandler.CreateService)
			services.PATCH("/:serviceId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), serviceUpdateHandler.UpdateService)
		}
	}

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
