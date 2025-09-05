package adminTimeSlotTemplateItem

import (
	"context"
	"time"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminTimeSlotTemplateItemModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time_slot_template_item"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	queries *dbgen.Queries
}

func NewUpdate(queries *dbgen.Queries) *Update {
	return &Update{
		queries: queries,
	}
}

func (s *Update) Update(ctx context.Context, templateID int64, itemID int64, req adminTimeSlotTemplateItemModel.UpdateRequest) (*adminTimeSlotTemplateItemModel.UpdateResponse, error) {
	// Validate time format and range
	startTime, err := utils.TimeStringToTime(req.StartTime)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValFieldTimeFormat, "invalid start time format", err)
	}

	endTime, err := utils.TimeStringToTime(req.EndTime)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValFieldTimeFormat, "invalid end time format", err)
	}

	// Validate time range
	if !endTime.After(startTime) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotEndBeforeStart)
	}

	// Check if template exists
	_, err = s.queries.CheckTimeSlotTemplateExists(ctx, templateID)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotTemplateNotFound)
	}

	// Check if item exists and belongs to the template
	exists, err := s.queries.CheckTimeSlotTemplateItemExistsByIDAndTemplateID(ctx, dbgen.CheckTimeSlotTemplateItemExistsByIDAndTemplateIDParams{
		ID:         itemID,
		TemplateID: templateID,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check if item exists and belongs to the template", err)
	}
	if !exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotTemplateItemNotFound)
	}

	// Get other existing template items to check for conflicts (excluding current item)
	otherItems, err := s.queries.GetTimeSlotTemplateItemsByTemplateIDExcluding(ctx, dbgen.GetTimeSlotTemplateItemsByTemplateIDExcludingParams{
		TemplateID: templateID,
		ID:         itemID,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get existing template items", err)
	}

	// Check for time conflicts with other existing items
	if err := s.checkTimeConflicts(startTime, endTime, otherItems); err != nil {
		return nil, err
	}

	// Update time slot template item
	updatedItem, err := s.queries.UpdateTimeSlotTemplateItem(ctx, dbgen.UpdateTimeSlotTemplateItemParams{
		ID:        itemID,
		StartTime: utils.TimeToPgTime(startTime),
		EndTime:   utils.TimeToPgTime(endTime),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update time slot template item", err)
	}

	return &adminTimeSlotTemplateItemModel.UpdateResponse{
		ID:         utils.FormatID(updatedItem.ID),
		TemplateID: utils.FormatID(updatedItem.TemplateID),
		StartTime:  utils.PgTimeToTimeString(updatedItem.StartTime),
		EndTime:    utils.PgTimeToTimeString(updatedItem.EndTime),
	}, nil
}

// checkTimeConflicts validates that the updated time slot doesn't overlap with other existing items
func (s *Update) checkTimeConflicts(startTime, endTime time.Time, existingItems []dbgen.GetTimeSlotTemplateItemsByTemplateIDExcludingRow) error {
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
