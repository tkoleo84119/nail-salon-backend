package adminStylist

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStylistModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stylist"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateMyStylistService struct {
	repo *sqlxRepo.Repositories
}

func NewUpdateMyStylistService(repo *sqlxRepo.Repositories) *UpdateMyStylistService {
	return &UpdateMyStylistService{
		repo: repo,
	}
}

func (s *UpdateMyStylistService) UpdateMyStylist(ctx context.Context, req adminStylistModel.UpdateMyStylistRequest, staffUserID int64) (*adminStylistModel.UpdateMyStylistResponse, error) {
	// ensure at least one field is provided for update
	if !req.HasUpdate() {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValAllFieldsEmpty)
	}

	// Check if stylist exists for this staff user
	_, err := s.repo.Stylist.GetStylistByStaffUserID(ctx, staffUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get stylist by staff user id", err)
	}

	// Update stylist record using repository
	updateStylist, err := s.repo.Stylist.UpdateStylist(ctx, staffUserID, sqlxRepo.UpdateStylistParams{
		Name:         req.Name,
		GoodAtShapes: req.GoodAtShapes,
		GoodAtColors: req.GoodAtColors,
		GoodAtStyles: req.GoodAtStyles,
		IsIntrovert:  req.IsIntrovert,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update stylist", err)
	}

	response := &adminStylistModel.UpdateMyStylistResponse{
		ID:           utils.FormatID(updateStylist.ID),
		StaffUserID:  utils.FormatID(staffUserID),
		Name:         utils.PgTextToString(updateStylist.Name),
		GoodAtShapes: updateStylist.GoodAtShapes,
		GoodAtColors: updateStylist.GoodAtColors,
		GoodAtStyles: updateStylist.GoodAtStyles,
		IsIntrovert:  utils.PgBoolToBool(updateStylist.IsIntrovert),
		CreatedAt:    utils.PgTimestamptzToTimeString(updateStylist.CreatedAt),
		UpdatedAt:    utils.PgTimestamptzToTimeString(updateStylist.UpdatedAt),
	}

	return response, nil
}
