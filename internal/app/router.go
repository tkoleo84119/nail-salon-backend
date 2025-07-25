package app

import (
	"github.com/gin-gonic/gin"
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	"github.com/tkoleo84119/nail-salon-backend/internal/handler"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
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
		// Public/Customer routes
		setupPublicAuthRoutes(api, handlers)
		setupPublicCustomerRoutes(api, cfg, queries, handlers)
		setupPublicBookingRoutes(api, cfg, queries, handlers)
		setupPublicStoreRoutes(api, cfg, queries, handlers)
		setupPublicServiceRoutes(api, cfg, queries, handlers)
		setupPublicScheduleRoutes(api, cfg, queries, handlers)

		// Admin routes
		admin := api.Group("/admin")
		{
			setupAdminAuthRoutes(admin, handlers)
			setupAdminStaffRoutes(admin, cfg, queries, handlers)
			setupAdminStylistRoutes(admin, cfg, queries, handlers)
			setupAdminStoreRoutes(admin, cfg, queries, handlers)
			setupAdminServiceRoutes(admin, cfg, queries, handlers)
			setupAdminScheduleRoutes(admin, cfg, queries, handlers)
			setupAdminTimeSlotTemplateRoutes(admin, cfg, queries, handlers)
			setupAdminBookingRoutes(admin, cfg, queries, handlers)
			setupAdminCustomerRoutes(admin, cfg, queries, handlers)
		}
	}

	return router
}

func setupPublicAuthRoutes(api *gin.RouterGroup, handlers Handlers) {
	auth := api.Group("/auth")
	{
		line := auth.Group("/line")
		{
			line.POST("/login", handlers.AuthCustomerLineLogin.CustomerLineLogin)
			line.POST("/register", handlers.AuthCustomerLineRegister.CustomerLineRegister)
		}
	}
}

func setupAdminAuthRoutes(admin *gin.RouterGroup, handlers Handlers) {
	auth := admin.Group("/auth")
	{
		auth.POST("/login", handlers.AuthStaffLogin.StaffLogin)
	}
}

func setupPublicBookingRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	bookings := api.Group("/bookings")
	{
		// Customer booking operations
		bookings.POST("", middleware.CustomerJWTAuth(*cfg, queries), handlers.BookingCreateMy.CreateMyBooking)
		bookings.PATCH("/:bookingId", middleware.CustomerJWTAuth(*cfg, queries), handlers.BookingUpdateMy.UpdateMyBooking)
		bookings.PATCH("/:bookingId/cancel", middleware.CustomerJWTAuth(*cfg, queries), handlers.BookingCancelMy.CancelMyBooking)
	}
}

func setupAdminBookingRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
}

func setupPublicCustomerRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	customers := api.Group("/customers")
	{
		// Customer self-service
		customers.GET("/me", middleware.CustomerJWTAuth(*cfg, queries), handlers.CustomerGetMy.GetMyCustomer)
		customers.PATCH("/me", middleware.CustomerJWTAuth(*cfg, queries), handlers.CustomerUpdateMy.UpdateMyCustomer)
	}
}

func setupAdminCustomerRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
}

func setupPublicScheduleRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
}

func setupAdminScheduleRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	schedules := admin.Group("/schedules")
	{
		// Bulk operations
		schedules.POST("/bulk", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.AdminScheduleCreateBulk.CreateSchedulesBulk)
		schedules.DELETE("/bulk", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.AdminScheduleDeleteBulk.DeleteSchedulesBulk)

		// Time slot operations
		schedules.POST("/:scheduleId/time-slots", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.AdminScheduleCreateTimeSlot.CreateTimeSlot)
		schedules.PATCH("/:scheduleId/time-slots/:timeSlotId", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.AdminScheduleUpdateTimeSlot.UpdateTimeSlot)
		schedules.DELETE("/:scheduleId/time-slots/:timeSlotId", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.AdminScheduleDeleteTimeSlot.DeleteTimeSlot)
	}
}

func setupPublicServiceRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
}

func setupAdminServiceRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	services := admin.Group("/services")
	{
		services.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.AdminServiceCreate.CreateService)
		services.PATCH("/:serviceId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.AdminServiceUpdate.UpdateService)
	}
}

func setupAdminStaffRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	staff := admin.Group("/staff")
	{
		// Staff management
		staff.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.AdminStaffCreate.CreateStaff)
		staff.PATCH("/:staffId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.AdminStaffUpdate.UpdateStaff)
		staff.PATCH("/me", middleware.JWTAuth(*cfg, queries), handlers.AdminStaffUpdateMe.UpdateMyStaff)

		// Store access management
		staff.POST("/:staffId/store-access", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.AdminStaffStoreAccess.CreateStoreAccess)
		staff.DELETE("/:staffId/store-access/bulk", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.AdminStaffDeleteStoreAccess.DeleteStoreAccessBulk)

	}
}

func setupPublicStoreRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
}

func setupAdminStoreRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	stores := admin.Group("/stores")
	{
		stores.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.AdminStoreCreate.CreateStore)
		stores.PATCH("/:storeId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.AdminStoreUpdate.UpdateStore)
	}
}

func setupAdminStylistRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	stylists := admin.Group("/stylists")
	{
		// Self-service stylist operations
		stylists.POST("/me", middleware.JWTAuth(*cfg, queries), middleware.RequireRoles(adminStaffModel.RoleAdmin, adminStaffModel.RoleManager, adminStaffModel.RoleStylist), handlers.AdminStylistCreate.CreateMyStylist)
		stylists.PATCH("/me", middleware.JWTAuth(*cfg, queries), middleware.RequireRoles(adminStaffModel.RoleAdmin, adminStaffModel.RoleManager, adminStaffModel.RoleStylist), handlers.AdminStylistUpdate.UpdateMyStylist)
	}
}

func setupAdminTimeSlotTemplateRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	timeSlotTemplates := admin.Group("/time-slot-templates")
	{
		// Template management
		timeSlotTemplates.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.AdminTimeSlotTemplateCreate.CreateTimeSlotTemplate)
		timeSlotTemplates.PATCH("/:templateId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.AdminTimeSlotTemplateUpdate.UpdateTimeSlotTemplate)
		timeSlotTemplates.DELETE("/:templateId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.AdminTimeSlotTemplateDelete.DeleteTimeSlotTemplate)

		// Template item management
		timeSlotTemplates.POST("/:templateId/items", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.AdminTimeSlotTemplateCreateItem.CreateTimeSlotTemplateItem)
		timeSlotTemplates.PATCH("/:templateId/items/:itemId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.AdminTimeSlotTemplateUpdateItem.UpdateTimeSlotTemplateItem)
		timeSlotTemplates.DELETE("/:templateId/items/:itemId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.AdminTimeSlotTemplateDeleteItem.DeleteTimeSlotTemplateItem)
	}
}
