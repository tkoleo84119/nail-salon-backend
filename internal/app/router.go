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

	// Apply CORS middleware globally
	router.Use(middleware.CORSMiddleware(cfg.CORS))

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
		}
	}

	return router
}

// Public route setup functions
func setupPublicAuthRoutes(api *gin.RouterGroup, handlers Handlers) {
	auth := api.Group("/auth")
	{
		line := auth.Group("/line")
		{
			line.POST("/login", handlers.Public.AuthCustomerLineLogin.CustomerLineLogin)
			line.POST("/register", handlers.Public.AuthCustomerLineRegister.CustomerLineRegister)
		}

		token := auth.Group("/token")
		{
			token.POST("/refresh", handlers.Public.AuthRefreshToken.RefreshToken)
		}
	}
}

func setupPublicCustomerRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	customers := api.Group("/customers")
	{
		// Customer self-service
		customers.GET("/me", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.CustomerGetMy.GetMyCustomer)
		customers.PATCH("/me", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.CustomerUpdateMy.UpdateMyCustomer)
	}
}

func setupPublicBookingRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	bookings := api.Group("/bookings")
	{
		// Customer booking operations
		bookings.GET("", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.BookingGetMy.GetMyBookings)
		bookings.GET("/:bookingId", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.BookingGetMySingle.GetMyBooking)
		bookings.POST("", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.BookingCreateMy.CreateMyBooking)
		bookings.PATCH("/:bookingId", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.BookingUpdateMy.UpdateMyBooking)
		bookings.PATCH("/:bookingId/cancel", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.BookingCancelMy.CancelMyBooking)
	}
}

func setupPublicStoreRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	stores := api.Group("/stores")
	{
		// Store listing
		stores.GET("", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.StoreGetStores.GetStores)

		// Single store detail
		stores.GET("/:storeId", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.StoreGetStore.GetStore)

		// Store services browsing
		stores.GET("/:storeId/services", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.StoreGetServices.GetStoreServices)

		// Store stylists browsing
		stores.GET("/:storeId/stylists", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.StoreGetStylists.GetStoreStylists)

		// Store schedule routes
		stores.GET("/:storeId/stylists/:stylistId/schedules", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.ScheduleGetStoreSchedules.GetStoreSchedules)
	}
}

func setupPublicServiceRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	// TODO: Add public service routes for browsing
	// services := api.Group("/services")
	// {
	// 	services.GET("", handlers.Public.ServiceList.List)
	// 	services.GET("/:serviceId", handlers.Public.ServiceGet.Get)
	// }
}

func setupPublicScheduleRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	schedules := api.Group("/schedules")
	{
		// Schedule time slot routes
		schedules.GET("/:scheduleId/time-slots", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.ScheduleGetTimeSlots.GetTimeSlots)
	}
}

// Admin route setup functions
func setupAdminAuthRoutes(admin *gin.RouterGroup, handlers Handlers) {
	auth := admin.Group("/auth")
	{
		auth.POST("/login", handlers.Admin.AuthStaffLogin.StaffLogin)

		token := auth.Group("/token")
		{
			token.POST("/refresh", handlers.Admin.AuthStaffRefreshToken.StaffRefreshToken)
			token.POST("/revoke", handlers.Admin.AuthStaffLogout.StaffLogout)
		}
	}
}

func setupAdminStaffRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	staff := admin.Group("/staff")
	{
		// Staff management
		staff.GET("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StaffGetList.GetStaffList)
		staff.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StaffCreate.CreateStaff)
		staff.GET("/:staffId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StaffGet.GetStaff)
		staff.PATCH("/:staffId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StaffUpdate.UpdateStaff)
		staff.GET("/me", middleware.JWTAuth(*cfg, queries), handlers.Admin.StaffGetMe.GetMyStaff)
		staff.PATCH("/me", middleware.JWTAuth(*cfg, queries), handlers.Admin.StaffUpdateMe.UpdateMyStaff)

		// Store access management
		staff.GET("/:staffId/store-access", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StaffGetStoreAccess.GetStaffStoreAccess)
		staff.POST("/:staffId/store-access", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StaffStoreAccess.CreateStoreAccess)
		staff.DELETE("/:staffId/store-access/bulk", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StaffDeleteStoreAccess.DeleteStoreAccessBulk)
	}
}

func setupAdminStylistRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	stylists := admin.Group("/stylists")
	{
		// Self-service stylist operations
		stylists.PATCH("/me", middleware.JWTAuth(*cfg, queries), middleware.RequireRoles(adminStaffModel.RoleAdmin, adminStaffModel.RoleManager, adminStaffModel.RoleStylist), handlers.Admin.StylistUpdate.UpdateMyStylist)
	}
}

func setupAdminStoreRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	stores := admin.Group("/stores")
	{
		stores.GET("", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.StoreGetList.GetStoreList)
		stores.GET("/:storeId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StoreGet.GetStore)
		stores.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StoreCreate.CreateStore)
		stores.PATCH("/:storeId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StoreUpdate.UpdateStore)

		// Store stylists routes
		stores.GET("/:storeId/stylists", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.StylistGetList.GetStylistList)

		// Store schedules routes
		stores.GET("/:storeId/schedules", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.ScheduleGetAll.GetAll)
		stores.GET("/:storeId/schedules/:scheduleId", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.ScheduleGet.GetSchedule)
		stores.POST("/:storeId/schedules/bulk", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.ScheduleCreateBulk.CreateBulk)
		stores.DELETE("/:storeId/schedules/bulk", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.ScheduleDeleteBulk.DeleteBulk)

		// Store bookings routes
		stores.GET("/:storeId/bookings", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.BookingGetList.GetBookingList)
		stores.POST("/:storeId/bookings", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.BookingCreate.CreateBooking)
		stores.PATCH("/:storeId/bookings/:bookingId", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.BookingUpdateByStaff.UpdateBookingByStaff)
		stores.PATCH("/:storeId/bookings/:bookingId/cancel", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.BookingCancel.CancelBooking)
	}
}

func setupAdminServiceRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	services := admin.Group("/services")
	{
		services.GET("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.ServiceGetList.GetServiceList)
		services.GET("/:serviceId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.ServiceGet.GetService)
		services.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.ServiceCreate.CreateService)
		services.PATCH("/:serviceId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.ServiceUpdate.UpdateService)
	}
}

func setupAdminScheduleRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	schedules := admin.Group("/schedules")
	{
		// Time slot operations
		schedules.POST("/:scheduleId/time-slots", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.ScheduleCreateTimeSlot.Create)
		schedules.PATCH("/:scheduleId/time-slots/:timeSlotId", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.ScheduleUpdateTimeSlot.Update)
		schedules.DELETE("/:scheduleId/time-slots/:timeSlotId", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.ScheduleDeleteTimeSlot.Delete)
	}
}

func setupAdminTimeSlotTemplateRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	timeSlotTemplates := admin.Group("/time-slot-templates")
	{
		// Template management
		timeSlotTemplates.GET("", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.TimeSlotTemplateGetAll.GetAll)
		timeSlotTemplates.GET("/:templateId", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.TimeSlotTemplateGet.Get)
		timeSlotTemplates.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.TimeSlotTemplateCreate.Create)
		timeSlotTemplates.PATCH("/:templateId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.TimeSlotTemplateUpdate.Update)
		timeSlotTemplates.DELETE("/:templateId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.TimeSlotTemplateDelete.Delete)

		// Template item management
		timeSlotTemplates.POST("/:templateId/items", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.TimeSlotTemplateItemCreate.Create)
		timeSlotTemplates.PATCH("/:templateId/items/:itemId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.TimeSlotTemplateUpdateItem.Update)
		timeSlotTemplates.DELETE("/:templateId/items/:itemId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.TimeSlotTemplateDeleteItem.Delete)
	}
}
