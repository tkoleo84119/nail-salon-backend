package adminTimeSlotTemplate

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	queries *dbgen.Queries
	repo    *sqlx.Repositories
}

func NewUpdate(queries *dbgen.Queries, repo *sqlx.Repositories) *Update {
	return &Update{
		queries: queries,
		repo:    repo,
	}
}

func (s *Update) Update(ctx context.Context, templateID int64, req adminTimeSlotTemplateModel.UpdateRequest) (*adminTimeSlotTemplateModel.UpdateResponse, error) {
	// Check if template exists
	exist, err := s.queries.CheckTimeSlotTemplateExists(ctx, templateID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check time slot template exists", err)
	}
	if !exist {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotTemplateNotFound)
	}

	// Update the template using sqlx repository
	response, err := s.repo.Template.UpdateTimeSlotTemplate(ctx, templateID, sqlx.UpdateTimeSlotTemplateParams{
		Name: req.Name,
		Note: req.Note,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update time slot template", err)
	}

	return &adminTimeSlotTemplateModel.UpdateResponse{
		ID:        utils.FormatID(response.ID),
		Name:      response.Name,
		Note:      utils.PgTextToString(response.Note),
		Updater:   utils.PgInt8ToIDString(response.Updater),
		CreatedAt: utils.PgTimestamptzToTimeString(response.CreatedAt),
		UpdatedAt: utils.PgTimestamptzToTimeString(response.UpdatedAt),
	}, nil
}
