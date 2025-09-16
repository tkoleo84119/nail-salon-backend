package job

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	"github.com/tkoleo84119/nail-salon-backend/internal/infra/redis"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
)

const (
	RefreshRevokeJobLockKey = "refresh_revoke_job_lock"
	RefreshRevokeLockTTL    = 30 * time.Minute
	RefreshRevokeBatchSize  = 200
)

type RefreshRevokeJob struct {
	cfg            *config.Config
	queries        *dbgen.Queries
	redisClient    *redis.Client
	cron           *cron.Cron
	taiwanLocation *time.Location
}

func NewRefreshRevokeJob(cfg *config.Config, queries *dbgen.Queries, redisClient *redis.Client) (*RefreshRevokeJob, error) {
	taiwanLocation, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return nil, fmt.Errorf("failed to load Taiwan timezone: %w", err)
	}

	c := cron.New(cron.WithLocation(taiwanLocation))

	return &RefreshRevokeJob{
		cfg:            cfg,
		queries:        queries,
		redisClient:    redisClient,
		cron:           c,
		taiwanLocation: taiwanLocation,
	}, nil
}

func (j *RefreshRevokeJob) Start() error {
	_, err := j.cron.AddFunc(j.cfg.Scheduler.RefreshRevokeCron, j.executeRefreshRevokeJob)
	if err != nil {
		return fmt.Errorf("failed to schedule refresh revoke job: %w", err)
	}

	j.cron.Start()
	log.Printf("Refresh revoke job started with schedule: %s (Taiwan timezone)", j.cfg.Scheduler.RefreshRevokeCron)

	return nil
}

func (j *RefreshRevokeJob) Stop() {
	j.cron.Stop()
	log.Println("Refresh revoke job stopped")
}

func (j *RefreshRevokeJob) executeRefreshRevokeJob() {
	ctx := context.Background()

	lockAcquired, err := j.redisClient.SetLock(ctx, RefreshRevokeJobLockKey, "locked", RefreshRevokeLockTTL)
	if err != nil {
		log.Printf("Failed to acquire lock for refresh revoke job: %v", err)
		return
	}

	if !lockAcquired {
		log.Println("Another instance is already running refresh revoke job, skipping...")
		return
	}

	defer func() {
		if err := j.redisClient.ReleaseLock(ctx, RefreshRevokeJobLockKey); err != nil {
			log.Printf("Failed to release lock for refresh revoke job: %v", err)
		}
	}()

	if err := j.processStaffUserRefreshTokenRevoke(ctx); err != nil {
		log.Printf("failed to process staff user refresh token revoke: %v", err)
		return
	}

	if err := j.processCustomerRefreshTokenRevoke(ctx); err != nil {
		log.Printf("failed to process customer refresh token revoke: %v", err)
		return
	}

	log.Println("Refresh revoke job execution completed successfully")
}

func (j *RefreshRevokeJob) processStaffUserRefreshTokenRevoke(ctx context.Context) error {
	count, err := j.queries.CountExpiredOrRevokedStaffUserTokens(ctx)
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	for i := 0; i < int(count); i += RefreshRevokeBatchSize {
		batchSize := RefreshRevokeBatchSize
		if remaining := int(count) - i; remaining < RefreshRevokeBatchSize {
			batchSize = remaining
		}

		err := j.queries.DeleteStaffUserTokensBatch(ctx, int32(batchSize))
		if err != nil {
			return err
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (j *RefreshRevokeJob) processCustomerRefreshTokenRevoke(ctx context.Context) error {
	count, err := j.queries.CountExpiredOrRevokedCustomerTokens(ctx)
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	for i := 0; i < int(count); i += RefreshRevokeBatchSize {
		batchSize := RefreshRevokeBatchSize
		if remaining := int(count) - i; remaining < RefreshRevokeBatchSize {
			batchSize = remaining
		}

		err := j.queries.DeleteCustomerTokensBatch(ctx, int32(batchSize))
		if err != nil {
			return err
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}
