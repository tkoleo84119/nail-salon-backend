package app

import (
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	authHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/auth"
	bookingHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/booking"
	customerHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/customer"
	scheduleHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/schedule"
	serviceHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/service"
	staffHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/staff"
	storeHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/store"
	stylistHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/stylist"
	timeSlotTemplateHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/infra/db"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	authService "github.com/tkoleo84119/nail-salon-backend/internal/service/auth"
	bookingService "github.com/tkoleo84119/nail-salon-backend/internal/service/booking"
	customerService "github.com/tkoleo84119/nail-salon-backend/internal/service/customer"
	scheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/schedule"
	serviceService "github.com/tkoleo84119/nail-salon-backend/internal/service/service"
	staffService "github.com/tkoleo84119/nail-salon-backend/internal/service/staff"
	storeService "github.com/tkoleo84119/nail-salon-backend/internal/service/store"
	stylistService "github.com/tkoleo84119/nail-salon-backend/internal/service/stylist"
	timeSlotTemplateService "github.com/tkoleo84119/nail-salon-backend/internal/service/time-slot-template"
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
	AuthStaffLogin             *authService.StaffLoginService
	AuthCustomerLineLogin      *authService.CustomerLineLoginService
	AuthCustomerLineRegister   *authService.CustomerLineRegisterService
	BookingCreateMy            *bookingService.CreateMyBookingService
	BookingUpdateMy            *bookingService.UpdateMyBookingService
	BookingCancelMy            bookingService.CancelMyBookingServiceInterface
	CustomerUpdateMy           *customerService.UpdateMyCustomerService
	StaffCreate                *staffService.CreateStaffService
	StaffUpdate                *staffService.UpdateStaffService
	StaffUpdateMe              *staffService.UpdateMyStaffService
	StaffStoreAccess           *staffService.CreateStoreAccessService
	StaffDeleteStoreAccess     *staffService.DeleteStoreAccessBulkService
	StoreCreate                *storeService.CreateStoreService
	StoreUpdate                *storeService.UpdateStoreService
	StylistCreate              *stylistService.CreateMyStylistService
	StylistUpdate              *stylistService.UpdateMyStylistService
	ScheduleCreateBulk         *scheduleService.CreateSchedulesBulkService
	ScheduleDeleteBulk         *scheduleService.DeleteSchedulesBulkService
	ScheduleCreateTimeSlot     *scheduleService.CreateTimeSlotService
	ScheduleUpdateTimeSlot     *scheduleService.UpdateTimeSlotService
	ScheduleDeleteTimeSlot     *scheduleService.DeleteTimeSlotService
	TimeSlotTemplateCreate     *timeSlotTemplateService.CreateTimeSlotTemplateService
	TimeSlotTemplateUpdate     *timeSlotTemplateService.UpdateTimeSlotTemplateService
	TimeSlotTemplateDelete     *timeSlotTemplateService.DeleteTimeSlotTemplateService
	TimeSlotTemplateCreateItem *timeSlotTemplateService.CreateTimeSlotTemplateItemService
	TimeSlotTemplateUpdateItem *timeSlotTemplateService.UpdateTimeSlotTemplateItemService
	TimeSlotTemplateDeleteItem *timeSlotTemplateService.DeleteTimeSlotTemplateItemService
	ServiceCreate              *serviceService.CreateServiceService
	ServiceUpdate              *serviceService.UpdateServiceService
}

type Handlers struct {
	AuthStaffLogin             *authHandler.StaffLoginHandler
	AuthCustomerLineLogin      *authHandler.CustomerLineLoginHandler
	AuthCustomerLineRegister   *authHandler.CustomerLineRegisterHandler
	BookingCreateMy            *bookingHandler.CreateMyBookingHandler
	BookingUpdateMy            *bookingHandler.UpdateMyBookingHandler
	BookingCancelMy            *bookingHandler.CancelMyBookingHandler
	CustomerUpdateMy           *customerHandler.UpdateMyCustomerHandler
	StaffCreate                *staffHandler.CreateStaffHandler
	StaffUpdate                *staffHandler.UpdateStaffHandler
	StaffUpdateMe              *staffHandler.UpdateMyStaffHandler
	StaffStoreAccess           *staffHandler.CreateStoreAccessHandler
	StaffDeleteStoreAccess     *staffHandler.DeleteStoreAccessBulkHandler
	StoreCreate                *storeHandler.CreateStoreHandler
	StoreUpdate                *storeHandler.UpdateStoreHandler
	StylistCreate              *stylistHandler.CreateMyStylistHandler
	StylistUpdate              *stylistHandler.UpdateMyStylistHandler
	ScheduleCreateBulk         *scheduleHandler.CreateSchedulesBulkHandler
	ScheduleDeleteBulk         *scheduleHandler.DeleteSchedulesBulkHandler
	ScheduleCreateTimeSlot     *scheduleHandler.CreateTimeSlotHandler
	ScheduleUpdateTimeSlot     *scheduleHandler.UpdateTimeSlotHandler
	ScheduleDeleteTimeSlot     *scheduleHandler.DeleteTimeSlotHandler
	TimeSlotTemplateCreate     *timeSlotTemplateHandler.CreateTimeSlotTemplateHandler
	TimeSlotTemplateUpdate     *timeSlotTemplateHandler.UpdateTimeSlotTemplateHandler
	TimeSlotTemplateDelete     *timeSlotTemplateHandler.DeleteTimeSlotTemplateHandler
	TimeSlotTemplateCreateItem *timeSlotTemplateHandler.CreateTimeSlotTemplateItemHandler
	TimeSlotTemplateUpdateItem *timeSlotTemplateHandler.UpdateTimeSlotTemplateItemHandler
	TimeSlotTemplateDeleteItem *timeSlotTemplateHandler.DeleteTimeSlotTemplateItemHandler
	ServiceCreate              *serviceHandler.CreateServiceHandler
	ServiceUpdate              *serviceHandler.UpdateServiceHandler
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
		AuthStaffLogin:             authService.NewStaffLoginService(queries, cfg.JWT),
		AuthCustomerLineLogin:      authService.NewCustomerLineLoginService(queries, cfg.Line, cfg.JWT),
		AuthCustomerLineRegister:   authService.NewCustomerLineRegisterService(queries, database.PgxPool, cfg.Line, cfg.JWT),
		BookingCreateMy:            bookingService.NewCreateMyBookingService(queries, database.PgxPool),
		BookingUpdateMy:            bookingService.NewUpdateMyBookingService(queries, repositories.Booking, database.PgxPool),
		BookingCancelMy:            bookingService.NewCancelMyBookingService(queries),
		CustomerUpdateMy:           customerService.NewUpdateMyCustomerService(queries, repositories.Customer),
		StaffCreate:                staffService.NewCreateStaffService(queries, database.PgxPool),
		StaffUpdate:                staffService.NewUpdateStaffService(queries, database.Sqlx),
		StaffUpdateMe:              staffService.NewUpdateMyStaffService(queries, repositories.StaffUser),
		StaffStoreAccess:           staffService.NewCreateStoreAccessService(queries),
		StaffDeleteStoreAccess:     staffService.NewDeleteStoreAccessBulkService(queries),
		StoreCreate:                storeService.NewCreateStoreService(queries, database.PgxPool),
		StoreUpdate:                storeService.NewUpdateStoreService(queries, repositories.Store),
		StylistCreate:              stylistService.NewCreateMyStylistService(queries),
		StylistUpdate:              stylistService.NewUpdateMyStylistService(queries, repositories.Stylist),
		ScheduleCreateBulk:         scheduleService.NewCreateSchedulesBulkService(queries, database.PgxPool),
		ScheduleDeleteBulk:         scheduleService.NewDeleteSchedulesBulkService(queries, database.PgxPool),
		ScheduleCreateTimeSlot:     scheduleService.NewCreateTimeSlotService(queries),
		ScheduleUpdateTimeSlot:     scheduleService.NewUpdateTimeSlotService(queries, repositories.TimeSlot),
		ScheduleDeleteTimeSlot:     scheduleService.NewDeleteTimeSlotService(queries),
		TimeSlotTemplateCreate:     timeSlotTemplateService.NewCreateTimeSlotTemplateService(queries, database.PgxPool),
		TimeSlotTemplateUpdate:     timeSlotTemplateService.NewUpdateTimeSlotTemplateService(queries, repositories.TimeSlotTemplate),
		TimeSlotTemplateDelete:     timeSlotTemplateService.NewDeleteTimeSlotTemplateService(queries),
		TimeSlotTemplateCreateItem: timeSlotTemplateService.NewCreateTimeSlotTemplateItemService(queries),
		TimeSlotTemplateUpdateItem: timeSlotTemplateService.NewUpdateTimeSlotTemplateItemService(queries),
		TimeSlotTemplateDeleteItem: timeSlotTemplateService.NewDeleteTimeSlotTemplateItemService(queries),
		ServiceCreate:              serviceService.NewCreateServiceService(queries),
		ServiceUpdate:              serviceService.NewUpdateServiceService(queries, repositories.Service),
	}

	handlers := Handlers{
		AuthStaffLogin:             authHandler.NewStaffLoginHandler(services.AuthStaffLogin),
		AuthCustomerLineLogin:      authHandler.NewCustomerLineLoginHandler(services.AuthCustomerLineLogin),
		AuthCustomerLineRegister:   authHandler.NewCustomerLineRegisterHandler(services.AuthCustomerLineRegister),
		BookingCreateMy:            bookingHandler.NewCreateMyBookingHandler(services.BookingCreateMy),
		BookingUpdateMy:            bookingHandler.NewUpdateMyBookingHandler(services.BookingUpdateMy),
		BookingCancelMy:            bookingHandler.NewCancelMyBookingHandler(services.BookingCancelMy),
		CustomerUpdateMy:           customerHandler.NewUpdateMyCustomerHandler(services.CustomerUpdateMy),
		StaffCreate:                staffHandler.NewCreateStaffHandler(services.StaffCreate),
		StaffUpdate:                staffHandler.NewUpdateStaffHandler(services.StaffUpdate),
		StaffUpdateMe:              staffHandler.NewUpdateMyStaffHandler(services.StaffUpdateMe),
		StaffStoreAccess:           staffHandler.NewCreateStoreAccessHandler(services.StaffStoreAccess),
		StaffDeleteStoreAccess:     staffHandler.NewDeleteStoreAccessBulkHandler(services.StaffDeleteStoreAccess),
		StoreCreate:                storeHandler.NewCreateStoreHandler(services.StoreCreate),
		StoreUpdate:                storeHandler.NewUpdateStoreHandler(services.StoreUpdate),
		StylistCreate:              stylistHandler.NewCreateMyStylistHandler(services.StylistCreate),
		StylistUpdate:              stylistHandler.NewUpdateMyStylistHandler(services.StylistUpdate),
		ScheduleCreateBulk:         scheduleHandler.NewCreateSchedulesBulkHandler(services.ScheduleCreateBulk),
		ScheduleDeleteBulk:         scheduleHandler.NewDeleteSchedulesBulkHandler(services.ScheduleDeleteBulk),
		ScheduleCreateTimeSlot:     scheduleHandler.NewCreateTimeSlotHandler(services.ScheduleCreateTimeSlot),
		ScheduleUpdateTimeSlot:     scheduleHandler.NewUpdateTimeSlotHandler(services.ScheduleUpdateTimeSlot),
		ScheduleDeleteTimeSlot:     scheduleHandler.NewDeleteTimeSlotHandler(services.ScheduleDeleteTimeSlot),
		TimeSlotTemplateCreate:     timeSlotTemplateHandler.NewCreateTimeSlotTemplateHandler(services.TimeSlotTemplateCreate),
		TimeSlotTemplateUpdate:     timeSlotTemplateHandler.NewUpdateTimeSlotTemplateHandler(services.TimeSlotTemplateUpdate),
		TimeSlotTemplateDelete:     timeSlotTemplateHandler.NewDeleteTimeSlotTemplateHandler(services.TimeSlotTemplateDelete),
		TimeSlotTemplateCreateItem: timeSlotTemplateHandler.NewCreateTimeSlotTemplateItemHandler(services.TimeSlotTemplateCreateItem),
		TimeSlotTemplateUpdateItem: timeSlotTemplateHandler.NewUpdateTimeSlotTemplateItemHandler(services.TimeSlotTemplateUpdateItem),
		TimeSlotTemplateDeleteItem: timeSlotTemplateHandler.NewDeleteTimeSlotTemplateItemHandler(services.TimeSlotTemplateDeleteItem),
		ServiceCreate:              serviceHandler.NewCreateServiceHandler(services.ServiceCreate),
		ServiceUpdate:              serviceHandler.NewUpdateServiceHandler(services.ServiceUpdate),
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
