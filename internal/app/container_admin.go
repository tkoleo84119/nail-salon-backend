package app

import (
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	"github.com/tkoleo84119/nail-salon-backend/internal/infra/db"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"

	// Admin handlers
	adminAuthHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/auth"
	adminBookingHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/booking"
	adminBookingProductHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/booking_product"
	adminBrandHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/brand"
	adminCheckoutHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/checkout"
	adminCouponHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/coupon"
	adminCustomerHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/customer"
	adminCustomerCouponHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/customer_coupon"
	adminProductHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/product"
	adminProductCategoryHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/product_category"
	adminReportHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/report"
	adminScheduleHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/schedule"
	adminServiceHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/service"
	adminStaffHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/staff"
	adminStoreHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/store"
	adminStoreAccessHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/store_access"
	adminStylistHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/stylist"
	adminTimeSlotTemplateHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/time-slot-template"
	adminTimeSlotHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/time_slot"
	adminTimeSlotTemplateItemHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/time_slot_template_item"

	// Admin services
	adminAuthService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/auth"
	adminBookingService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/booking"
	adminBookingProductService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/booking_product"
	adminBrandService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/brand"
	adminCheckoutService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/checkout"
	adminCouponService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/coupon"
	adminCustomerService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/customer"
	adminCustomerCouponService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/customer_coupon"
	adminProductService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/product"
	adminProductCategoryService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/product_category"
	adminReportService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/report"
	adminScheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/schedule"
	adminServiceService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/service"
	adminStaffService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/staff"
	adminStoreService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/store"
	adminStoreAccessService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/store_access"
	adminStylistService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/stylist"
	adminTimeSlotTemplateService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/time-slot-template"
	adminTimeSlotService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/time_slot"
	adminTimeSlotTemplateItemService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/time_slot_template_item"
)

// AdminServices contains all admin-facing services
type AdminServices struct {
	// Authentication services
	AuthStaffLogin        adminAuthService.LoginInterface
	AuthStaffRefreshToken adminAuthService.RefreshTokenInterface
	AuthStaffLogout       adminAuthService.LogoutInterface
	AuthStaffPermission   adminAuthService.PermissionInterface

	// Staff management services
	StaffCreate   adminStaffService.CreateInterface
	StaffGetAll   adminStaffService.GetAllInterface
	StaffUpdate   adminStaffService.UpdateInterface
	StaffUpdateMe adminStaffService.UpdateMeInterface
	StaffGet      adminStaffService.GetInterface
	StaffGetMe    adminStaffService.GetMeInterface

	// Store access services
	StaffGetStoreAccess        adminStoreAccessService.GetInterface
	StaffCreateStoreAccess     adminStoreAccessService.CreateInterface
	StaffDeleteBulkStoreAccess adminStoreAccessService.DeleteBulkInterface

	// Store management services
	StoreGetList adminStoreService.GetAllInterface
	StoreGet     adminStoreService.GetInterface
	StoreCreate  adminStoreService.CreateInterface
	StoreUpdate  adminStoreService.UpdateInterface

	// Brand management services
	BrandCreate adminBrandService.CreateInterface
	BrandGetAll adminBrandService.GetAllInterface
	BrandUpdate adminBrandService.UpdateInterface

	// Product management services
	ProductCreate adminProductService.CreateInterface
	ProductGetAll adminProductService.GetAllInterface
	ProductUpdate adminProductService.UpdateInterface

	// Product category management services
	ProductCategoryCreate adminProductCategoryService.CreateInterface
	ProductCategoryGetAll adminProductCategoryService.GetAllInterface
	ProductCategoryUpdate adminProductCategoryService.UpdateInterface

	// Service management services
	ServiceGetList adminServiceService.GetAllInterface
	ServiceGet     adminServiceService.GetInterface
	ServiceCreate  adminServiceService.CreateInterface
	ServiceUpdate  adminServiceService.UpdateInterface

	// Stylist management services
	StylistUpdateMe adminStylistService.UpdateMeInterface
	StylistGetAll   adminStylistService.GetAllInterface

	// Customer management services
	CustomerGetAll adminCustomerService.GetAllInterface
	CustomerGet    adminCustomerService.GetInterface
	CustomerUpdate adminCustomerService.UpdateInterface

	// Booking management services
	BookingCreate          adminBookingService.CreateInterface
	BookingGetAll          adminBookingService.GetAllInterface
	BookingUpdate          adminBookingService.UpdateInterface
	BookingCancel          adminBookingService.CancelInterface
	BookingGet             adminBookingService.GetInterface
	BookingUpdateCompleted adminBookingService.UpdateCompletedInterface

	// Booking product services
	BookingProductBulkCreate adminBookingProductService.BulkCreateInterface
	BookingProductBulkDelete adminBookingProductService.BulkDeleteInterface
	BookingProductGetAll     adminBookingProductService.GetAllInterface

	// Schedule management services
	ScheduleCreateBulk     adminScheduleService.CreateBulkInterface
	ScheduleDeleteBulk     adminScheduleService.DeleteBulkInterface
	ScheduleUpdate         adminScheduleService.UpdateInterface
	ScheduleCreateTimeSlot adminTimeSlotService.CreateInterface
	ScheduleUpdateTimeSlot adminTimeSlotService.UpdateInterface
	ScheduleDeleteTimeSlot adminTimeSlotService.DeleteInterface
	ScheduleGetAll         adminScheduleService.GetAllInterface
	ScheduleGet            adminScheduleService.GetInterface

	// Time slot template services
	TimeSlotTemplateGetAll adminTimeSlotTemplateService.GetAllInterface
	TimeSlotTemplateGet    adminTimeSlotTemplateService.GetInterface
	TimeSlotTemplateCreate adminTimeSlotTemplateService.CreateInterface
	TimeSlotTemplateDelete adminTimeSlotTemplateService.DeleteInterface
	TimeSlotTemplateUpdate adminTimeSlotTemplateService.UpdateInterface

	// Time slot template item services
	TimeSlotTemplateItemCreate adminTimeSlotTemplateItemService.CreateInterface
	TimeSlotTemplateUpdateItem adminTimeSlotTemplateItemService.UpdateInterface
	TimeSlotTemplateDeleteItem adminTimeSlotTemplateItemService.DeleteInterface

	// Coupon management services
	CouponCreate adminCouponService.CreateInterface
	CouponGetAll adminCouponService.GetAllInterface
	CouponUpdate adminCouponService.UpdateInterface

	// Customer coupon services
	CustomerCouponGetAll adminCustomerCouponService.GetAllInterface
	CustomerCouponCreate adminCustomerCouponService.CreateInterface

	// Checkout services
	CheckoutCreate adminCheckoutService.CreateInterface

	// Report services
	ReportGetPerformanceMe    adminReportService.GetPerformanceMeInterface
	ReportGetStorePerformance adminReportService.GetStorePerformanceInterface
}

// AdminHandlers contains all admin-facing handlers
type AdminHandlers struct {
	// Authentication handlers
	AuthStaffLogin        *adminAuthHandler.Login
	AuthStaffRefreshToken *adminAuthHandler.RefreshToken
	AuthStaffLogout       *adminAuthHandler.Logout
	AuthStaffPermission   *adminAuthHandler.Permission

	// Staff management handlers
	StaffCreate   *adminStaffHandler.Create
	StaffUpdate   *adminStaffHandler.Update
	StaffUpdateMe *adminStaffHandler.UpdateMe
	StaffGet      *adminStaffHandler.Get
	StaffGetMe    *adminStaffHandler.GetMe
	StaffGetAll   *adminStaffHandler.GetAll

	// Store access handlers
	StaffGetStoreAccess        *adminStoreAccessHandler.Get
	StaffCreateStoreAccess     *adminStoreAccessHandler.Create
	StaffDeleteBulkStoreAccess *adminStoreAccessHandler.DeleteBulk

	// Store management handlers
	StoreGetList *adminStoreHandler.GetAll
	StoreGet     *adminStoreHandler.Get
	StoreCreate  *adminStoreHandler.Create
	StoreUpdate  *adminStoreHandler.Update

	// Brand management handlers
	BrandCreate *adminBrandHandler.Create
	BrandGetAll *adminBrandHandler.GetAll
	BrandUpdate *adminBrandHandler.Update

	// Product management handlers
	ProductCreate *adminProductHandler.Create
	ProductGetAll *adminProductHandler.GetAll
	ProductUpdate *adminProductHandler.Update

	// Product category management handlers
	ProductCategoryCreate *adminProductCategoryHandler.Create
	ProductCategoryGetAll *adminProductCategoryHandler.GetAll
	ProductCategoryUpdate *adminProductCategoryHandler.Update

	// Service management handlers
	ServiceGetList *adminServiceHandler.GetAll
	ServiceGet     *adminServiceHandler.Get
	ServiceCreate  *adminServiceHandler.Create
	ServiceUpdate  *adminServiceHandler.Update

	// Stylist management handlers
	StylistUpdateMe *adminStylistHandler.UpdateMe
	StylistGetAll   *adminStylistHandler.GetAll

	// Customer management handlers
	CustomerGetAll *adminCustomerHandler.GetAll
	CustomerGet    *adminCustomerHandler.Get
	CustomerUpdate *adminCustomerHandler.Update

	// Booking management handlers
	BookingCreate          *adminBookingHandler.Create
	BookingGetAll          *adminBookingHandler.GetAll
	BookingUpdate          *adminBookingHandler.Update
	BookingCancel          *adminBookingHandler.Cancel
	BookingGet             *adminBookingHandler.Get
	BookingUpdateCompleted *adminBookingHandler.UpdateCompleted

	// Booking product handlers
	BookingProductBulkCreate *adminBookingProductHandler.BulkCreate
	BookingProductBulkDelete *adminBookingProductHandler.BulkDelete
	BookingProductGetAll     *adminBookingProductHandler.GetAll

	// Schedule management handlers
	ScheduleCreateBulk     *adminScheduleHandler.CreateBulk
	ScheduleDeleteBulk     *adminScheduleHandler.DeleteBulk
	ScheduleUpdate         *adminScheduleHandler.Update
	ScheduleCreateTimeSlot *adminTimeSlotHandler.Create
	ScheduleUpdateTimeSlot *adminTimeSlotHandler.Update
	ScheduleDeleteTimeSlot *adminTimeSlotHandler.Delete
	ScheduleGetAll         *adminScheduleHandler.GetAll
	ScheduleGet            *adminScheduleHandler.Get

	// Time slot template handlers
	TimeSlotTemplateGetAll *adminTimeSlotTemplateHandler.GetAll
	TimeSlotTemplateGet    *adminTimeSlotTemplateHandler.Get
	TimeSlotTemplateCreate *adminTimeSlotTemplateHandler.Create
	TimeSlotTemplateUpdate *adminTimeSlotTemplateHandler.Update
	TimeSlotTemplateDelete *adminTimeSlotTemplateHandler.Delete

	// Time slot template item handlers
	TimeSlotTemplateItemCreate *adminTimeSlotTemplateItemHandler.Create
	TimeSlotTemplateUpdateItem *adminTimeSlotTemplateItemHandler.Update
	TimeSlotTemplateDeleteItem *adminTimeSlotTemplateItemHandler.Delete

	// Coupon management handlers
	CouponCreate *adminCouponHandler.Create
	CouponGetAll *adminCouponHandler.GetAll
	CouponUpdate *adminCouponHandler.Update

	// Customer coupon handlers
	CustomerCouponGetAll *adminCustomerCouponHandler.GetAll
	CustomerCouponCreate *adminCustomerCouponHandler.Create

	// Checkout handlers
	CheckoutCreate *adminCheckoutHandler.Create

	// Report handlers
	ReportGetPerformanceMe    *adminReportHandler.GetPerformanceMe
	ReportGetStorePerformance *adminReportHandler.GetStorePerformance
}

// NewAdminServices creates and initializes all admin services
func NewAdminServices(queries *dbgen.Queries, database *db.Database, repositories Repositories, cfg *config.Config) AdminServices {
	return AdminServices{
		// Authentication services
		AuthStaffLogin:        adminAuthService.NewLogin(queries, cfg.JWT),
		AuthStaffRefreshToken: adminAuthService.NewRefreshToken(queries, cfg.JWT),
		AuthStaffLogout:       adminAuthService.NewLogout(queries),
		AuthStaffPermission:   adminAuthService.NewPermission(),

		// Staff management services
		StaffCreate:                adminStaffService.NewCreate(queries, database.PgxPool),
		StaffUpdate:                adminStaffService.NewUpdate(queries, repositories.SQLX),
		StaffUpdateMe:              adminStaffService.NewUpdateMe(queries, repositories.SQLX),
		StaffGet:                   adminStaffService.NewGet(queries),
		StaffGetMe:                 adminStaffService.NewGetMe(queries),
		StaffGetAll:                adminStaffService.NewGetAll(repositories.SQLX),
		StaffGetStoreAccess:        adminStoreAccessService.NewGet(queries),
		StaffCreateStoreAccess:     adminStoreAccessService.NewCreate(queries),
		StaffDeleteBulkStoreAccess: adminStoreAccessService.NewDeleteBulk(queries),

		// Store management services
		StoreGetList: adminStoreService.NewGetAll(repositories.SQLX),
		StoreGet:     adminStoreService.NewGet(queries),
		StoreCreate:  adminStoreService.NewCreate(queries, database.PgxPool),
		StoreUpdate:  adminStoreService.NewUpdate(queries, repositories.SQLX),

		// Brand management services
		BrandCreate: adminBrandService.NewCreate(queries),
		BrandGetAll: adminBrandService.NewGetAll(repositories.SQLX),
		BrandUpdate: adminBrandService.NewUpdate(queries, repositories.SQLX),

		// Product management services
		ProductCreate: adminProductService.NewCreate(queries),
		ProductGetAll: adminProductService.NewGetAll(repositories.SQLX),
		ProductUpdate: adminProductService.NewUpdate(queries, repositories.SQLX),

		// Product category management services
		ProductCategoryCreate: adminProductCategoryService.NewCreate(queries),
		ProductCategoryGetAll: adminProductCategoryService.NewGetAll(repositories.SQLX),
		ProductCategoryUpdate: adminProductCategoryService.NewUpdate(queries, repositories.SQLX),

		// Service management services
		ServiceGetList: adminServiceService.NewGetAll(repositories.SQLX),
		ServiceGet:     adminServiceService.NewGet(queries),
		ServiceCreate:  adminServiceService.NewCreate(queries),
		ServiceUpdate:  adminServiceService.NewUpdate(queries, repositories.SQLX),

		// Stylist management services
		StylistUpdateMe: adminStylistService.NewUpdateMe(queries, repositories.SQLX),
		StylistGetAll:   adminStylistService.NewGetAll(repositories.SQLX),

		// Customer management services
		CustomerGetAll: adminCustomerService.NewGetAll(repositories.SQLX),
		CustomerGet:    adminCustomerService.NewGet(queries),
		CustomerUpdate: adminCustomerService.NewUpdate(queries, repositories.SQLX),
		// Booking management services
		BookingCreate:          adminBookingService.NewCreate(queries, database.PgxPool),
		BookingGetAll:          adminBookingService.NewGetAll(queries, repositories.SQLX),
		BookingUpdate:          adminBookingService.NewUpdate(queries, repositories.SQLX, database.Sqlx),
		BookingCancel:          adminBookingService.NewCancel(queries, database.Sqlx, repositories.SQLX),
		BookingGet:             adminBookingService.NewGet(queries),
		BookingUpdateCompleted: adminBookingService.NewUpdateCompleted(queries),

		// Booking product services
		BookingProductBulkCreate: adminBookingProductService.NewBulkCreate(queries),
		BookingProductBulkDelete: adminBookingProductService.NewBulkDelete(queries),
		BookingProductGetAll:     adminBookingProductService.NewGetAll(queries, repositories.SQLX),

		// Schedule management services
		ScheduleCreateBulk:     adminScheduleService.NewCreateBulk(queries, database.PgxPool),
		ScheduleDeleteBulk:     adminScheduleService.NewDeleteBulk(queries),
		ScheduleUpdate:         adminScheduleService.NewUpdate(queries, repositories.SQLX),
		ScheduleCreateTimeSlot: adminTimeSlotService.NewCreate(queries),
		ScheduleUpdateTimeSlot: adminTimeSlotService.NewUpdate(queries, repositories.SQLX),
		ScheduleDeleteTimeSlot: adminTimeSlotService.NewDelete(queries),
		ScheduleGetAll:         adminScheduleService.NewGetAll(queries, repositories.SQLX),
		ScheduleGet:            adminScheduleService.NewGet(queries),

		// Time slot template services
		TimeSlotTemplateGetAll: adminTimeSlotTemplateService.NewGetAll(repositories.SQLX),
		TimeSlotTemplateGet:    adminTimeSlotTemplateService.NewGet(queries),
		TimeSlotTemplateCreate: adminTimeSlotTemplateService.NewCreate(queries, database.PgxPool),
		TimeSlotTemplateUpdate: adminTimeSlotTemplateService.NewUpdate(queries, repositories.SQLX),
		TimeSlotTemplateDelete: adminTimeSlotTemplateService.NewDelete(queries),

		// Time slot template item services
		TimeSlotTemplateItemCreate: adminTimeSlotTemplateItemService.NewCreate(queries),
		TimeSlotTemplateUpdateItem: adminTimeSlotTemplateItemService.NewUpdate(queries),
		TimeSlotTemplateDeleteItem: adminTimeSlotTemplateItemService.NewDelete(queries),

		// Coupon management services
		CouponCreate: adminCouponService.NewCreate(queries, repositories.SQLX),
		CouponGetAll: adminCouponService.NewGetAll(repositories.SQLX),
		CouponUpdate: adminCouponService.NewUpdate(queries, repositories.SQLX),

		// Customer coupon services
		CustomerCouponGetAll: adminCustomerCouponService.NewGetAll(queries, repositories.SQLX),
		CustomerCouponCreate: adminCustomerCouponService.NewCreate(queries),

		// Checkout services
		CheckoutCreate: adminCheckoutService.NewCreate(queries, repositories.SQLX, database.PgxPool),

		// Report services
		ReportGetPerformanceMe:    adminReportService.NewGetPerformanceMe(queries),
		ReportGetStorePerformance: adminReportService.NewGetStorePerformance(queries),
	}
}

// NewAdminHandlers creates and initializes all admin handlers
func NewAdminHandlers(services AdminServices) AdminHandlers {
	return AdminHandlers{
		// Authentication handlers
		AuthStaffLogin:        adminAuthHandler.NewLogin(services.AuthStaffLogin),
		AuthStaffRefreshToken: adminAuthHandler.NewRefreshToken(services.AuthStaffRefreshToken),
		AuthStaffLogout:       adminAuthHandler.NewLogout(services.AuthStaffLogout),
		AuthStaffPermission:   adminAuthHandler.NewPermission(services.AuthStaffPermission),

		// Staff management handlers
		StaffCreate:                adminStaffHandler.NewCreate(services.StaffCreate),
		StaffUpdate:                adminStaffHandler.NewUpdate(services.StaffUpdate),
		StaffUpdateMe:              adminStaffHandler.NewUpdateMe(services.StaffUpdateMe),
		StaffGet:                   adminStaffHandler.NewGet(services.StaffGet),
		StaffGetMe:                 adminStaffHandler.NewGetMe(services.StaffGetMe),
		StaffGetAll:                adminStaffHandler.NewGetAll(services.StaffGetAll),
		StaffGetStoreAccess:        adminStoreAccessHandler.NewGet(services.StaffGetStoreAccess),
		StaffCreateStoreAccess:     adminStoreAccessHandler.NewCreate(services.StaffCreateStoreAccess),
		StaffDeleteBulkStoreAccess: adminStoreAccessHandler.NewDeleteBulk(services.StaffDeleteBulkStoreAccess),

		// Store management handlers
		StoreGetList: adminStoreHandler.NewGetAll(services.StoreGetList),
		StoreGet:     adminStoreHandler.NewGet(services.StoreGet),
		StoreCreate:  adminStoreHandler.NewCreate(services.StoreCreate),
		StoreUpdate:  adminStoreHandler.NewUpdate(services.StoreUpdate),

		// Brand management handlers
		BrandCreate: adminBrandHandler.NewCreate(services.BrandCreate),
		BrandGetAll: adminBrandHandler.NewGetAll(services.BrandGetAll),
		BrandUpdate: adminBrandHandler.NewUpdate(services.BrandUpdate),

		// Product management handlers
		ProductCreate: adminProductHandler.NewCreate(services.ProductCreate),
		ProductGetAll: adminProductHandler.NewGetAll(services.ProductGetAll),
		ProductUpdate: adminProductHandler.NewUpdate(services.ProductUpdate),

		// Product category management handlers
		ProductCategoryCreate: adminProductCategoryHandler.NewCreate(services.ProductCategoryCreate),
		ProductCategoryGetAll: adminProductCategoryHandler.NewGetAll(services.ProductCategoryGetAll),
		ProductCategoryUpdate: adminProductCategoryHandler.NewUpdate(services.ProductCategoryUpdate),

		// Service management handlers
		ServiceGetList: adminServiceHandler.NewGetAll(services.ServiceGetList),
		ServiceGet:     adminServiceHandler.NewGet(services.ServiceGet),
		ServiceCreate:  adminServiceHandler.NewCreate(services.ServiceCreate),
		ServiceUpdate:  adminServiceHandler.NewUpdate(services.ServiceUpdate),

		// Stylist management handlers
		StylistUpdateMe: adminStylistHandler.NewUpdateMe(services.StylistUpdateMe),
		StylistGetAll:   adminStylistHandler.NewGetAll(services.StylistGetAll),

		// Customer management handlers
		CustomerGetAll: adminCustomerHandler.NewGetAll(services.CustomerGetAll),
		CustomerGet:    adminCustomerHandler.NewGet(services.CustomerGet),
		CustomerUpdate: adminCustomerHandler.NewUpdate(services.CustomerUpdate),

		// Booking management handlers
		BookingCreate:          adminBookingHandler.NewCreate(services.BookingCreate),
		BookingGetAll:          adminBookingHandler.NewGetAll(services.BookingGetAll),
		BookingUpdate:          adminBookingHandler.NewUpdate(services.BookingUpdate),
		BookingCancel:          adminBookingHandler.NewCancel(services.BookingCancel),
		BookingGet:             adminBookingHandler.NewGet(services.BookingGet),
		BookingUpdateCompleted: adminBookingHandler.NewUpdateCompleted(services.BookingUpdateCompleted),

		// Booking product handlers
		BookingProductBulkCreate: adminBookingProductHandler.NewBulkCreate(services.BookingProductBulkCreate),
		BookingProductBulkDelete: adminBookingProductHandler.NewBulkDelete(services.BookingProductBulkDelete),
		BookingProductGetAll:     adminBookingProductHandler.NewGetAll(services.BookingProductGetAll),

		// Schedule management handlers
		ScheduleCreateBulk:     adminScheduleHandler.NewCreateBulk(services.ScheduleCreateBulk),
		ScheduleDeleteBulk:     adminScheduleHandler.NewDeleteBulk(services.ScheduleDeleteBulk),
		ScheduleUpdate:         adminScheduleHandler.NewUpdate(services.ScheduleUpdate),
		ScheduleCreateTimeSlot: adminTimeSlotHandler.NewCreate(services.ScheduleCreateTimeSlot),
		ScheduleUpdateTimeSlot: adminTimeSlotHandler.NewUpdate(services.ScheduleUpdateTimeSlot),
		ScheduleDeleteTimeSlot: adminTimeSlotHandler.NewDelete(services.ScheduleDeleteTimeSlot),
		ScheduleGetAll:         adminScheduleHandler.NewGetAll(services.ScheduleGetAll),
		ScheduleGet:            adminScheduleHandler.NewGet(services.ScheduleGet),

		// Time slot template handlers
		TimeSlotTemplateGetAll: adminTimeSlotTemplateHandler.NewGetAll(services.TimeSlotTemplateGetAll),
		TimeSlotTemplateGet:    adminTimeSlotTemplateHandler.NewGet(services.TimeSlotTemplateGet),
		TimeSlotTemplateCreate: adminTimeSlotTemplateHandler.NewCreate(services.TimeSlotTemplateCreate),
		TimeSlotTemplateUpdate: adminTimeSlotTemplateHandler.NewUpdate(services.TimeSlotTemplateUpdate),
		TimeSlotTemplateDelete: adminTimeSlotTemplateHandler.NewDelete(services.TimeSlotTemplateDelete),

		// Time slot template item handlers
		TimeSlotTemplateItemCreate: adminTimeSlotTemplateItemHandler.NewCreate(services.TimeSlotTemplateItemCreate),
		TimeSlotTemplateUpdateItem: adminTimeSlotTemplateItemHandler.NewUpdate(services.TimeSlotTemplateUpdateItem),
		TimeSlotTemplateDeleteItem: adminTimeSlotTemplateItemHandler.NewDelete(services.TimeSlotTemplateDeleteItem),

		// Coupon management handlers
		CouponCreate: adminCouponHandler.NewCreate(services.CouponCreate),
		CouponGetAll: adminCouponHandler.NewGetAll(services.CouponGetAll),
		CouponUpdate: adminCouponHandler.NewUpdate(services.CouponUpdate),

		// Customer coupon handlers
		CustomerCouponGetAll: adminCustomerCouponHandler.NewGetAll(services.CustomerCouponGetAll),
		CustomerCouponCreate: adminCustomerCouponHandler.NewCreate(services.CustomerCouponCreate),

		// Checkout handlers
		CheckoutCreate: adminCheckoutHandler.NewCreate(services.CheckoutCreate),

		// Report handlers
		ReportGetPerformanceMe:    adminReportHandler.NewGetPerformanceMe(services.ReportGetPerformanceMe),
		ReportGetStorePerformance: adminReportHandler.NewGetStorePerformance(services.ReportGetStorePerformance),
	}
}
