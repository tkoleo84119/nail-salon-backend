package adminTimeSlotTemplate

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Delete struct {
	queries *dbgen.Queries
}

func NewDelete(queries *dbgen.Queries) DeleteInterface {
	return &Delete{
		queries: queries,
	}
}

func (s *Delete) Delete(ctx context.Context, templateID int64) (*adminTimeSlotTemplateModel.DeleteResponse, error) {
	// Check if template exists
	exists, err := s.queries.CheckTimeSlotTemplateExists(ctx, templateID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check time slot template exists", err)
	}
	if !exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotTemplateNotFound)
	}

	// Delete the time slot template
	err = s.queries.DeleteTimeSlotTemplate(ctx, templateID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete time slot template", err)
	}

	return &adminTimeSlotTemplateModel.DeleteResponse{
		Deleted: utils.FormatID(templateID),
	}, nil
}
