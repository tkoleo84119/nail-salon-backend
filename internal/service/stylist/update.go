package stylist

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/stylist"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
)

type UpdateStylistService struct {
	queries           dbgen.Querier
	stylistRepository sqlx.StylistRepositoryInterface
}

func NewUpdateStylistService(queries dbgen.Querier, stylistRepository sqlx.StylistRepositoryInterface) *UpdateStylistService {
	return &UpdateStylistService{
		queries:           queries,
		stylistRepository: stylistRepository,
	}
}

func (s *UpdateStylistService) UpdateStylist(ctx context.Context, req stylist.UpdateStylistRequest, staffUserID int64) (*stylist.UpdateStylistResponse, error) {
	// Get staff user info to check role
	staffUser, err := s.queries.GetStaffUserByID(ctx, staffUserID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.AuthStaffFailed, "failed to get staff user", err)
	}

	// Check if user is SUPER_ADMIN (not allowed to update stylist)
	if staffUser.Role == staff.RoleSuperAdmin {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Check if request has any fields to update
	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValAllFieldsEmpty)
	}

	// Check if stylist exists for this staff user
	exists, err := s.queries.CheckStylistExistsByStaffUserID(ctx, pgtype.Int8{Int64: staffUserID, Valid: true})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check stylist existence", err)
	}
	if !exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotCreated)
	}

	// Update stylist record using repository
	response, err := s.stylistRepository.UpdateStylist(ctx, staffUserID, req)
	if err != nil {
		if strings.Contains(err.Error(), "no rows returned") {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update stylist", err)
	}

	return response, nil
}