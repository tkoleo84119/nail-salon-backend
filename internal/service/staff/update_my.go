package staff

import (
	"context"
	"strings"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
)

type UpdateMyStaffService struct {
	queries             dbgen.Querier
	staffUserRepository sqlx.StaffUserRepositoryInterface
}

func NewUpdateMyStaffService(queries dbgen.Querier, staffUserRepository sqlx.StaffUserRepositoryInterface) *UpdateMyStaffService {
	return &UpdateMyStaffService{
		queries:             queries,
		staffUserRepository: staffUserRepository,
	}
}

func (s *UpdateMyStaffService) UpdateMyStaff(ctx context.Context, req staff.UpdateMyStaffRequest, staffUserID int64) (*staff.UpdateMyStaffResponse, error) {
	// Check if request has any fields to update
	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValAllFieldsEmpty)
	}

	// Check if staff user exists
	_, err := s.queries.GetStaffUserByID(ctx, staffUserID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.AuthStaffFailed, "failed to get staff user", err)
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
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.UserEmailExists)
		}
	}

	// Update staff user record using repository
	response, err := s.staffUserRepository.UpdateMyStaff(ctx, staffUserID, req)
	if err != nil {
		if strings.Contains(err.Error(), "no rows returned") {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthStaffFailed)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update staff user", err)
	}

	return response, nil
}
