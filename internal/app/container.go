package app

import (
	"fmt"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	"github.com/tkoleo84119/nail-salon-backend/internal/infra/db"
	"github.com/tkoleo84119/nail-salon-backend/internal/infra/redis"
	"github.com/tkoleo84119/nail-salon-backend/internal/job"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/service/cache"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Container struct {
	cfg           *config.Config
	database      *db.Database
	redis         *redis.Client
	queries       *dbgen.Queries
	lineMessenger *utils.LineMessageClient
	authCache     cache.AuthCacheInterface
	activityLog   cache.ActivityLogCacheInterface

	repositories Repositories
	services     Services
	handlers     Handlers
	jobs         Jobs
}

type Repositories struct {
	SQLX *sqlx.Repositories
}

type Services struct {
	Public PublicServices
	Admin  AdminServices
}

type Handlers struct {
	Public PublicHandlers
	Admin  AdminHandlers
}

type Jobs struct {
	LineReminderJob *job.LineReminderJob
}

func NewContainer(cfg *config.Config, database *db.Database, redisClient *redis.Client) (*Container, error) {
	queries := dbgen.New(database.PgxPool)
	lineMessenger := utils.NewLineMessenger(cfg.Line.MessageAccessToken)
	authCache := cache.NewAuthCache(redisClient)
	activityLog := cache.NewActivityLogCache(redisClient)

	repositories := Repositories{
		SQLX: sqlx.NewRepositories(database.Sqlx),
	}

	// Initialize services using separated containers
	publicServices := NewPublicServices(queries, database, repositories, cfg, lineMessenger, authCache, activityLog)
	adminServices := NewAdminServices(queries, database, repositories, cfg, lineMessenger, authCache, activityLog)

	services := Services{
		Public: publicServices,
		Admin:  adminServices,
	}

	// Initialize handlers using separated containers
    publicHandlers := NewPublicHandlers(publicServices, cfg)
    adminHandlers := NewAdminHandlers(adminServices, cfg)

	handlers := Handlers{
		Public: publicHandlers,
		Admin:  adminHandlers,
	}

	// Initialize jobs
	lineReminderJob, err := job.NewLineReminderJob(cfg, queries, redisClient, lineMessenger)
	if err != nil {
		return nil, fmt.Errorf("failed to create line reminder job: %w", err)
	}

	jobs := Jobs{
		LineReminderJob: lineReminderJob,
	}

	return &Container{
		cfg:           cfg,
		database:      database,
		redis:         redisClient,
		queries:       queries,
		lineMessenger: lineMessenger,
		authCache:     authCache,
		activityLog:   activityLog,
		repositories:  repositories,
		services:      services,
		handlers:      handlers,
		jobs:          jobs,
	}, nil
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

func (c *Container) GetRedis() *redis.Client {
	return c.redis
}

func (c *Container) GetQueries() *dbgen.Queries {
	return c.queries
}

func (c *Container) GetJobs() Jobs {
	return c.jobs
}

func (c *Container) GetLineMessenger() *utils.LineMessageClient {
	return c.lineMessenger
}

func (c *Container) GetAuthCache() cache.AuthCacheInterface {
	return c.authCache
}
