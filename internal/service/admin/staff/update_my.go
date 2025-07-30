package adminStaff

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
)

type UpdateMyStaffService struct {
	queries dbgen.Querier
	repo    *sqlx.Repositories
}

func NewUpdateMyStaffService(queries dbgen.Querier, repo *sqlx.Repositories) *UpdateMyStaffService {
	return &UpdateMyStaffService{
		queries: queries,
		repo:    repo,
	}
}

func (s *UpdateMyStaffService) UpdateMyStaff(ctx context.Context, req adminStaffModel.UpdateMyStaffRequest, staffUserID int64) (*adminStaffModel.UpdateMyStaffResponse, error) {
	// Check if request has any fields to update
	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValAllFieldsEmpty)
	}

	// Check if staff user exists
	_, err := s.queries.GetStaffUserByID(ctx, staffUserID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.StaffNotFound, "failed to get staff user", err)
	}

	// Check email uniqueness if email is being updated
	if req.Email != nil {
		exists, err := s.queries.CheckEmailUniqueForUpdate(ctx, dbgen.CheckEmailUniqueForUpdateParams{
			Email: *req.Email,
			ID:    staffUserID,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check email uniqueness", err)
		}
		if exists {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffEmailExists)
		}
	}

	// Update staff user record using repository
	response, err := s.repo.Staff.UpdateMyStaff(ctx, staffUserID, req)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update staff own information", err)
	}

	return response, nil
}
