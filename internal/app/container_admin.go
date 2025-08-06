package app

import (
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	"github.com/tkoleo84119/nail-salon-backend/internal/infra/db"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"

	// Admin handlers
	adminAuthHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/auth"
	adminBookingHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/booking"
	adminScheduleHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/schedule"
	adminServiceHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/service"
	adminStaffHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/staff"
	adminStoreHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/store"
	adminStylistHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/stylist"
	adminTimeSlotTemplateHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/time-slot-template"
	adminTimeSlotHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/time_slot"
	adminTimeSlotTemplateItemHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/time_slot_template_item"

	// Admin services
	adminAuthService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/auth"
	adminBookingService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/booking"
	adminScheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/schedule"
	adminServiceService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/service"
	adminStaffService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/staff"
	adminStoreService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/store"
	adminStylistService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/stylist"
	adminTimeSlotTemplateService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/time-slot-template"
	adminTimeSlotService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/time_slot"
	adminTimeSlotTemplateItemService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/time_slot_template_item"
)

// AdminServices contains all admin-facing services
type AdminServices struct {
	// Authentication services
	AuthStaffLogin        adminAuthService.StaffLoginServiceInterface
	AuthStaffRefreshToken adminAuthService.StaffRefreshTokenServiceInterface
	AuthStaffLogout       adminAuthService.StaffLogoutServiceInterface

	// Staff management services
	StaffCreate            adminStaffService.CreateStaffServiceInterface
	StaffUpdate            adminStaffService.UpdateStaffServiceInterface
	StaffUpdateMe          adminStaffService.UpdateMyStaffServiceInterface
	StaffGet               adminStaffService.GetStaffServiceInterface
	StaffGetMe             adminStaffService.GetMyStaffServiceInterface
	StaffGetList           adminStaffService.GetStaffListServiceInterface
	StaffGetStoreAccess    adminStaffService.GetStaffStoreAccessServiceInterface
	StaffStoreAccess       adminStaffService.CreateStoreAccessServiceInterface
	StaffDeleteStoreAccess adminStaffService.DeleteStoreAccessBulkServiceInterface

	// Store management services
	StoreGetList adminStoreService.GetStoreListServiceInterface
	StoreGet     adminStoreService.GetStoreServiceInterface
	StoreCreate  adminStoreService.CreateStoreServiceInterface
	StoreUpdate  adminStoreService.UpdateStoreServiceInterface

	// Service management services
	ServiceGetList adminServiceService.GetServiceListServiceInterface
	ServiceGet     adminServiceService.GetServiceServiceInterface
	ServiceCreate  adminServiceService.CreateServiceInterface
	ServiceUpdate  adminServiceService.UpdateServiceInterface

	// Stylist management services
	StylistUpdate  adminStylistService.UpdateMyStylistServiceInterface
	StylistGetList adminStylistService.GetStylistListServiceInterface

	// Booking management services
	BookingCreate        adminBookingService.CreateBookingServiceInterface
	BookingGetList       adminBookingService.GetBookingListServiceInterface
	BookingUpdateByStaff adminBookingService.UpdateBookingByStaffServiceInterface
	BookingCancel        adminBookingService.CancelBookingServiceInterface

	// Schedule management services
	ScheduleCreateBulk     adminScheduleService.CreateBulkInterface
	ScheduleDeleteBulk     adminScheduleService.DeleteBulkInterface
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
}

// AdminHandlers contains all admin-facing handlers
type AdminHandlers struct {
	// Authentication handlers
	AuthStaffLogin        *adminAuthHandler.StaffLoginHandler
	AuthStaffRefreshToken *adminAuthHandler.StaffRefreshTokenHandler
	AuthStaffLogout       *adminAuthHandler.StaffLogoutHandler

	// Staff management handlers
	StaffCreate            *adminStaffHandler.CreateStaffHandler
	StaffUpdate            *adminStaffHandler.UpdateStaffHandler
	StaffUpdateMe          *adminStaffHandler.UpdateMyStaffHandler
	StaffGet               *adminStaffHandler.GetStaffHandler
	StaffGetMe             *adminStaffHandler.GetMyStaffHandler
	StaffGetList           *adminStaffHandler.GetStaffListHandler
	StaffGetStoreAccess    *adminStaffHandler.GetStaffStoreAccessHandler
	StaffStoreAccess       *adminStaffHandler.CreateStoreAccessHandler
	StaffDeleteStoreAccess *adminStaffHandler.DeleteStoreAccessBulkHandler

	// Store management handlers
	StoreGetList *adminStoreHandler.GetStoreListHandler
	StoreGet     *adminStoreHandler.GetStoreHandler
	StoreCreate  *adminStoreHandler.CreateStoreHandler
	StoreUpdate  *adminStoreHandler.UpdateStoreHandler

	// Service management handlers
	ServiceGetList *adminServiceHandler.GetServiceListHandler
	ServiceGet     *adminServiceHandler.GetServiceHandler
	ServiceCreate  *adminServiceHandler.CreateServiceHandler
	ServiceUpdate  *adminServiceHandler.UpdateServiceHandler

	// Stylist management handlers
	StylistUpdate  *adminStylistHandler.UpdateMyStylistHandler
	StylistGetList *adminStylistHandler.GetStylistListHandler

	// Booking management handlers
	BookingCreate        *adminBookingHandler.CreateBookingHandler
	BookingGetList       *adminBookingHandler.GetBookingListHandler
	BookingUpdateByStaff *adminBookingHandler.UpdateBookingByStaffHandler
	BookingCancel        *adminBookingHandler.CancelBookingHandler

	// Schedule management handlers
	ScheduleCreateBulk     *adminScheduleHandler.CreateBulk
	ScheduleDeleteBulk     *adminScheduleHandler.DeleteBulk
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
}

// NewAdminServices creates and initializes all admin services
func NewAdminServices(queries *dbgen.Queries, database *db.Database, repositories Repositories, cfg *config.Config) AdminServices {
	return AdminServices{
		// Authentication services
		AuthStaffLogin:        adminAuthService.NewStaffLoginService(repositories.SQLX, cfg.JWT),
		AuthStaffRefreshToken: adminAuthService.NewStaffRefreshTokenService(repositories.SQLX, cfg.JWT),
		AuthStaffLogout:       adminAuthService.NewStaffLogoutService(repositories.SQLX),

		// Staff management services
		StaffCreate:            adminStaffService.NewCreateStaffService(database.Sqlx, repositories.SQLX),
		StaffUpdate:            adminStaffService.NewUpdateStaffService(database.Sqlx),
		StaffUpdateMe:          adminStaffService.NewUpdateMyStaffService(repositories.SQLX),
		StaffGet:               adminStaffService.NewGetStaffService(repositories.SQLX),
		StaffGetMe:             adminStaffService.NewGetMyStaffService(repositories.SQLX),
		StaffGetList:           adminStaffService.NewGetStaffListService(repositories.SQLX),
		StaffGetStoreAccess:    adminStaffService.NewGetStaffStoreAccessService(repositories.SQLX),
		StaffStoreAccess:       adminStaffService.NewCreateStoreAccessService(repositories.SQLX),
		StaffDeleteStoreAccess: adminStaffService.NewDeleteStoreAccessBulkService(repositories.SQLX),

		// Store management services
		StoreGetList: adminStoreService.NewGetStoreListService(repositories.SQLX),
		StoreGet:     adminStoreService.NewGetStoreService(repositories.SQLX),
		StoreCreate:  adminStoreService.NewCreateStoreService(database.Sqlx, repositories.SQLX),
		StoreUpdate:  adminStoreService.NewUpdateStoreService(repositories.SQLX),

		// Service management services
		ServiceGetList: adminServiceService.NewGetServiceListService(repositories.SQLX),
		ServiceGet:     adminServiceService.NewGetServiceService(repositories.SQLX),
		ServiceCreate:  adminServiceService.NewCreateServiceService(repositories.SQLX),
		ServiceUpdate:  adminServiceService.NewUpdateServiceService(repositories.SQLX),

		// Stylist management services
		StylistUpdate:  adminStylistService.NewUpdateMyStylistService(repositories.SQLX),
		StylistGetList: adminStylistService.NewGetStylistListService(repositories.SQLX),

		// Booking management services
		BookingCreate:        adminBookingService.NewCreateBookingService(queries, database.PgxPool),
		BookingGetList:       adminBookingService.NewGetBookingListService(queries, repositories.SQLX),
		BookingUpdateByStaff: adminBookingService.NewUpdateBookingByStaffService(queries, database.PgxPool, repositories.SQLX),
		BookingCancel:        adminBookingService.NewCancelBookingService(database.Sqlx, repositories.SQLX),

		// Schedule management services
		ScheduleCreateBulk:     adminScheduleService.NewCreateBulk(queries, database.PgxPool),
		ScheduleDeleteBulk:     adminScheduleService.NewDeleteBulk(queries),
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
	}
}

// NewAdminHandlers creates and initializes all admin handlers
func NewAdminHandlers(services AdminServices) AdminHandlers {
	return AdminHandlers{
		// Authentication handlers
		AuthStaffLogin:        adminAuthHandler.NewStaffLoginHandler(services.AuthStaffLogin),
		AuthStaffRefreshToken: adminAuthHandler.NewStaffRefreshTokenHandler(services.AuthStaffRefreshToken),
		AuthStaffLogout:       adminAuthHandler.NewStaffLogoutHandler(services.AuthStaffLogout),

		// Staff management handlers
		StaffCreate:            adminStaffHandler.NewCreateStaffHandler(services.StaffCreate),
		StaffUpdate:            adminStaffHandler.NewUpdateStaffHandler(services.StaffUpdate),
		StaffUpdateMe:          adminStaffHandler.NewUpdateMyStaffHandler(services.StaffUpdateMe),
		StaffGet:               adminStaffHandler.NewGetStaffHandler(services.StaffGet),
		StaffGetMe:             adminStaffHandler.NewGetMyStaffHandler(services.StaffGetMe),
		StaffGetList:           adminStaffHandler.NewGetStaffListHandler(services.StaffGetList),
		StaffGetStoreAccess:    adminStaffHandler.NewGetStaffStoreAccessHandler(services.StaffGetStoreAccess),
		StaffStoreAccess:       adminStaffHandler.NewCreateStoreAccessHandler(services.StaffStoreAccess),
		StaffDeleteStoreAccess: adminStaffHandler.NewDeleteStoreAccessBulkHandler(services.StaffDeleteStoreAccess),

		// Store management handlers
		StoreGetList: adminStoreHandler.NewGetStoreListHandler(services.StoreGetList),
		StoreGet:     adminStoreHandler.NewGetStoreHandler(services.StoreGet),
		StoreCreate:  adminStoreHandler.NewCreateStoreHandler(services.StoreCreate),
		StoreUpdate:  adminStoreHandler.NewUpdateStoreHandler(services.StoreUpdate),

		// Service management handlers
		ServiceGetList: adminServiceHandler.NewGetServiceListHandler(services.ServiceGetList),
		ServiceGet:     adminServiceHandler.NewGetServiceHandler(services.ServiceGet),
		ServiceCreate:  adminServiceHandler.NewCreateServiceHandler(services.ServiceCreate),
		ServiceUpdate:  adminServiceHandler.NewUpdateServiceHandler(services.ServiceUpdate),

		// Stylist management handlers
		StylistUpdate:  adminStylistHandler.NewUpdateMyStylistHandler(services.StylistUpdate),
		StylistGetList: adminStylistHandler.NewGetStylistListHandler(services.StylistGetList),

		// Booking management handlers
		BookingCreate:        adminBookingHandler.NewCreateBookingHandler(services.BookingCreate),
		BookingGetList:       adminBookingHandler.NewGetBookingListHandler(services.BookingGetList),
		BookingUpdateByStaff: adminBookingHandler.NewUpdateBookingByStaffHandler(services.BookingUpdateByStaff),
		BookingCancel:        adminBookingHandler.NewCancelBookingHandler(services.BookingCancel),

		// Schedule management handlers
		ScheduleCreateBulk:     adminScheduleHandler.NewCreateBulk(services.ScheduleCreateBulk),
		ScheduleDeleteBulk:     adminScheduleHandler.NewDeleteBulk(services.ScheduleDeleteBulk),
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
	}
}
