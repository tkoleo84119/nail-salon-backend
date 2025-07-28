package app

import (
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	"github.com/tkoleo84119/nail-salon-backend/internal/infra/db"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
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
	Schedule         *sqlx.ScheduleRepository
}

type Services struct {
	Public PublicServices
	Admin  AdminServices
}

type Handlers struct {
	Public PublicHandlers
	Admin  AdminHandlers
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
		Schedule:         sqlx.NewScheduleRepository(database.Sqlx),
	}

	// Initialize services using separated containers
	publicServices := NewPublicServices(queries, database, repositories, cfg)
	adminServices := NewAdminServices(queries, database, repositories, cfg)

	services := Services{
		Public: publicServices,
		Admin:  adminServices,
	}

	// Initialize handlers using separated containers
	publicHandlers := NewPublicHandlers(publicServices)
	adminHandlers := NewAdminHandlers(adminServices)

	handlers := Handlers{
		Public: publicHandlers,
		Admin:  adminHandlers,
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
