package adminStaff

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetMyStaffService struct {
	queries *dbgen.Queries
}

func NewGetMyStaffService(queries *dbgen.Queries) *GetMyStaffService {
	return &GetMyStaffService{
		queries: queries,
	}
}

func (s *GetMyStaffService) GetMyStaff(ctx context.Context, staffUserID int64) (*adminStaffModel.GetMyStaffResponse, error) {
	// Get staff user information
	staffUser, err := s.queries.GetStaffUserByID(ctx, staffUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "Failed to get staff user", err)
	}

	// Convert to response format
	response := &adminStaffModel.GetMyStaffResponse{
		ID:       utils.FormatID(staffUser.ID),
		Username: staffUser.Username,
		Email:    staffUser.Email,
		Role:     staffUser.Role,
		IsActive: staffUser.IsActive.Bool,
	}

	return response, nil
}
