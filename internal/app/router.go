package app

import (
	"github.com/gin-gonic/gin"
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	"github.com/tkoleo84119/nail-salon-backend/internal/handler"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	staffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
)

func SetupRoutes(container *Container) *gin.Engine {
	cfg := container.GetConfig()
	database := container.GetDatabase()
	handlers := container.GetHandlers()
	
	queries := dbgen.New(database.PgxPool)
	router := gin.Default()

	router.GET("/health", handler.Health)
	
	api := router.Group("/api")
	{
		setupAuthRoutes(api, handlers)
		setupCustomerRoutes(api, cfg, queries, handlers)
		setupBookingRoutes(api, cfg, queries, handlers)
		setupStaffRoutes(api, cfg, queries, handlers)
		setupStylistRoutes(api, cfg, queries, handlers)
		setupScheduleRoutes(api, cfg, queries, handlers)
		setupStoreRoutes(api, cfg, queries, handlers)
		setupTimeSlotTemplateRoutes(api, cfg, queries, handlers)
		setupServiceRoutes(api, cfg, queries, handlers)
	}

	return router
}

func setupAuthRoutes(api *gin.RouterGroup, handlers Handlers) {
	auth := api.Group("/auth")
	{
		customer := auth.Group("/customer")
		{
			customer.POST("/line/login", handlers.CustomerLineLogin.LineLogin)
			customer.POST("/line/register", handlers.CustomerLineRegister.LineRegister)
		}
	}
}

func setupCustomerRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	customer := api.Group("/customers")
	{
		customer.PATCH("/me", middleware.JWTAuth(*cfg, queries), handlers.CustomerUpdateMy.UpdateMyCustomer)
	}
}

func setupBookingRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	bookings := api.Group("/bookings")
	{
		bookings.POST("/me", middleware.JWTAuth(*cfg, queries), handlers.BookingCreateMy.CreateMyBooking)
	}
}

func setupStaffRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	staff := api.Group("/staff")
	{
		staff.POST("/login", handlers.AuthLogin.Login)
		staff.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.StaffCreate.CreateStaff)
		staff.PATCH("/:id", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.StaffUpdate.UpdateStaff)
		staff.PATCH("/me", middleware.JWTAuth(*cfg, queries), handlers.StaffUpdateMe.UpdateMyStaff)
		staff.POST("/:id/store-access", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.StaffStoreAccess.CreateStoreAccess)
		staff.DELETE("/:id/store-access", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.StaffDeleteStoreAccess.DeleteStoreAccessBulk)
	}
}

func setupStylistRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	stylists := api.Group("/stylists")
	{
		stylists.POST("/me", middleware.JWTAuth(*cfg, queries), middleware.RequireRoles(staffModel.RoleAdmin, staffModel.RoleManager, staffModel.RoleStylist), handlers.StylistCreate.CreateMyStylist)
		stylists.PATCH("/me", middleware.JWTAuth(*cfg, queries), middleware.RequireRoles(staffModel.RoleAdmin, staffModel.RoleManager, staffModel.RoleStylist), handlers.StylistUpdate.UpdateMyStylist)
	}
}

func setupScheduleRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	schedules := api.Group("/schedules")
	{
		schedules.POST("/bulk", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.ScheduleCreateBulk.CreateSchedulesBulk)
		schedules.DELETE("/bulk", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.ScheduleDeleteBulk.DeleteSchedulesBulk)
		schedules.POST("/:scheduleId/time-slots", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.ScheduleCreateTimeSlot.CreateTimeSlot)
		schedules.PATCH("/:scheduleId/time-slots/:timeSlotId", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.ScheduleUpdateTimeSlot.UpdateTimeSlot)
		schedules.DELETE("/:scheduleId/time-slots/:timeSlotId", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.ScheduleDeleteTimeSlot.DeleteTimeSlot)
	}
}

func setupStoreRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	stores := api.Group("/stores")
	{
		stores.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.StoreCreate.CreateStore)
		stores.PATCH("/:storeId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.StoreUpdate.UpdateStore)
	}
}

func setupTimeSlotTemplateRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	timeSlotTemplates := api.Group("/time-slot-templates")
	{
		timeSlotTemplates.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.TimeSlotTemplateCreate.CreateTimeSlotTemplate)
		timeSlotTemplates.PATCH("/:templateId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.TimeSlotTemplateUpdate.UpdateTimeSlotTemplate)
		timeSlotTemplates.DELETE("/:templateId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.TimeSlotTemplateDelete.DeleteTimeSlotTemplate)
		timeSlotTemplates.POST("/:templateId/items", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.TimeSlotTemplateCreateItem.CreateTimeSlotTemplateItem)
		timeSlotTemplates.PATCH("/:templateId/items/:itemId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.TimeSlotTemplateUpdateItem.UpdateTimeSlotTemplateItem)
		timeSlotTemplates.DELETE("/:templateId/items/:itemId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.TimeSlotTemplateDeleteItem.DeleteTimeSlotTemplateItem)
	}
}

func setupServiceRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	services := api.Group("/services")
	{
		services.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.ServiceCreate.CreateService)
		services.PATCH("/:serviceId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.ServiceUpdate.UpdateService)
	}
}