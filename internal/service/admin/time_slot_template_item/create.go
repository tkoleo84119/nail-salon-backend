package adminTimeSlotTemplateItem

import (
	"context"
	"time"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminTimeSlotTemplateItemModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time_slot_template_item"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	queries *dbgen.Queries
}

func NewCreate(queries *dbgen.Queries) *Create {
	return &Create{
		queries: queries,
	}
}

func (s *Create) Create(ctx context.Context, templateID int64, req adminTimeSlotTemplateItemModel.CreateRequest) (*adminTimeSlotTemplateItemModel.CreateResponse, error) {
	// Validate time format and range
	startTime, err := utils.TimeStringToTime(req.StartTime)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValFieldTimeFormat, "invalid start time format", err)
	}

	endTime, err := utils.TimeStringToTime(req.EndTime)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValFieldTimeFormat, "invalid end time format", err)
	}

	if !endTime.After(startTime) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotEndBeforeStart)
	}

	// Check if template exists
	_, err = s.queries.CheckTimeSlotTemplateExists(ctx, templateID)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotTemplateNotFound)
	}

	// Get existing template items to check for conflicts
	existingItems, err := s.queries.GetTimeSlotTemplateItemsByTemplateID(ctx, templateID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get existing template items", err)
	}

	// Check for time conflicts with existing items
	if err := s.checkTimeConflicts(startTime, endTime, existingItems); err != nil {
		return nil, err
	}

	// Generate item ID
	itemID := utils.GenerateID()

	// Create time slot template item
	item, err := s.queries.CreateTimeSlotTemplateItem(ctx, dbgen.CreateTimeSlotTemplateItemParams{
		ID:         itemID,
		TemplateID: templateID,
		StartTime:  utils.TimePtrToPgTime(&startTime),
		EndTime:    utils.TimePtrToPgTime(&endTime),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create time slot template item", err)
	}

	return &adminTimeSlotTemplateItemModel.CreateResponse{
		ID:         utils.FormatID(item.ID),
		TemplateID: utils.FormatID(item.TemplateID),
		StartTime:  utils.PgTimeToTimeString(item.StartTime),
		EndTime:    utils.PgTimeToTimeString(item.EndTime),
	}, nil
}

// checkTimeConflicts validates that the new time slot doesn't overlap with existing items
func (s *Create) checkTimeConflicts(startTime, endTime time.Time, existingItems []dbgen.GetTimeSlotTemplateItemsByTemplateIDRow) error {
	for _, item := range existingItems {
		// Convert pgtype.Time to time.Time for comparison
		existingStart, err := utils.PgTimeToTime(item.StartTime)
		if err != nil {
			return errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert time", err)
		}
		existingEnd, err := utils.PgTimeToTime(item.EndTime)
		if err != nil {
			return errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert time", err)
		}

		// Check if new slot overlaps with existing slot
		if startTime.Before(existingEnd) && endTime.After(existingStart) {
			return errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotConflict)
		}
	}
	return nil
}
