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

	// Admin services
	adminAuthService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/auth"
	adminBookingService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/booking"
	adminScheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/schedule"
	adminServiceService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/service"
	adminStaffService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/staff"
	adminStoreService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/store"
	adminStylistService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/stylist"
	adminTimeSlotTemplateService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/time-slot-template"
)

// AdminServices contains all admin-facing services
type AdminServices struct {
	// Authentication services
	AuthStaffLogin        *adminAuthService.StaffLoginService
	AuthStaffRefreshToken adminAuthService.StaffRefreshTokenServiceInterface

	// Staff management services
	StaffCreate            *adminStaffService.CreateStaffService
	StaffUpdate            *adminStaffService.UpdateStaffService
	StaffUpdateMe          *adminStaffService.UpdateMyStaffService
	StaffGet               adminStaffService.GetStaffServiceInterface
	StaffGetMe             adminStaffService.GetMyStaffServiceInterface
	StaffGetList           adminStaffService.GetStaffListServiceInterface
	StaffGetStoreAccess    adminStaffService.GetStaffStoreAccessServiceInterface
	StaffStoreAccess       *adminStaffService.CreateStoreAccessService
	StaffDeleteStoreAccess *adminStaffService.DeleteStoreAccessBulkService

	// Store management services
	StoreGetList adminStoreService.GetStoreListServiceInterface
	StoreGet     adminStoreService.GetStoreServiceInterface
	StoreCreate  adminStoreService.CreateStoreServiceInterface
	StoreUpdate  adminStoreService.UpdateStoreServiceInterface

	// Service management services
	ServiceGetList adminServiceService.GetServiceListServiceInterface
	ServiceGet     adminServiceService.GetServiceServiceInterface
	ServiceCreate  *adminServiceService.CreateServiceService
	ServiceUpdate  *adminServiceService.UpdateServiceService

	// Stylist management services
	StylistCreate  *adminStylistService.CreateMyStylistService
	StylistUpdate  *adminStylistService.UpdateMyStylistService
	StylistGetList adminStylistService.GetStylistListServiceInterface

	// Booking management services
	BookingCreate        adminBookingService.CreateBookingServiceInterface
	BookingGetList       adminBookingService.GetBookingListServiceInterface
	BookingUpdateByStaff adminBookingService.UpdateBookingByStaffServiceInterface

	// Schedule management services
	ScheduleCreateBulk     adminScheduleService.CreateSchedulesBulkServiceInterface
	ScheduleDeleteBulk     adminScheduleService.DeleteSchedulesBulkServiceInterface
	ScheduleCreateTimeSlot adminScheduleService.CreateTimeSlotServiceInterface
	ScheduleUpdateTimeSlot adminScheduleService.UpdateTimeSlotServiceInterface
	ScheduleDeleteTimeSlot adminScheduleService.DeleteTimeSlotServiceInterface
	ScheduleGetList        adminScheduleService.GetScheduleListServiceInterface
	ScheduleGet            adminScheduleService.GetScheduleServiceInterface

	// Time slot template services
	TimeSlotTemplateGetList    adminTimeSlotTemplateService.GetTimeSlotTemplateListServiceInterface
	TimeSlotTemplateGet        adminTimeSlotTemplateService.GetTimeSlotTemplateServiceInterface
	TimeSlotTemplateCreate     *adminTimeSlotTemplateService.CreateTimeSlotTemplateService
	TimeSlotTemplateUpdate     *adminTimeSlotTemplateService.UpdateTimeSlotTemplateService
	TimeSlotTemplateDelete     *adminTimeSlotTemplateService.DeleteTimeSlotTemplateService
	TimeSlotTemplateCreateItem *adminTimeSlotTemplateService.CreateTimeSlotTemplateItemService
	TimeSlotTemplateUpdateItem *adminTimeSlotTemplateService.UpdateTimeSlotTemplateItemService
	TimeSlotTemplateDeleteItem *adminTimeSlotTemplateService.DeleteTimeSlotTemplateItemService
}

// AdminHandlers contains all admin-facing handlers
type AdminHandlers struct {
	// Authentication handlers
	AuthStaffLogin        *adminAuthHandler.StaffLoginHandler
	AuthStaffRefreshToken *adminAuthHandler.StaffRefreshTokenHandler

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
	StylistCreate  *adminStylistHandler.CreateMyStylistHandler
	StylistUpdate  *adminStylistHandler.UpdateMyStylistHandler
	StylistGetList *adminStylistHandler.GetStylistListHandler

	// Booking management handlers
	BookingCreate        *adminBookingHandler.CreateBookingHandler
	BookingGetList       *adminBookingHandler.GetBookingListHandler
	BookingUpdateByStaff *adminBookingHandler.UpdateBookingByStaffHandler

	// Schedule management handlers
	ScheduleCreateBulk     *adminScheduleHandler.CreateSchedulesBulkHandler
	ScheduleDeleteBulk     *adminScheduleHandler.DeleteSchedulesBulkHandler
	ScheduleCreateTimeSlot *adminScheduleHandler.CreateTimeSlotHandler
	ScheduleUpdateTimeSlot *adminScheduleHandler.UpdateTimeSlotHandler
	ScheduleDeleteTimeSlot *adminScheduleHandler.DeleteTimeSlotHandler
	ScheduleGetList        *adminScheduleHandler.GetScheduleListHandler
	ScheduleGet            *adminScheduleHandler.GetScheduleHandler

	// Time slot template handlers
	TimeSlotTemplateGetList    *adminTimeSlotTemplateHandler.GetTimeSlotTemplateListHandler
	TimeSlotTemplateGet        *adminTimeSlotTemplateHandler.GetTimeSlotTemplateHandler
	TimeSlotTemplateCreate     *adminTimeSlotTemplateHandler.CreateTimeSlotTemplateHandler
	TimeSlotTemplateUpdate     *adminTimeSlotTemplateHandler.UpdateTimeSlotTemplateHandler
	TimeSlotTemplateDelete     *adminTimeSlotTemplateHandler.DeleteTimeSlotTemplateHandler
	TimeSlotTemplateCreateItem *adminTimeSlotTemplateHandler.CreateTimeSlotTemplateItemHandler
	TimeSlotTemplateUpdateItem *adminTimeSlotTemplateHandler.UpdateTimeSlotTemplateItemHandler
	TimeSlotTemplateDeleteItem *adminTimeSlotTemplateHandler.DeleteTimeSlotTemplateItemHandler
}

// NewAdminServices creates and initializes all admin services
func NewAdminServices(queries *dbgen.Queries, database *db.Database, repositories Repositories, cfg *config.Config) AdminServices {
	return AdminServices{
		// Authentication services
		AuthStaffLogin:        adminAuthService.NewStaffLoginService(queries, cfg.JWT),
		AuthStaffRefreshToken: adminAuthService.NewStaffRefreshTokenService(queries, cfg.JWT),

		// Staff management services
		StaffCreate:            adminStaffService.NewCreateStaffService(queries, database.PgxPool),
		StaffUpdate:            adminStaffService.NewUpdateStaffService(queries, database.Sqlx),
		StaffUpdateMe:          adminStaffService.NewUpdateMyStaffService(queries, repositories.StaffUser),
		StaffGet:               adminStaffService.NewGetStaffService(queries),
		StaffGetMe:             adminStaffService.NewGetMyStaffService(queries),
		StaffGetList:           adminStaffService.NewGetStaffListService(repositories.StaffUser),
		StaffGetStoreAccess:    adminStaffService.NewGetStaffStoreAccessService(queries),
		StaffStoreAccess:       adminStaffService.NewCreateStoreAccessService(queries),
		StaffDeleteStoreAccess: adminStaffService.NewDeleteStoreAccessBulkService(queries),

		// Store management services
		StoreGetList: adminStoreService.NewGetStoreListService(repositories.Store),
		StoreGet:     adminStoreService.NewGetStoreService(queries),
		StoreCreate:  adminStoreService.NewCreateStoreService(queries, database.PgxPool),
		StoreUpdate:  adminStoreService.NewUpdateStoreService(queries, repositories.Store),

		// Service management services
		ServiceGetList: adminServiceService.NewGetServiceListService(queries, repositories.Service),
		ServiceGet:     adminServiceService.NewGetServiceService(queries),
		ServiceCreate:  adminServiceService.NewCreateServiceService(queries),
		ServiceUpdate:  adminServiceService.NewUpdateServiceService(queries, repositories.Service),

		// Stylist management services
		StylistCreate:  adminStylistService.NewCreateMyStylistService(queries),
		StylistUpdate:  adminStylistService.NewUpdateMyStylistService(queries, repositories.Stylist),
		StylistGetList: adminStylistService.NewGetStylistListService(queries, repositories.Stylist),

		// Booking management services
		BookingCreate:        adminBookingService.NewCreateBookingService(queries, database.PgxPool),
		BookingGetList:       adminBookingService.NewGetBookingListService(queries, repositories.Booking),
		BookingUpdateByStaff: adminBookingService.NewUpdateBookingByStaffService(queries, database.PgxPool, repositories.Booking),

		// Schedule management services
		ScheduleCreateBulk:     adminScheduleService.NewCreateSchedulesBulkService(queries, database.PgxPool),
		ScheduleDeleteBulk:     adminScheduleService.NewDeleteSchedulesBulkService(queries, database.PgxPool),
		ScheduleCreateTimeSlot: adminScheduleService.NewCreateTimeSlotService(queries),
		ScheduleUpdateTimeSlot: adminScheduleService.NewUpdateTimeSlotService(queries, repositories.TimeSlot),
		ScheduleDeleteTimeSlot: adminScheduleService.NewDeleteTimeSlotService(queries),
		ScheduleGetList:        adminScheduleService.NewGetScheduleListService(queries, repositories.Schedule),
		ScheduleGet:            adminScheduleService.NewGetScheduleService(queries),

		// Time slot template services
		TimeSlotTemplateGetList:    adminTimeSlotTemplateService.NewGetTimeSlotTemplateListService(queries, repositories.TimeSlotTemplate),
		TimeSlotTemplateGet:        adminTimeSlotTemplateService.NewGetTimeSlotTemplateService(queries),
		TimeSlotTemplateCreate:     adminTimeSlotTemplateService.NewCreateTimeSlotTemplateService(queries, database.PgxPool),
		TimeSlotTemplateUpdate:     adminTimeSlotTemplateService.NewUpdateTimeSlotTemplateService(queries, repositories.TimeSlotTemplate),
		TimeSlotTemplateDelete:     adminTimeSlotTemplateService.NewDeleteTimeSlotTemplateService(queries),
		TimeSlotTemplateCreateItem: adminTimeSlotTemplateService.NewCreateTimeSlotTemplateItemService(queries),
		TimeSlotTemplateUpdateItem: adminTimeSlotTemplateService.NewUpdateTimeSlotTemplateItemService(queries),
		TimeSlotTemplateDeleteItem: adminTimeSlotTemplateService.NewDeleteTimeSlotTemplateItemService(queries),
	}
}

// NewAdminHandlers creates and initializes all admin handlers
func NewAdminHandlers(services AdminServices) AdminHandlers {
	return AdminHandlers{
		// Authentication handlers
		AuthStaffLogin:        adminAuthHandler.NewStaffLoginHandler(services.AuthStaffLogin),
		AuthStaffRefreshToken: adminAuthHandler.NewStaffRefreshTokenHandler(services.AuthStaffRefreshToken),

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
		StylistCreate:  adminStylistHandler.NewCreateMyStylistHandler(services.StylistCreate),
		StylistUpdate:  adminStylistHandler.NewUpdateMyStylistHandler(services.StylistUpdate),
		StylistGetList: adminStylistHandler.NewGetStylistListHandler(services.StylistGetList),

		// Booking management handlers
		BookingCreate:        adminBookingHandler.NewCreateBookingHandler(services.BookingCreate),
		BookingGetList:       adminBookingHandler.NewGetBookingListHandler(services.BookingGetList),
		BookingUpdateByStaff: adminBookingHandler.NewUpdateBookingByStaffHandler(services.BookingUpdateByStaff),

		// Schedule management handlers
		ScheduleCreateBulk:     adminScheduleHandler.NewCreateSchedulesBulkHandler(services.ScheduleCreateBulk),
		ScheduleDeleteBulk:     adminScheduleHandler.NewDeleteSchedulesBulkHandler(services.ScheduleDeleteBulk),
		ScheduleCreateTimeSlot: adminScheduleHandler.NewCreateTimeSlotHandler(services.ScheduleCreateTimeSlot),
		ScheduleUpdateTimeSlot: adminScheduleHandler.NewUpdateTimeSlotHandler(services.ScheduleUpdateTimeSlot),
		ScheduleDeleteTimeSlot: adminScheduleHandler.NewDeleteTimeSlotHandler(services.ScheduleDeleteTimeSlot),
		ScheduleGetList:        adminScheduleHandler.NewGetScheduleListHandler(services.ScheduleGetList),
		ScheduleGet:            adminScheduleHandler.NewGetScheduleHandler(services.ScheduleGet),

		// Time slot template handlers
		TimeSlotTemplateGetList:    adminTimeSlotTemplateHandler.NewGetTimeSlotTemplateListHandler(services.TimeSlotTemplateGetList),
		TimeSlotTemplateGet:        adminTimeSlotTemplateHandler.NewGetTimeSlotTemplateHandler(services.TimeSlotTemplateGet),
		TimeSlotTemplateCreate:     adminTimeSlotTemplateHandler.NewCreateTimeSlotTemplateHandler(services.TimeSlotTemplateCreate),
		TimeSlotTemplateUpdate:     adminTimeSlotTemplateHandler.NewUpdateTimeSlotTemplateHandler(services.TimeSlotTemplateUpdate),
		TimeSlotTemplateDelete:     adminTimeSlotTemplateHandler.NewDeleteTimeSlotTemplateHandler(services.TimeSlotTemplateDelete),
		TimeSlotTemplateCreateItem: adminTimeSlotTemplateHandler.NewCreateTimeSlotTemplateItemHandler(services.TimeSlotTemplateCreateItem),
		TimeSlotTemplateUpdateItem: adminTimeSlotTemplateHandler.NewUpdateTimeSlotTemplateItemHandler(services.TimeSlotTemplateUpdateItem),
		TimeSlotTemplateDeleteItem: adminTimeSlotTemplateHandler.NewDeleteTimeSlotTemplateItemHandler(services.TimeSlotTemplateDeleteItem),
	}
}
