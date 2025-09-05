package adminTimeSlotTemplateItem

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminTimeSlotTemplateItemModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time_slot_template_item"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Delete struct {
	queries *dbgen.Queries
}

func NewDelete(queries *dbgen.Queries) *Delete {
	return &Delete{
		queries: queries,
	}
}

func (s *Delete) Delete(ctx context.Context, templateID int64, itemID int64) (*adminTimeSlotTemplateItemModel.DeleteResponse, error) {
	// Check if template exists
	exists, err := s.queries.CheckTimeSlotTemplateExists(ctx, templateID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check if template exists", err)
	}
	if !exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotTemplateNotFound)
	}

	// Check if item exists and belongs to the template
	exists, err = s.queries.CheckTimeSlotTemplateItemExistsByIDAndTemplateID(ctx, dbgen.CheckTimeSlotTemplateItemExistsByIDAndTemplateIDParams{
		ID:         itemID,
		TemplateID: templateID,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check if item exists and belongs to the template", err)
	}
	if !exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotTemplateItemNotFound)
	}

	// Delete the time slot template item
	err = s.queries.DeleteTimeSlotTemplateItem(ctx, itemID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete time slot template item", err)
	}

	return &adminTimeSlotTemplateItemModel.DeleteResponse{
		Deleted: utils.FormatID(itemID),
	}, nil
}
