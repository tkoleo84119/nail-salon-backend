package app

import (
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	"github.com/tkoleo84119/nail-salon-backend/internal/infra/db"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"

	// Public handlers
	authHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/auth"
	bookingHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/booking"
	customerHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/customer"
	scheduleHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/schedule"
	storeHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/store"

	// Public services
	authService "github.com/tkoleo84119/nail-salon-backend/internal/service/auth"
	bookingService "github.com/tkoleo84119/nail-salon-backend/internal/service/booking"
	customerService "github.com/tkoleo84119/nail-salon-backend/internal/service/customer"
	scheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/schedule"
	storeService "github.com/tkoleo84119/nail-salon-backend/internal/service/store"
)

// PublicServices contains all public/customer-facing services
type PublicServices struct {
	// Authentication services
	AuthCustomerLineLogin    authService.CustomerLineLoginServiceInterface
	AuthCustomerLineRegister authService.CustomerLineRegisterServiceInterface
	AuthRefreshToken         authService.RefreshTokenServiceInterface

	// Customer services
	CustomerGetMy    *customerService.GetMyCustomerService
	CustomerUpdateMy *customerService.UpdateMyCustomerService

	// Booking services
	BookingCreateMy    *bookingService.CreateMyBookingService
	BookingUpdateMy    *bookingService.UpdateMyBookingService
	BookingCancelMy    bookingService.CancelMyBookingServiceInterface
	BookingGetMy       bookingService.GetMyBookingsServiceInterface
	BookingGetMySingle bookingService.GetMyBookingServiceInterface

	// Schedule services
	ScheduleGetStoreSchedules scheduleService.GetStoreScheduleServiceInterface
	ScheduleGetTimeSlots      scheduleService.GetTimeSlotServiceInterface

	// Store services
	StoreGetServices storeService.GetStoreServicesServiceInterface
	StoreGetStylists storeService.GetStoreStylistsServiceInterface
	StoreGetStores   storeService.GetStoresServiceInterface
	StoreGetStore    storeService.GetStoreServiceInterface
}

// PublicHandlers contains all public/customer-facing handlers
type PublicHandlers struct {
	// Authentication handlers
	AuthCustomerLineLogin    *authHandler.CustomerLineLoginHandler
	AuthCustomerLineRegister *authHandler.CustomerLineRegisterHandler
	AuthRefreshToken         *authHandler.RefreshTokenHandler

	// Customer handlers
	CustomerGetMy    *customerHandler.GetMyCustomerHandler
	CustomerUpdateMy *customerHandler.UpdateMyCustomerHandler

	// Booking handlers
	BookingCreateMy    *bookingHandler.CreateMyBookingHandler
	BookingUpdateMy    *bookingHandler.UpdateMyBookingHandler
	BookingCancelMy    *bookingHandler.CancelMyBookingHandler
	BookingGetMy       *bookingHandler.GetMyBookingsHandler
	BookingGetMySingle *bookingHandler.GetMyBookingHandler

	// Schedule handlers
	ScheduleGetStoreSchedules *scheduleHandler.GetStoreScheduleHandler
	ScheduleGetTimeSlots      *scheduleHandler.GetTimeSlotHandler

	// Store handlers
	StoreGetServices *storeHandler.GetStoreServicesHandler
	StoreGetStylists *storeHandler.GetStoreStylistsHandler
	StoreGetStores   *storeHandler.GetStoresHandler
	StoreGetStore    *storeHandler.GetStoreHandler
}

// NewPublicServices creates and initializes all public services
func NewPublicServices(queries *dbgen.Queries, database *db.Database, repositories Repositories, cfg *config.Config) PublicServices {
	return PublicServices{
		// Authentication services
		AuthCustomerLineLogin:    authService.NewCustomerLineLoginService(queries, cfg.Line, cfg.JWT),
		AuthCustomerLineRegister: authService.NewCustomerLineRegisterService(queries, database.PgxPool, cfg.Line, cfg.JWT),
		AuthRefreshToken:         authService.NewRefreshTokenService(queries, cfg.JWT),

		// Customer services
		CustomerGetMy:    customerService.NewGetMyCustomerService(queries),
		CustomerUpdateMy: customerService.NewUpdateMyCustomerService(queries, repositories.Customer),

		// Booking services
		BookingCreateMy:    bookingService.NewCreateMyBookingService(queries, database.PgxPool),
		BookingUpdateMy:    bookingService.NewUpdateMyBookingService(queries, repositories.Booking, database.PgxPool),
		BookingCancelMy:    bookingService.NewCancelMyBookingService(queries),
		BookingGetMy:       bookingService.NewGetMyBookingsService(repositories.Booking),
		BookingGetMySingle: bookingService.NewGetMyBookingService(queries),

		// Schedule services
		ScheduleGetStoreSchedules: scheduleService.NewGetStoreScheduleService(queries),
		ScheduleGetTimeSlots:      scheduleService.NewGetTimeSlotService(queries),

		// Store services
		StoreGetServices: storeService.NewGetStoreServicesService(queries, repositories.Service),
		StoreGetStylists: storeService.NewGetStoreStylistsService(queries, repositories.Stylist),
		StoreGetStores:   storeService.NewGetStoresService(repositories.Store),
		StoreGetStore:    storeService.NewGetStoreService(queries),
	}
}

// NewPublicHandlers creates and initializes all public handlers
func NewPublicHandlers(services PublicServices) PublicHandlers {
	return PublicHandlers{
		// Authentication handlers
		AuthCustomerLineLogin:    authHandler.NewCustomerLineLoginHandler(services.AuthCustomerLineLogin),
		AuthCustomerLineRegister: authHandler.NewCustomerLineRegisterHandler(services.AuthCustomerLineRegister),
		AuthRefreshToken:         authHandler.NewRefreshTokenHandler(services.AuthRefreshToken),

		// Customer handlers
		CustomerGetMy:    customerHandler.NewGetMyCustomerHandler(services.CustomerGetMy),
		CustomerUpdateMy: customerHandler.NewUpdateMyCustomerHandler(services.CustomerUpdateMy),

		// Booking handlers
		BookingCreateMy:    bookingHandler.NewCreateMyBookingHandler(services.BookingCreateMy),
		BookingUpdateMy:    bookingHandler.NewUpdateMyBookingHandler(services.BookingUpdateMy),
		BookingCancelMy:    bookingHandler.NewCancelMyBookingHandler(services.BookingCancelMy),
		BookingGetMy:       bookingHandler.NewGetMyBookingsHandler(services.BookingGetMy),
		BookingGetMySingle: bookingHandler.NewGetMyBookingHandler(services.BookingGetMySingle),

		// Schedule handlers
		ScheduleGetStoreSchedules: scheduleHandler.NewGetStoreScheduleHandler(services.ScheduleGetStoreSchedules),
		ScheduleGetTimeSlots:      scheduleHandler.NewGetTimeSlotHandler(services.ScheduleGetTimeSlots),

		// Store handlers
		StoreGetServices: storeHandler.NewGetStoreServicesHandler(services.StoreGetServices),
		StoreGetStylists: storeHandler.NewGetStoreStylistsHandler(services.StoreGetStylists),
		StoreGetStores:   storeHandler.NewGetStoresHandler(services.StoreGetStores),
		StoreGetStore:    storeHandler.NewGetStoreHandler(services.StoreGetStore),
	}
}
