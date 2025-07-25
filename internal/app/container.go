package app

import (
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	adminAuthHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/auth"
	adminScheduleHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/schedule"
	adminServiceHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/service"
	adminStaffHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/staff"
	adminStoreHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/store"
	adminStylistHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/stylist"
	adminTimeSlotTemplateHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/admin/time-slot-template"
	authHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/auth"
	bookingHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/booking"
	customerHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/infra/db"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	adminAuthService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/auth"
	adminScheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/schedule"
	adminServiceService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/service"
	adminStaffService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/staff"
	adminStoreService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/store"
	adminStylistService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/stylist"
	adminTimeSlotTemplateService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/time-slot-template"
	authService "github.com/tkoleo84119/nail-salon-backend/internal/service/auth"
	bookingService "github.com/tkoleo84119/nail-salon-backend/internal/service/booking"
	customerService "github.com/tkoleo84119/nail-salon-backend/internal/service/customer"
)

type Container struct {
	cfg      *config.Config
	database *db.Database

	repositories Repositories
	services     Services
	handlers     Handlers
}

type Repositories struct {
	Stylist          *sqlx.StylistRepository
	StaffUser        *sqlx.StaffUserRepository
	Store            *sqlx.StoreRepository
	Service          *sqlx.ServiceRepository
	Customer         *sqlx.CustomerRepository
	TimeSlot         *sqlx.TimeSlotRepository
	TimeSlotTemplate *sqlx.TimeSlotTemplateRepository
	Booking          *sqlx.BookingRepository
}

type Services struct {
	AuthStaffLogin                  *adminAuthService.StaffLoginService
	AuthCustomerLineLogin           *authService.CustomerLineLoginService
	AuthCustomerLineRegister        *authService.CustomerLineRegisterService
	BookingCreateMy                 *bookingService.CreateMyBookingService
	BookingUpdateMy                 *bookingService.UpdateMyBookingService
	BookingCancelMy                 bookingService.CancelMyBookingServiceInterface
	CustomerGetMy                   *customerService.GetMyCustomerService
	CustomerUpdateMy                *customerService.UpdateMyCustomerService
	AdminStaffCreate                *adminStaffService.CreateStaffService
	AdminStaffUpdate                *adminStaffService.UpdateStaffService
	AdminStaffUpdateMe              *adminStaffService.UpdateMyStaffService
	AdminStaffStoreAccess           *adminStaffService.CreateStoreAccessService
	AdminStaffDeleteStoreAccess     *adminStaffService.DeleteStoreAccessBulkService
	AdminStoreCreate                *adminStoreService.CreateStoreService
	AdminStoreUpdate                *adminStoreService.UpdateStoreService
	AdminStylistCreate              *adminStylistService.CreateMyStylistService
	AdminStylistUpdate              *adminStylistService.UpdateMyStylistService
	AdminScheduleCreateBulk         *adminScheduleService.CreateSchedulesBulkService
	AdminScheduleDeleteBulk         *adminScheduleService.DeleteSchedulesBulkService
	AdminScheduleCreateTimeSlot     *adminScheduleService.CreateTimeSlotService
	AdminScheduleUpdateTimeSlot     *adminScheduleService.UpdateTimeSlotService
	AdminScheduleDeleteTimeSlot     *adminScheduleService.DeleteTimeSlotService
	AdminTimeSlotTemplateCreate     *adminTimeSlotTemplateService.CreateTimeSlotTemplateService
	AdminTimeSlotTemplateUpdate     *adminTimeSlotTemplateService.UpdateTimeSlotTemplateService
	AdminTimeSlotTemplateDelete     *adminTimeSlotTemplateService.DeleteTimeSlotTemplateService
	AdminTimeSlotTemplateCreateItem *adminTimeSlotTemplateService.CreateTimeSlotTemplateItemService
	AdminTimeSlotTemplateUpdateItem *adminTimeSlotTemplateService.UpdateTimeSlotTemplateItemService
	AdminTimeSlotTemplateDeleteItem *adminTimeSlotTemplateService.DeleteTimeSlotTemplateItemService
	AdminServiceCreate              *adminServiceService.CreateServiceService
	AdminServiceUpdate              *adminServiceService.UpdateServiceService
}

type Handlers struct {
	AuthStaffLogin                  *adminAuthHandler.StaffLoginHandler
	AuthCustomerLineLogin           *authHandler.CustomerLineLoginHandler
	AuthCustomerLineRegister        *authHandler.CustomerLineRegisterHandler
	BookingCreateMy                 *bookingHandler.CreateMyBookingHandler
	BookingUpdateMy                 *bookingHandler.UpdateMyBookingHandler
	BookingCancelMy                 *bookingHandler.CancelMyBookingHandler
	CustomerGetMy                   *customerHandler.GetMyCustomerHandler
	CustomerUpdateMy                *customerHandler.UpdateMyCustomerHandler
	AdminStaffCreate                *adminStaffHandler.CreateStaffHandler
	AdminStaffUpdate                *adminStaffHandler.UpdateStaffHandler
	AdminStaffUpdateMe              *adminStaffHandler.UpdateMyStaffHandler
	AdminStaffStoreAccess           *adminStaffHandler.CreateStoreAccessHandler
	AdminStaffDeleteStoreAccess     *adminStaffHandler.DeleteStoreAccessBulkHandler
	AdminStoreCreate                *adminStoreHandler.CreateStoreHandler
	AdminStoreUpdate                *adminStoreHandler.UpdateStoreHandler
	AdminStylistCreate              *adminStylistHandler.CreateMyStylistHandler
	AdminStylistUpdate              *adminStylistHandler.UpdateMyStylistHandler
	AdminScheduleCreateBulk         *adminScheduleHandler.CreateSchedulesBulkHandler
	AdminScheduleDeleteBulk         *adminScheduleHandler.DeleteSchedulesBulkHandler
	AdminScheduleCreateTimeSlot     *adminScheduleHandler.CreateTimeSlotHandler
	AdminScheduleUpdateTimeSlot     *adminScheduleHandler.UpdateTimeSlotHandler
	AdminScheduleDeleteTimeSlot     *adminScheduleHandler.DeleteTimeSlotHandler
	AdminTimeSlotTemplateCreate     *adminTimeSlotTemplateHandler.CreateTimeSlotTemplateHandler
	AdminTimeSlotTemplateUpdate     *adminTimeSlotTemplateHandler.UpdateTimeSlotTemplateHandler
	AdminTimeSlotTemplateDelete     *adminTimeSlotTemplateHandler.DeleteTimeSlotTemplateHandler
	AdminTimeSlotTemplateCreateItem *adminTimeSlotTemplateHandler.CreateTimeSlotTemplateItemHandler
	AdminTimeSlotTemplateUpdateItem *adminTimeSlotTemplateHandler.UpdateTimeSlotTemplateItemHandler
	AdminTimeSlotTemplateDeleteItem *adminTimeSlotTemplateHandler.DeleteTimeSlotTemplateItemHandler
	AdminServiceCreate              *adminServiceHandler.CreateServiceHandler
	AdminServiceUpdate              *adminServiceHandler.UpdateServiceHandler
}

func NewContainer(cfg *config.Config, database *db.Database) *Container {
	queries := dbgen.New(database.PgxPool)

	repositories := Repositories{
		Stylist:          sqlx.NewStylistRepository(database.Sqlx),
		StaffUser:        sqlx.NewStaffUserRepository(database.Sqlx),
		Store:            sqlx.NewStoreRepository(database.Sqlx),
		Service:          sqlx.NewServiceRepository(database.Sqlx),
		Customer:         sqlx.NewCustomerRepository(database.Sqlx),
		TimeSlot:         sqlx.NewTimeSlotRepository(database.Sqlx),
		TimeSlotTemplate: sqlx.NewTimeSlotTemplateRepository(database.Sqlx),
		Booking:          sqlx.NewBookingRepository(database.Sqlx),
	}

	services := Services{
		AuthStaffLogin:                  adminAuthService.NewStaffLoginService(queries, cfg.JWT),
		AuthCustomerLineLogin:           authService.NewCustomerLineLoginService(queries, cfg.Line, cfg.JWT),
		AuthCustomerLineRegister:        authService.NewCustomerLineRegisterService(queries, database.PgxPool, cfg.Line, cfg.JWT),
		BookingCreateMy:                 bookingService.NewCreateMyBookingService(queries, database.PgxPool),
		BookingUpdateMy:                 bookingService.NewUpdateMyBookingService(queries, repositories.Booking, database.PgxPool),
		BookingCancelMy:                 bookingService.NewCancelMyBookingService(queries),
		CustomerGetMy:                   customerService.NewGetMyCustomerService(queries),
		CustomerUpdateMy:                customerService.NewUpdateMyCustomerService(queries, repositories.Customer),
		AdminStaffCreate:                adminStaffService.NewCreateStaffService(queries, database.PgxPool),
		AdminStaffUpdate:                adminStaffService.NewUpdateStaffService(queries, database.Sqlx),
		AdminStaffUpdateMe:              adminStaffService.NewUpdateMyStaffService(queries, repositories.StaffUser),
		AdminStaffStoreAccess:           adminStaffService.NewCreateStoreAccessService(queries),
		AdminStaffDeleteStoreAccess:     adminStaffService.NewDeleteStoreAccessBulkService(queries),
		AdminStoreCreate:                adminStoreService.NewCreateStoreService(queries, database.PgxPool),
		AdminStoreUpdate:                adminStoreService.NewUpdateStoreService(queries, repositories.Store),
		AdminStylistCreate:              adminStylistService.NewCreateMyStylistService(queries),
		AdminStylistUpdate:              adminStylistService.NewUpdateMyStylistService(queries, repositories.Stylist),
		AdminScheduleCreateBulk:         adminScheduleService.NewCreateSchedulesBulkService(queries, database.PgxPool),
		AdminScheduleDeleteBulk:         adminScheduleService.NewDeleteSchedulesBulkService(queries, database.PgxPool),
		AdminScheduleCreateTimeSlot:     adminScheduleService.NewCreateTimeSlotService(queries),
		AdminScheduleUpdateTimeSlot:     adminScheduleService.NewUpdateTimeSlotService(queries, repositories.TimeSlot),
		AdminScheduleDeleteTimeSlot:     adminScheduleService.NewDeleteTimeSlotService(queries),
		AdminTimeSlotTemplateCreate:     adminTimeSlotTemplateService.NewCreateTimeSlotTemplateService(queries, database.PgxPool),
		AdminTimeSlotTemplateUpdate:     adminTimeSlotTemplateService.NewUpdateTimeSlotTemplateService(queries, repositories.TimeSlotTemplate),
		AdminTimeSlotTemplateDelete:     adminTimeSlotTemplateService.NewDeleteTimeSlotTemplateService(queries),
		AdminTimeSlotTemplateCreateItem: adminTimeSlotTemplateService.NewCreateTimeSlotTemplateItemService(queries),
		AdminTimeSlotTemplateUpdateItem: adminTimeSlotTemplateService.NewUpdateTimeSlotTemplateItemService(queries),
		AdminTimeSlotTemplateDeleteItem: adminTimeSlotTemplateService.NewDeleteTimeSlotTemplateItemService(queries),
		AdminServiceCreate:              adminServiceService.NewCreateServiceService(queries),
		AdminServiceUpdate:              adminServiceService.NewUpdateServiceService(queries, repositories.Service),
	}

	handlers := Handlers{
		AuthStaffLogin:                  adminAuthHandler.NewStaffLoginHandler(services.AuthStaffLogin),
		AuthCustomerLineLogin:           authHandler.NewCustomerLineLoginHandler(services.AuthCustomerLineLogin),
		AuthCustomerLineRegister:        authHandler.NewCustomerLineRegisterHandler(services.AuthCustomerLineRegister),
		BookingCreateMy:                 bookingHandler.NewCreateMyBookingHandler(services.BookingCreateMy),
		BookingUpdateMy:                 bookingHandler.NewUpdateMyBookingHandler(services.BookingUpdateMy),
		BookingCancelMy:                 bookingHandler.NewCancelMyBookingHandler(services.BookingCancelMy),
		CustomerGetMy:                   customerHandler.NewGetMyCustomerHandler(services.CustomerGetMy),
		CustomerUpdateMy:                customerHandler.NewUpdateMyCustomerHandler(services.CustomerUpdateMy),
		AdminStaffCreate:                adminStaffHandler.NewCreateStaffHandler(services.AdminStaffCreate),
		AdminStaffUpdate:                adminStaffHandler.NewUpdateStaffHandler(services.AdminStaffUpdate),
		AdminStaffUpdateMe:              adminStaffHandler.NewUpdateMyStaffHandler(services.AdminStaffUpdateMe),
		AdminStaffStoreAccess:           adminStaffHandler.NewCreateStoreAccessHandler(services.AdminStaffStoreAccess),
		AdminStaffDeleteStoreAccess:     adminStaffHandler.NewDeleteStoreAccessBulkHandler(services.AdminStaffDeleteStoreAccess),
		AdminStoreCreate:                adminStoreHandler.NewCreateStoreHandler(services.AdminStoreCreate),
		AdminStoreUpdate:                adminStoreHandler.NewUpdateStoreHandler(services.AdminStoreUpdate),
		AdminStylistCreate:              adminStylistHandler.NewCreateMyStylistHandler(services.AdminStylistCreate),
		AdminStylistUpdate:              adminStylistHandler.NewUpdateMyStylistHandler(services.AdminStylistUpdate),
		AdminScheduleCreateBulk:         adminScheduleHandler.NewCreateSchedulesBulkHandler(services.AdminScheduleCreateBulk),
		AdminScheduleDeleteBulk:         adminScheduleHandler.NewDeleteSchedulesBulkHandler(services.AdminScheduleDeleteBulk),
		AdminScheduleCreateTimeSlot:     adminScheduleHandler.NewCreateTimeSlotHandler(services.AdminScheduleCreateTimeSlot),
		AdminScheduleUpdateTimeSlot:     adminScheduleHandler.NewUpdateTimeSlotHandler(services.AdminScheduleUpdateTimeSlot),
		AdminScheduleDeleteTimeSlot:     adminScheduleHandler.NewDeleteTimeSlotHandler(services.AdminScheduleDeleteTimeSlot),
		AdminTimeSlotTemplateCreate:     adminTimeSlotTemplateHandler.NewCreateTimeSlotTemplateHandler(services.AdminTimeSlotTemplateCreate),
		AdminTimeSlotTemplateUpdate:     adminTimeSlotTemplateHandler.NewUpdateTimeSlotTemplateHandler(services.AdminTimeSlotTemplateUpdate),
		AdminTimeSlotTemplateDelete:     adminTimeSlotTemplateHandler.NewDeleteTimeSlotTemplateHandler(services.AdminTimeSlotTemplateDelete),
		AdminTimeSlotTemplateCreateItem: adminTimeSlotTemplateHandler.NewCreateTimeSlotTemplateItemHandler(services.AdminTimeSlotTemplateCreateItem),
		AdminTimeSlotTemplateUpdateItem: adminTimeSlotTemplateHandler.NewUpdateTimeSlotTemplateItemHandler(services.AdminTimeSlotTemplateUpdateItem),
		AdminTimeSlotTemplateDeleteItem: adminTimeSlotTemplateHandler.NewDeleteTimeSlotTemplateItemHandler(services.AdminTimeSlotTemplateDeleteItem),
		AdminServiceCreate:              adminServiceHandler.NewCreateServiceHandler(services.AdminServiceCreate),
		AdminServiceUpdate:              adminServiceHandler.NewUpdateServiceHandler(services.AdminServiceUpdate),
	}

	return &Container{
		cfg:          cfg,
		database:     database,
		repositories: repositories,
		services:     services,
		handlers:     handlers,
	}
}

func (c *Container) GetConfig() *config.Config {
	return c.cfg
}

func (c *Container) GetDatabase() *db.Database {
	return c.database
}

func (c *Container) GetRepositories() Repositories {
	return c.repositories
}

func (c *Container) GetServices() Services {
	return c.services
}

func (c *Container) GetHandlers() Handlers {
	return c.handlers
}
