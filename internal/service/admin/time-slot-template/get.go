package adminTimeSlotTemplate

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetTimeSlotTemplateService struct {
	queries *dbgen.Queries
}

func NewGetTimeSlotTemplateService(queries *dbgen.Queries) *GetTimeSlotTemplateService {
	return &GetTimeSlotTemplateService{
		queries: queries,
	}
}

func (s *GetTimeSlotTemplateService) GetTimeSlotTemplate(ctx context.Context, templateID int64) (*adminTimeSlotTemplateModel.GetTimeSlotTemplateResponse, error) {
	// Get time slot template by ID
	rows, err := s.queries.GetTimeSlotTemplateWithItemsByID(ctx, templateID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get time slot template", err)
	}
	if len(rows) == 0 {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotTemplateNotFound)
	}

	// Build time slot template items response
	response := &adminTimeSlotTemplateModel.GetTimeSlotTemplateResponse{
		Items: make([]adminTimeSlotTemplateModel.GetTimeSlotTemplateItemInfo, 0, len(rows)),
	}

	for i, row := range rows {
		if i == 0 {
			response.ID = utils.FormatID(row.ID)
			response.Name = row.Name
			response.Note = utils.PgTextToString(row.Note)
			response.Updater = utils.PgInt8ToIDString(row.Updater)
			response.CreatedAt = utils.PgTimestamptzToTimeString(row.CreatedAt)
			response.UpdatedAt = utils.PgTimestamptzToTimeString(row.UpdatedAt)
		}

		if row.ItemID.Valid {
			response.Items = append(response.Items, adminTimeSlotTemplateModel.GetTimeSlotTemplateItemInfo{
				ID:        utils.FormatID(row.ItemID.Int64),
				StartTime: utils.PgTimeToTimeString(row.StartTime),
				EndTime:   utils.PgTimeToTimeString(row.EndTime),
			})
		}
	}

	return response, nil
}
