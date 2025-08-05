package adminTimeSlotTemplate

import (
	"context"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	queries dbgen.Querier
	db      *pgxpool.Pool
}

func NewCreate(queries dbgen.Querier, db *pgxpool.Pool) *Create {
	return &Create{
		queries: queries,
		db:      db,
	}
}

func (s *Create) Create(ctx context.Context, req adminTimeSlotTemplateModel.CreateRequest, creatorID int64) (*adminTimeSlotTemplateModel.CreateResponse, error) {
	// Validate time slots
	templateID, templateItems, responseTimeSlots, err := s.validateTimeSlotsAndPrepareData(req.TimeSlots)
	if err != nil {
		return nil, err
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
		Note:    utils.StringPtrToPgText(&req.Note, true),
		Updater: utils.Int64ToPgInt8(creatorID),
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

	return &adminTimeSlotTemplateModel.CreateResponse{
		ID:        utils.FormatID(template.ID),
		Name:      template.Name,
		Note:      template.Note.String,
		TimeSlots: responseTimeSlots,
	}, nil
}

// validateTimeSlots validates that time slots don't overlap and have valid time ranges
func (s *Create) validateTimeSlotsAndPrepareData(timeSlots []adminTimeSlotTemplateModel.CreateTimeSlotItem) (int64, []dbgen.BatchCreateTimeSlotTemplateItemsParams, []adminTimeSlotTemplateModel.CreateTimeSlotItemResponse, error) {
	// Generate template ID
	templateID := utils.GenerateID()

	// Prepare template items for batch insert
	templateItems := make([]dbgen.BatchCreateTimeSlotTemplateItemsParams, len(timeSlots))
	responseTimeSlots := make([]adminTimeSlotTemplateModel.CreateTimeSlotItemResponse, len(timeSlots))

	// Parse and validate each time slot
	parsedTimeSlots := make([]struct {
		start time.Time
		end   time.Time
	}, len(timeSlots))

	for i, timeSlot := range timeSlots {
		startTime, err := utils.TimeStringToTime(timeSlot.StartTime)
		if err != nil {
			return 0, nil, nil, errorCodes.NewServiceError(errorCodes.ValFieldTimeFormat, "invalid start time format", err)
		}

		endTime, err := utils.TimeStringToTime(timeSlot.EndTime)
		if err != nil {
			return 0, nil, nil, errorCodes.NewServiceError(errorCodes.ValFieldTimeFormat, "invalid end time format", err)
		}

		// Validate time range
		if !endTime.After(startTime) {
			return 0, nil, nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotEndBeforeStart)
		}

		itemID := utils.GenerateID()
		now := time.Now()
		templateItems[i] = dbgen.BatchCreateTimeSlotTemplateItemsParams{
			ID:         itemID,
			TemplateID: templateID,
			StartTime:  utils.TimeToPgTime(startTime),
			EndTime:    utils.TimeToPgTime(endTime),
			CreatedAt:  utils.TimeToPgTimestamptz(now),
			UpdatedAt:  utils.TimeToPgTimestamptz(now),
		}

		responseTimeSlots[i] = adminTimeSlotTemplateModel.CreateTimeSlotItemResponse{
			ID:        utils.FormatID(itemID),
			StartTime: timeSlot.StartTime,
			EndTime:   timeSlot.EndTime,
		}

		parsedTimeSlots[i] = struct {
			start time.Time
			end   time.Time
		}{startTime, endTime}
	}

	// Sort by start time
	sort.Slice(parsedTimeSlots, func(i, j int) bool {
		return parsedTimeSlots[i].start.Before(parsedTimeSlots[j].start)
	})

	// Check for overlaps
	for i := 1; i < len(parsedTimeSlots); i++ {
		if parsedTimeSlots[i].start.Before(parsedTimeSlots[i-1].end) {
			return 0, nil, nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotConflict)
		}
	}

	return templateID, templateItems, responseTimeSlots, nil
}
