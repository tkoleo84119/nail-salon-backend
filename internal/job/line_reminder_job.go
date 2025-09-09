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
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

const (
	LineReminderJobLockKey = "line_reminder_job_lock"
	LineReminderLockTTL    = 30 * time.Minute
	BatchSize              = 100
)

type LineReminderJob struct {
	cfg            *config.Config
	queries        *dbgen.Queries
	redisClient    *redis.Client
	lineMessenger  *utils.LineMessageClient
	cron           *cron.Cron
	taiwanLocation *time.Location
}

type BookingReminderData struct {
	ID              int64
	CustomerLineUID string
	CustomerName    string
	StoreName       string
	StoreAddress    string
	WorkDate        string
	StartTime       string
	EndTime         string
}

func NewLineReminderJob(cfg *config.Config, queries *dbgen.Queries, redisClient *redis.Client, lineMessenger *utils.LineMessageClient) (*LineReminderJob, error) {
	taiwanLocation, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return nil, fmt.Errorf("failed to load Taiwan timezone: %w", err)
	}

	c := cron.New(cron.WithLocation(taiwanLocation))

	return &LineReminderJob{
		cfg:            cfg,
		queries:        queries,
		redisClient:    redisClient,
		lineMessenger:  lineMessenger,
		cron:           c,
		taiwanLocation: taiwanLocation,
	}, nil
}

func (j *LineReminderJob) Start() error {
	_, err := j.cron.AddFunc(j.cfg.Scheduler.LineReminderCron, j.executeLineReminderJob)
	if err != nil {
		return fmt.Errorf("failed to schedule reminder job: %w", err)
	}

	j.cron.Start()
	log.Printf("Line reminder job started with schedule: %s (Taiwan timezone)", j.cfg.Scheduler.LineReminderCron)

	return nil
}

func (j *LineReminderJob) Stop() {
	j.cron.Stop()
	log.Println("Line reminder job stopped")
}

func (j *LineReminderJob) executeLineReminderJob() {
	ctx := context.Background()

	lockAcquired, err := j.redisClient.SetLock(ctx, LineReminderJobLockKey, "locked", LineReminderLockTTL)
	if err != nil {
		log.Printf("Failed to acquire lock for reminder job: %v", err)
		return
	}

	if !lockAcquired {
		log.Println("Another instance is already running reminder job, skipping...")
		return
	}

	defer func() {
		if err := j.redisClient.ReleaseLock(ctx, LineReminderJobLockKey); err != nil {
			log.Printf("Failed to release lock for reminder job: %v", err)
		}
	}()

	tomorrow := time.Now().In(j.taiwanLocation).AddDate(0, 0, 1)
	tomorrowDateString := tomorrow.Format("2006-01-02")

	log.Printf("Starting line reminder job execution for date: %s", tomorrowDateString)

	if err := j.processBookingReminders(ctx, tomorrowDateString); err != nil {
		log.Printf("Error processing booking reminders: %v", err)
		return
	}

	log.Println("Line reminder job execution completed successfully")
}

func (j *LineReminderJob) processBookingReminders(ctx context.Context, date string) error {
	pgDate, err := utils.DateStringToPgDate(date)
	if err != nil {
		return fmt.Errorf("failed to convert date string: %w", err)
	}

	// get tomorrow bookings for reminder
	bookings, err := j.queries.GetTomorrowBookingsForReminder(ctx, pgDate)
	if err != nil {
		return fmt.Errorf("failed to get tomorrow bookings: %w", err)
	}

	if len(bookings) == 0 {
		return nil
	}

	// process bookings in batches
	for i := 0; i < len(bookings); i += BatchSize {
		end := i + BatchSize
		if end > len(bookings) {
			end = len(bookings)
		}

		batch := bookings[i:end]
		if err := j.processBatch(ctx, batch); err != nil {
			log.Printf("Error processing batch %d-%d: %v", i, end-1, err)
		}
	}

	return nil
}

func (j *LineReminderJob) processBatch(ctx context.Context, bookings []dbgen.GetTomorrowBookingsForReminderRow) error {
	for _, booking := range bookings {
		reminderData := &BookingReminderData{
			ID:              booking.ID,
			CustomerLineUID: booking.CustomerLineUid,
			CustomerName:    booking.CustomerName,
			StoreName:       booking.StoreName,
			StoreAddress:    utils.PgTextToString(booking.StoreAddress),
			WorkDate:        utils.PgDateToDateString(booking.WorkDate),
			StartTime:       utils.PgTimeToTimeString(booking.StartTime),
			EndTime:         utils.PgTimeToTimeString(booking.EndTime),
		}

		if err := j.sendReminderMessage(reminderData); err != nil {
			log.Printf("Failed to send reminder for booking ID %d: %v", booking.ID, err)
		} else {
			log.Printf("Reminder sent successfully for booking ID %d", booking.ID)
		}
	}

	return nil
}

func (j *LineReminderJob) sendReminderMessage(data *BookingReminderData) error {
	bookingData := &utils.BookingData{
		StoreName:    data.StoreName,
		StoreAddress: data.StoreAddress,
		Date:         data.WorkDate,
		StartTime:    data.StartTime,
		EndTime:      data.EndTime,
		CustomerName: &data.CustomerName,
	}

	return j.lineMessenger.SendBookingReminderMessage(data.CustomerLineUID, bookingData)
}
