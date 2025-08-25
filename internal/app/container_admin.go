package app

import (
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	"github.com/tkoleo84119/nail-salon-backend/internal/infra/db"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"

	// Admin handlers
	adminAuthHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/auth"
	adminBookingHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/booking"
	adminCouponHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/coupon"
	adminCustomerHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/customer"
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
	adminCouponService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/coupon"
	adminCustomerService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/customer"
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
	BookingCreate adminBookingService.CreateInterface
	BookingGetAll adminBookingService.GetAllInterface
	BookingUpdate adminBookingService.UpdateInterface
	BookingCancel adminBookingService.CancelInterface

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
	BookingCreate *adminBookingHandler.Create
	BookingGetAll *adminBookingHandler.GetAll
	BookingUpdate *adminBookingHandler.Update
	BookingCancel *adminBookingHandler.Cancel

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
		BookingCreate: adminBookingService.NewCreate(queries, database.PgxPool),
		BookingGetAll: adminBookingService.NewGetAll(queries, repositories.SQLX),
		BookingUpdate: adminBookingService.NewUpdate(queries, repositories.SQLX, database.Sqlx),
		BookingCancel: adminBookingService.NewCancel(queries, database.Sqlx, repositories.SQLX),

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
		BookingCreate: adminBookingHandler.NewCreate(services.BookingCreate),
		BookingGetAll: adminBookingHandler.NewGetAll(services.BookingGetAll),
		BookingUpdate: adminBookingHandler.NewUpdate(services.BookingUpdate),
		BookingCancel: adminBookingHandler.NewCancel(services.BookingCancel),

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
	}
}
