package app

import (
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	"github.com/tkoleo84119/nail-salon-backend/internal/infra/db"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"

	// Public handlers
	authHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/auth"
	bookingHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/booking"
	customerHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/customer"
	customerCouponHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/customer_coupon"
	scheduleHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/schedule"
	serviceHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/service"
	storeHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/store"
	stylistHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/stylist"
	timeSlotHandler "github.com/tkoleo84119/nail-salon-backend/internal/handler/time_slot"

	// Public services
	authService "github.com/tkoleo84119/nail-salon-backend/internal/service/auth"
	bookingService "github.com/tkoleo84119/nail-salon-backend/internal/service/booking"
	customerService "github.com/tkoleo84119/nail-salon-backend/internal/service/customer"
	customerCouponService "github.com/tkoleo84119/nail-salon-backend/internal/service/customer_coupon"
	scheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/schedule"
	serviceService "github.com/tkoleo84119/nail-salon-backend/internal/service/service"
	storeService "github.com/tkoleo84119/nail-salon-backend/internal/service/store"
	stylistService "github.com/tkoleo84119/nail-salon-backend/internal/service/stylist"
	timeSlotService "github.com/tkoleo84119/nail-salon-backend/internal/service/time_slot"
)

// PublicServices contains all public/customer-facing services
type PublicServices struct {
	// Authentication services
	AuthLineLogin    authService.LineLoginInterface
	AuthLineRegister authService.LineRegisterInterface
	AuthRefreshToken authService.RefreshTokenInterface

	// Customer services
	CustomerGetMe    customerService.GetMeInterface
	CustomerUpdateMe customerService.UpdateMeInterface

	// CustomerCoupon services
	CustomerCouponGetAll customerCouponService.GetAllInterface

	// Booking services
	BookingCreate      bookingService.CreateInterface
	BookingUpdate      bookingService.UpdateInterface
	BookingCancel      bookingService.CancelInterface
	BookingGetAll      bookingService.GetAllInterface
	BookingGetMySingle bookingService.GetInterface

	// Schedule services
	ScheduleGetAll scheduleService.GetAllInterface

	// TimeSlot services
	TimeSlotGetAll timeSlotService.GetAllInterface

	// Store services
	StoreGetAll storeService.GetAllInterface

	// Service services
	ServiceGetAll serviceService.GetAllInterface

	// Stylist services
	StylistGetAll stylistService.GetAllInterface
}

// PublicHandlers contains all public/customer-facing handlers
type PublicHandlers struct {
	// Authentication handlers
	AuthLineLogin    *authHandler.LineLogin
	AuthLineRegister *authHandler.LineRegister
	AuthRefreshToken *authHandler.RefreshToken

	// Customer handlers
	CustomerGetMe    *customerHandler.GetMe
	CustomerUpdateMe *customerHandler.UpdateMe

	// CustomerCoupon handlers
	CustomerCouponGetAll *customerCouponHandler.GetAll

	// Booking handlers
	BookingCreate      *bookingHandler.Create
	BookingUpdate      *bookingHandler.Update
	BookingCancel      *bookingHandler.Cancel
	BookingGetAll      *bookingHandler.GetAll
	BookingGetMySingle *bookingHandler.Get

	// Schedule handlers
	ScheduleGetAll *scheduleHandler.GetAll

	// TimeSlot handlers
	TimeSlotGetAll *timeSlotHandler.GetAll

	// Store handlers
	StoreGetAll *storeHandler.GetAll

	// Service handlers
	ServiceGetAll *serviceHandler.GetAll

	// Stylist handlers
	StylistGetAll *stylistHandler.GetAll
}

// NewPublicServices creates and initializes all public services
func NewPublicServices(queries *dbgen.Queries, database *db.Database, repositories Repositories, cfg *config.Config) PublicServices {
	return PublicServices{
		// Authentication services
		AuthLineLogin:    authService.NewLineLogin(queries, database.PgxPool, cfg.Line, cfg.JWT),
		AuthLineRegister: authService.NewLineRegister(queries, database.PgxPool, cfg.Line, cfg.JWT),
		AuthRefreshToken: authService.NewRefreshToken(queries, cfg.JWT),

		// Customer services
		CustomerGetMe:    customerService.NewGetMe(queries),
		CustomerUpdateMe: customerService.NewUpdateMe(queries, repositories.SQLX),

		// CustomerCoupon services
		CustomerCouponGetAll: customerCouponService.NewGetAll(queries, repositories.SQLX),

		// Booking services
		BookingCreate:      bookingService.NewCreate(queries, database.PgxPool),
		BookingUpdate:      bookingService.NewUpdate(queries, repositories.SQLX, database.Sqlx),
		BookingCancel:      bookingService.NewCancel(queries, database.PgxPool),
		BookingGetAll:      bookingService.NewGetAll(repositories.SQLX),
		BookingGetMySingle: bookingService.NewGet(queries),

		// Schedule services
		ScheduleGetAll: scheduleService.NewGetAll(queries),

		// TimeSlot services
		TimeSlotGetAll: timeSlotService.NewGetAll(queries),

		// Store services
		StoreGetAll: storeService.NewGetAll(repositories.SQLX),

		// Service services
		ServiceGetAll: serviceService.NewGetAll(repositories.SQLX),

		// Stylist services
		StylistGetAll: stylistService.NewGetAll(queries, repositories.SQLX),
	}
}

// NewPublicHandlers creates and initializes all public handlers
func NewPublicHandlers(services PublicServices) PublicHandlers {
	return PublicHandlers{
		// Authentication handlers
		AuthLineLogin:    authHandler.NewLineLogin(services.AuthLineLogin),
		AuthLineRegister: authHandler.NewLineRegister(services.AuthLineRegister),
		AuthRefreshToken: authHandler.NewRefreshToken(services.AuthRefreshToken),

		// Customer handlers
		CustomerGetMe:    customerHandler.NewGetMe(services.CustomerGetMe),
		CustomerUpdateMe: customerHandler.NewUpdateMe(services.CustomerUpdateMe),

		// CustomerCoupon handlers
		CustomerCouponGetAll: customerCouponHandler.NewGetAll(services.CustomerCouponGetAll),

		// Booking handlers
		BookingCreate:      bookingHandler.NewCreate(services.BookingCreate),
		BookingUpdate:      bookingHandler.NewUpdate(services.BookingUpdate),
		BookingCancel:      bookingHandler.NewCancel(services.BookingCancel),
		BookingGetAll:      bookingHandler.NewGetAll(services.BookingGetAll),
		BookingGetMySingle: bookingHandler.NewGet(services.BookingGetMySingle),

		// Schedule handlers
		ScheduleGetAll: scheduleHandler.NewGetAll(services.ScheduleGetAll),

		// TimeSlot handlers
		TimeSlotGetAll: timeSlotHandler.NewGetAll(services.TimeSlotGetAll),

		// Store handlers
		StoreGetAll: storeHandler.NewGetAll(services.StoreGetAll),

		// Service handlers
		ServiceGetAll: serviceHandler.NewGetAll(services.ServiceGetAll),

		// Stylist handlers
		StylistGetAll: stylistHandler.NewGetAll(services.StylistGetAll),
	}
}
