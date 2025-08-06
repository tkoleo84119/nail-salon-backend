package adminStaff

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateMe struct {
	queries *dbgen.Queries
	repo    *sqlxRepo.Repositories
}

func NewUpdateMe(queries *dbgen.Queries, repo *sqlxRepo.Repositories) *UpdateMe {
	return &UpdateMe{
		queries: queries,
		repo:    repo,
	}
}

func (s *UpdateMe) UpdateMe(ctx context.Context, req adminStaffModel.UpdateMeRequest, staffUserID int64) (*adminStaffModel.UpdateMeResponse, error) {
	// Check if request has any fields to update
	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValAllFieldsEmpty)
	}

	// Check if staff user exists
	_, err := s.queries.GetStaffUserByID(ctx, staffUserID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.StaffNotFound, "failed to get staff user", err)
	}

	// Update staff user record using repository
	updatedStaffUser, err := s.repo.Staff.UpdateStaffUser(ctx, staffUserID, sqlxRepo.UpdateStaffUserParams{
		Email: req.Email,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update staff own information", err)
	}

	response := &adminStaffModel.UpdateMeResponse{
		ID:        utils.FormatID(updatedStaffUser.ID),
		Username:  updatedStaffUser.Username,
		Email:     updatedStaffUser.Email,
		Role:      updatedStaffUser.Role,
		IsActive:  utils.PgBoolToBool(updatedStaffUser.IsActive),
		CreatedAt: utils.PgTimestamptzToTimeString(updatedStaffUser.CreatedAt),
		UpdatedAt: utils.PgTimestamptzToTimeString(updatedStaffUser.UpdatedAt),
	}

	return response, nil
}
