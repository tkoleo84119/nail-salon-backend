package timeSlotTemplate

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	timeSlotTemplate "github.com/tkoleo84119/nail-salon-backend/internal/model/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateTimeSlotTemplateService struct {
	queries dbgen.Querier
	db      *pgxpool.Pool
}

func NewCreateTimeSlotTemplateService(queries dbgen.Querier, db *pgxpool.Pool) *CreateTimeSlotTemplateService {
	return &CreateTimeSlotTemplateService{
		queries: queries,
		db:      db,
	}
}

func (s *CreateTimeSlotTemplateService) CreateTimeSlotTemplate(ctx context.Context, req timeSlotTemplate.CreateTimeSlotTemplateRequest, staffContext common.StaffContext) (*timeSlotTemplate.CreateTimeSlotTemplateResponse, error) {
	// Parse staff user ID
	staffUserID, err := utils.ParseID(staffContext.UserID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.AuthStaffFailed, "invalid staff user ID", err)
	}

	// Validate time slots
	if err := s.validateTimeSlots(req.TimeSlots); err != nil {
		return nil, err
	}

	// Generate template ID
	templateID := utils.GenerateID()

	// Prepare template items for batch insert
	templateItems := make([]dbgen.BatchCreateTimeSlotTemplateItemsParams, len(req.TimeSlots))
	responseTimeSlots := make([]timeSlotTemplate.TimeSlotItemResponse, len(req.TimeSlots))

	for i, timeSlot := range req.TimeSlots {
		// Parse time strings
		startTime, err := common.ParseTimeSlot(timeSlot.StartTime)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid start time format", err)
		}

		endTime, err := common.ParseTimeSlot(timeSlot.EndTime)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid end time format", err)
		}

		itemID := utils.GenerateID()

		templateItems[i] = dbgen.BatchCreateTimeSlotTemplateItemsParams{
			ID:         itemID,
			TemplateID: templateID,
			StartTime:  pgtype.Time{Microseconds: int64(startTime.Hour()*3600+startTime.Minute()*60+startTime.Second()) * 1000000, Valid: true},
			EndTime:    pgtype.Time{Microseconds: int64(endTime.Hour()*3600+endTime.Minute()*60+endTime.Second()) * 1000000, Valid: true},
			CreatedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
			UpdatedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		responseTimeSlots[i] = timeSlotTemplate.TimeSlotItemResponse{
			ID:        utils.FormatID(itemID),
			StartTime: timeSlot.StartTime,
			EndTime:   timeSlot.EndTime,
		}
	}

	// Begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := dbgen.New(tx)

	// Create time slot template
	template, err := qtx.CreateTimeSlotTemplate(ctx, dbgen.CreateTimeSlotTemplateParams{
		ID:      templateID,
		Name:    req.Name,
		Note:    pgtype.Text{String: req.Note, Valid: req.Note != ""},
		Updater: pgtype.Int8{Int64: staffUserID, Valid: true},
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create time slot template", err)
	}

	// Batch create template items
	if _, err := qtx.BatchCreateTimeSlotTemplateItems(ctx, templateItems); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create template items", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	return &timeSlotTemplate.CreateTimeSlotTemplateResponse{
		ID:        utils.FormatID(template.ID),
		Name:      template.Name,
		Note:      template.Note.String,
		TimeSlots: responseTimeSlots,
	}, nil
}

// validateTimeSlots validates that time slots don't overlap and have valid time ranges
func (s *CreateTimeSlotTemplateService) validateTimeSlots(timeSlots []timeSlotTemplate.TimeSlotItem) error {
	if len(timeSlots) == 0 {
		return errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "at least one time slot is required", nil)
	}

	// Parse and validate each time slot
	parsedTimeSlots := make([]struct {
		start time.Time
		end   time.Time
	}, len(timeSlots))

	for i, timeSlot := range timeSlots {
		startTime, err := common.ParseTimeSlot(timeSlot.StartTime)
		if err != nil {
			return errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid start time format", err)
		}

		endTime, err := common.ParseTimeSlot(timeSlot.EndTime)
		if err != nil {
			return errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid end time format", err)
		}

		// Validate time range
		if !endTime.After(startTime) {
			return errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "end time must be after start time", nil)
		}

		parsedTimeSlots[i] = struct {
			start time.Time
			end   time.Time
		}{startTime, endTime}
	}

	// Check for overlaps
	for i := 0; i < len(parsedTimeSlots); i++ {
		for j := i + 1; j < len(parsedTimeSlots); j++ {
			slot1 := parsedTimeSlots[i]
			slot2 := parsedTimeSlots[j]

			// Check if slots overlap
			if slot1.start.Before(slot2.end) && slot2.start.Before(slot1.end) {
				return errorCodes.NewServiceError(errorCodes.ScheduleTimeConflict, "time slots cannot overlap", nil)
			}
		}
	}

	return nil
}
