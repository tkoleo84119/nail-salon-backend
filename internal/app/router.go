package app

import (
	"github.com/gin-gonic/gin"
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	"github.com/tkoleo84119/nail-salon-backend/internal/handler"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
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
			setupAdminAuthRoutes(admin, cfg, queries, handlers)
			setupAdminStaffRoutes(admin, cfg, queries, handlers)
			setupAdminStylistRoutes(admin, cfg, queries, handlers)
			setupAdminCustomerRoutes(admin, cfg, queries, handlers)
			setupAdminStoreRoutes(admin, cfg, queries, handlers)
			setupAdminAccountRoutes(admin, cfg, queries, handlers)
			setupAdminBrandRoutes(admin, cfg, queries, handlers)
			setupAdminSupplierRoutes(admin, cfg, queries, handlers)
			setupAdminExpenseRoutes(admin, cfg, queries, handlers)
			setupAdminProductCategoryRoutes(admin, cfg, queries, handlers)
			setupAdminServiceRoutes(admin, cfg, queries, handlers)
			setupAdminScheduleRoutes(admin, cfg, queries, handlers)
			setupAdminTimeSlotTemplateRoutes(admin, cfg, queries, handlers)
			setupAdminCouponRoutes(admin, cfg, queries, handlers)
			setupAdminCustomerCouponRoutes(admin, cfg, queries, handlers)
			setupAdminReportRoutes(admin, cfg, queries, handlers)
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
			line.POST("/login", handlers.Public.AuthLineLogin.LineLogin)
			line.POST("/register", handlers.Public.AuthLineRegister.LineRegister)
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
		customers.GET("/me", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.CustomerGetMe.GetMe)
		customers.PATCH("/me", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.CustomerUpdateMe.UpdateMe)
	}

	// Customer coupons
	customerCoupons := api.Group("/customer_coupons")
	{
		customerCoupons.GET("", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.CustomerCouponGetAll.GetAll)
	}
}

func setupPublicBookingRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	bookings := api.Group("/bookings")
	{
		// Customer booking operations
		bookings.GET("", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.BookingGetAll.GetAll)
		bookings.GET("/:bookingId", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.BookingGetMySingle.Get)
		bookings.POST("", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.BookingCreate.Create)
		bookings.PATCH("/:bookingId", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.BookingUpdate.Update)
		bookings.PATCH("/:bookingId/cancel", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.BookingCancel.Cancel)
	}
}

func setupPublicStoreRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	stores := api.Group("/stores")
	{
		// Store listing
		stores.GET("", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.StoreGetAll.GetAll)

		// Store stylists browsing
		stores.GET("/:storeId/stylists", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.StylistGetAll.GetAll)

		// Store stylist schedule routes
		stores.GET("/:storeId/stylists/:stylistId/schedules", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.ScheduleGetAll.GetAll)
	}
}

func setupPublicServiceRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	services := api.Group("/services")
	{
		services.GET("", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.ServiceGetAll.GetAll)
	}
}

func setupPublicScheduleRoutes(api *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	schedules := api.Group("/schedules")
	{
		// Schedule time slot routes
		schedules.GET("/:scheduleId/time-slots", middleware.CustomerJWTAuth(*cfg, queries), handlers.Public.TimeSlotGetAll.GetAll)
	}
}

// Admin route setup functions
func setupAdminAuthRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	auth := admin.Group("/auth")
	{
		auth.POST("/login", handlers.Admin.AuthStaffLogin.Login)
		auth.GET("/permission", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.AuthStaffPermission.GetPermission)

		token := auth.Group("/token")
		{
			token.POST("/refresh", handlers.Admin.AuthStaffRefreshToken.RefreshToken)
			token.POST("/revoke", handlers.Admin.AuthStaffLogout.Logout)
		}
	}
}

func setupAdminStaffRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	staff := admin.Group("/staff")
	{
		// Staff management
		staff.GET("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StaffGetAll.GetAll)
		staff.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StaffCreate.Create)
		staff.GET("/:staffId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StaffGet.Get)
		staff.PATCH("/:staffId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StaffUpdate.Update)
		staff.GET("/me", middleware.JWTAuth(*cfg, queries), handlers.Admin.StaffGetMe.GetMe)
		staff.PATCH("/me", middleware.JWTAuth(*cfg, queries), handlers.Admin.StaffUpdateMe.UpdateMe)

		// Store access management
		staff.GET("/:staffId/store-access", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StaffGetStoreAccess.Get)
		staff.POST("/:staffId/store-access", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StaffCreateStoreAccess.Create)
		staff.DELETE("/:staffId/store-access/bulk", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StaffDeleteBulkStoreAccess.DeleteBulk)
	}
}

func setupAdminStylistRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	stylists := admin.Group("/stylists")
	{
		// Self-service stylist operations
		stylists.PATCH("/me", middleware.JWTAuth(*cfg, queries), middleware.RequireRoles(common.RoleAdmin, common.RoleManager, common.RoleStylist), handlers.Admin.StylistUpdateMe.UpdateMe)
	}
}

func setupAdminStoreRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	stores := admin.Group("/stores")
	{
		stores.GET("", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.StoreGetList.GetAll)
		stores.GET("/:storeId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StoreGet.Get)
		stores.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StoreCreate.Create)
		stores.PATCH("/:storeId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.StoreUpdate.Update)

		// Store stylists routes
		stores.GET("/:storeId/stylists", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.StylistGetAll.GetAll)

		// Store staff username routes
		stores.GET("/:storeId/staff/store-username", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.StaffGetStoreUsername.GetStoreUsername)

		// Store schedules routes
		stores.GET("/:storeId/schedules", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.ScheduleGetAll.GetAll)
		stores.GET("/:storeId/schedules/:scheduleId", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.ScheduleGet.Get)
		stores.PATCH("/:storeId/schedules/:scheduleId", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.ScheduleUpdate.Update)
		stores.POST("/:storeId/schedules/bulk", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.ScheduleCreateBulk.CreateBulk)
		stores.DELETE("/:storeId/schedules/bulk", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.ScheduleDeleteBulk.DeleteBulk)

		// Store bookings routes
		stores.GET("/:storeId/bookings", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.BookingGetAll.GetAll)
		stores.GET("/:storeId/bookings/:bookingId", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.BookingGet.Get)
		stores.POST("/:storeId/bookings", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.BookingCreate.Create)
		stores.PATCH("/:storeId/bookings/:bookingId", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.BookingUpdate.Update)
		stores.PATCH("/:storeId/bookings/:bookingId/cancel", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.BookingCancel.Cancel)
		stores.PATCH("/:storeId/bookings/:bookingId/completed", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.BookingUpdateCompleted.UpdateCompleted)

		// Store checkouts routes
		stores.POST("/:storeId/bookings/:bookingId/checkouts", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.CheckoutCreate.Create)

		// Store booking products routes
		stores.GET("/:storeId/bookings/:bookingId/products", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.BookingProductGetAll.GetAll)
		stores.POST("/:storeId/bookings/:bookingId/products/bulk", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.BookingProductBulkCreate.BulkCreate)
		stores.DELETE("/:storeId/bookings/:bookingId/products/bulk", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.BookingProductBulkDelete.BulkDelete)

		// Store products routes
		stores.GET("/:storeId/products", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.ProductGetAll.GetAll)
		stores.POST("/:storeId/products", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.ProductCreate.Create)
		stores.PATCH("/:storeId/products/:productId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.ProductUpdate.Update)

		// Store stock usages routes
		stores.GET("/:storeId/stock-usages", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.StockUsagesGetAll.GetAll)
		stores.POST("/:storeId/stock-usages", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.StockUsagesCreate.Create)
		stores.PATCH("/:storeId/stock-usages/:stockUsageId/finish", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.StockUsagesUpdateFinish.UpdateFinish)

		// Store accounts routes
		stores.GET("/:storeId/accounts", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.AccountGetAll.GetAll)

		// Store account transactions routes
		stores.GET("/:storeId/accounts/:accountId/transactions", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.AccountTransactionGetAll.GetAll)
	}
}

func setupAdminServiceRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	services := admin.Group("/services")
	{
		services.GET("", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.ServiceGetList.GetAll)
		services.GET("/:serviceId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.ServiceGet.Get)
		services.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.ServiceCreate.Create)
		services.PATCH("/:serviceId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.ServiceUpdate.Update)
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

func setupAdminCustomerRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	customers := admin.Group("/customers")
	{
		// Customer management - all staff can view customers
		customers.GET("", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.CustomerGetAll.GetAll)
		customers.GET("/:customerId", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.CustomerGet.Get)
		customers.PATCH("/:customerId", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.CustomerUpdate.Update)
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

func setupAdminCouponRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	coupons := admin.Group("/coupons")
	{
		coupons.GET("", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.CouponGetAll.GetAll)
		coupons.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.CouponCreate.Create)
		coupons.PATCH("/:couponId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.CouponUpdate.Update)
	}
}

func setupAdminCustomerCouponRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	customerCoupons := admin.Group("/customer_coupons")
	{
		customerCoupons.GET("", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.CustomerCouponGetAll.GetAll)
		customerCoupons.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.CustomerCouponCreate.Create)
	}
}

func setupAdminBrandRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	brands := admin.Group("/brands")
	{
		brands.GET("", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.BrandGetAll.GetAll)
		brands.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.BrandCreate.Create)
		brands.PATCH("/:brandId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.BrandUpdate.Update)
	}
}

func setupAdminSupplierRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	suppliers := admin.Group("/suppliers")
	{
		suppliers.GET("", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.SupplierGetAll.GetAll)
		suppliers.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.SupplierCreate.Create)
		suppliers.PATCH("/:supplierId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.SupplierUpdate.Update)
	}
}

func setupAdminExpenseRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	stores := admin.Group("/stores")
	{
		stores.POST("/:storeId/expenses", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.ExpenseCreate.Create)
		stores.GET("/:storeId/expenses", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.ExpenseGetAll.GetAll)
		stores.GET("/:storeId/expenses/:expenseId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.ExpenseGet.Get)
		stores.PATCH("/:storeId/expenses/:expenseId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.ExpenseUpdate.Update)

		// Expense items routes
		stores.POST("/:storeId/expenses/:expenseId/items", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.ExpenseItemCreate.Create)
		stores.PATCH("/:storeId/expenses/:expenseId/items/:expenseItemId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.ExpenseItemUpdate.Update)
	}
}

func setupAdminProductCategoryRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	productCategories := admin.Group("/product-categories")
	{
		productCategories.GET("", middleware.JWTAuth(*cfg, queries), middleware.RequireAnyStaffRole(), handlers.Admin.ProductCategoryGetAll.GetAll)
		productCategories.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.ProductCategoryCreate.Create)
		productCategories.PATCH("/:productCategoryId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.ProductCategoryUpdate.Update)
	}
}

func setupAdminAccountRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	accounts := admin.Group("/accounts")
	{
		accounts.POST("", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.AccountCreate.Create)
		accounts.PATCH("/:accountId", middleware.JWTAuth(*cfg, queries), middleware.RequireAdminRoles(), handlers.Admin.AccountUpdate.Update)
	}
}

func setupAdminReportRoutes(admin *gin.RouterGroup, cfg *config.Config, queries *dbgen.Queries, handlers Handlers) {
	reports := admin.Group("/reports")
	{
		// Performance report - all staff except SUPER_ADMIN can access
		reports.GET("/performance/me", middleware.JWTAuth(*cfg, queries), middleware.RequireNotSuperAdmin(), handlers.Admin.ReportGetPerformanceMe.GetPerformanceMe)
		// Store performance report - SUPER_ADMIN, ADMIN, and MANAGER can access
		reports.GET("/performance/store/:storeId", middleware.JWTAuth(*cfg, queries), middleware.RequireManagerOrAbove(), handlers.Admin.ReportGetStorePerformance.GetStorePerformance)
	}
}
