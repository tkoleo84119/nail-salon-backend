package app

import (
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	"github.com/tkoleo84119/nail-salon-backend/internal/infra/db"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"

	// Admin handlers
	adminAuthHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/auth"
	adminScheduleHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/schedule"
	adminServiceHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/service"
	adminStaffHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/staff"
	adminStoreHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/store"
	adminStylistHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/stylist"
	adminTimeSlotTemplateHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/time-slot-template"

	// Admin services
	adminAuthService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/auth"
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
	StaffGetList           adminStaffService.GetStaffListServiceInterface
	StaffStoreAccess       *adminStaffService.CreateStoreAccessService
	StaffDeleteStoreAccess *adminStaffService.DeleteStoreAccessBulkService

	// Store management services
	StoreCreate *adminStoreService.CreateStoreService
	StoreUpdate *adminStoreService.UpdateStoreService

	// Service management services
	ServiceCreate *adminServiceService.CreateServiceService
	ServiceUpdate *adminServiceService.UpdateServiceService

	// Stylist management services
	StylistCreate *adminStylistService.CreateMyStylistService
	StylistUpdate *adminStylistService.UpdateMyStylistService

	// Schedule management services
	ScheduleCreateBulk     *adminScheduleService.CreateSchedulesBulkService
	ScheduleDeleteBulk     *adminScheduleService.DeleteSchedulesBulkService
	ScheduleCreateTimeSlot *adminScheduleService.CreateTimeSlotService
	ScheduleUpdateTimeSlot *adminScheduleService.UpdateTimeSlotService
	ScheduleDeleteTimeSlot *adminScheduleService.DeleteTimeSlotService

	// Time slot template services
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
	StaffGetList           *adminStaffHandler.GetStaffListHandler
	StaffStoreAccess       *adminStaffHandler.CreateStoreAccessHandler
	StaffDeleteStoreAccess *adminStaffHandler.DeleteStoreAccessBulkHandler

	// Store management handlers
	StoreCreate *adminStoreHandler.CreateStoreHandler
	StoreUpdate *adminStoreHandler.UpdateStoreHandler

	// Service management handlers
	ServiceCreate *adminServiceHandler.CreateServiceHandler
	ServiceUpdate *adminServiceHandler.UpdateServiceHandler

	// Stylist management handlers
	StylistCreate *adminStylistHandler.CreateMyStylistHandler
	StylistUpdate *adminStylistHandler.UpdateMyStylistHandler

	// Schedule management handlers
	ScheduleCreateBulk     *adminScheduleHandler.CreateSchedulesBulkHandler
	ScheduleDeleteBulk     *adminScheduleHandler.DeleteSchedulesBulkHandler
	ScheduleCreateTimeSlot *adminScheduleHandler.CreateTimeSlotHandler
	ScheduleUpdateTimeSlot *adminScheduleHandler.UpdateTimeSlotHandler
	ScheduleDeleteTimeSlot *adminScheduleHandler.DeleteTimeSlotHandler

	// Time slot template handlers
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
		StaffGetList:           adminStaffService.NewGetStaffListService(repositories.StaffUser),
		StaffStoreAccess:       adminStaffService.NewCreateStoreAccessService(queries),
		StaffDeleteStoreAccess: adminStaffService.NewDeleteStoreAccessBulkService(queries),

		// Store management services
		StoreCreate: adminStoreService.NewCreateStoreService(queries, database.PgxPool),
		StoreUpdate: adminStoreService.NewUpdateStoreService(queries, repositories.Store),

		// Service management services
		ServiceCreate: adminServiceService.NewCreateServiceService(queries),
		ServiceUpdate: adminServiceService.NewUpdateServiceService(queries, repositories.Service),

		// Stylist management services
		StylistCreate: adminStylistService.NewCreateMyStylistService(queries),
		StylistUpdate: adminStylistService.NewUpdateMyStylistService(queries, repositories.Stylist),

		// Schedule management services
		ScheduleCreateBulk:     adminScheduleService.NewCreateSchedulesBulkService(queries, database.PgxPool),
		ScheduleDeleteBulk:     adminScheduleService.NewDeleteSchedulesBulkService(queries, database.PgxPool),
		ScheduleCreateTimeSlot: adminScheduleService.NewCreateTimeSlotService(queries),
		ScheduleUpdateTimeSlot: adminScheduleService.NewUpdateTimeSlotService(queries, repositories.TimeSlot),
		ScheduleDeleteTimeSlot: adminScheduleService.NewDeleteTimeSlotService(queries),

		// Time slot template services
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
		StaffGetList:           adminStaffHandler.NewGetStaffListHandler(services.StaffGetList),
		StaffStoreAccess:       adminStaffHandler.NewCreateStoreAccessHandler(services.StaffStoreAccess),
		StaffDeleteStoreAccess: adminStaffHandler.NewDeleteStoreAccessBulkHandler(services.StaffDeleteStoreAccess),

		// Store management handlers
		StoreCreate: adminStoreHandler.NewCreateStoreHandler(services.StoreCreate),
		StoreUpdate: adminStoreHandler.NewUpdateStoreHandler(services.StoreUpdate),

		// Service management handlers
		ServiceCreate: adminServiceHandler.NewCreateServiceHandler(services.ServiceCreate),
		ServiceUpdate: adminServiceHandler.NewUpdateServiceHandler(services.ServiceUpdate),

		// Stylist management handlers
		StylistCreate: adminStylistHandler.NewCreateMyStylistHandler(services.StylistCreate),
		StylistUpdate: adminStylistHandler.NewUpdateMyStylistHandler(services.StylistUpdate),

		// Schedule management handlers
		ScheduleCreateBulk:     adminScheduleHandler.NewCreateSchedulesBulkHandler(services.ScheduleCreateBulk),
		ScheduleDeleteBulk:     adminScheduleHandler.NewDeleteSchedulesBulkHandler(services.ScheduleDeleteBulk),
		ScheduleCreateTimeSlot: adminScheduleHandler.NewCreateTimeSlotHandler(services.ScheduleCreateTimeSlot),
		ScheduleUpdateTimeSlot: adminScheduleHandler.NewUpdateTimeSlotHandler(services.ScheduleUpdateTimeSlot),
		ScheduleDeleteTimeSlot: adminScheduleHandler.NewDeleteTimeSlotHandler(services.ScheduleDeleteTimeSlot),

		// Time slot template handlers
		TimeSlotTemplateCreate:     adminTimeSlotTemplateHandler.NewCreateTimeSlotTemplateHandler(services.TimeSlotTemplateCreate),
		TimeSlotTemplateUpdate:     adminTimeSlotTemplateHandler.NewUpdateTimeSlotTemplateHandler(services.TimeSlotTemplateUpdate),
		TimeSlotTemplateDelete:     adminTimeSlotTemplateHandler.NewDeleteTimeSlotTemplateHandler(services.TimeSlotTemplateDelete),
		TimeSlotTemplateCreateItem: adminTimeSlotTemplateHandler.NewCreateTimeSlotTemplateItemHandler(services.TimeSlotTemplateCreateItem),
		TimeSlotTemplateUpdateItem: adminTimeSlotTemplateHandler.NewUpdateTimeSlotTemplateItemHandler(services.TimeSlotTemplateUpdateItem),
		TimeSlotTemplateDeleteItem: adminTimeSlotTemplateHandler.NewDeleteTimeSlotTemplateItemHandler(services.TimeSlotTemplateDeleteItem),
	}
}
