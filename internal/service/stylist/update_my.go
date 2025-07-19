package stylist

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/stylist"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
)

type UpdateMyStylistService struct {
	queries           dbgen.Querier
	stylistRepository sqlx.StylistRepositoryInterface
}

func NewUpdateMyStylistService(queries dbgen.Querier, stylistRepository sqlx.StylistRepositoryInterface) *UpdateMyStylistService {
	return &UpdateMyStylistService{
		queries:           queries,
		stylistRepository: stylistRepository,
	}
}

func (s *UpdateMyStylistService) UpdateMyStylist(ctx context.Context, req stylist.UpdateMyStylistRequest, staffUserID int64) (*stylist.UpdateMyStylistResponse, error) {
	// ensure at least one field is provided for update
	if !req.HasUpdate() {
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
