package staff

import (
	"context"
	"strings"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
)

type UpdateStaffMeServiceInterface interface {
	UpdateStaffMe(ctx context.Context, req staff.UpdateStaffMeRequest, staffUserID int64) (*staff.UpdateStaffMeResponse, error)
}

type UpdateStaffMeService struct {
	queries              dbgen.Querier
	staffUserRepository  sqlx.StaffUserRepositoryInterface
}

func NewUpdateStaffMeService(queries dbgen.Querier, staffUserRepository sqlx.StaffUserRepositoryInterface) *UpdateStaffMeService {
	return &UpdateStaffMeService{
		queries:             queries,
		staffUserRepository: staffUserRepository,
	}
}

func (s *UpdateStaffMeService) UpdateStaffMe(ctx context.Context, req staff.UpdateStaffMeRequest, staffUserID int64) (*staff.UpdateStaffMeResponse, error) {
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
	response, err := s.staffUserRepository.UpdateStaffMe(ctx, staffUserID, req)
	if err != nil {
		if strings.Contains(err.Error(), "no rows returned") {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthStaffFailed)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update staff user", err)
	}

	return response, nil
}