package adminStaff

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetStaffStoreAccessService struct {
	repo *sqlxRepo.Repositories
}

func NewGetStaffStoreAccessService(repo *sqlxRepo.Repositories) *GetStaffStoreAccessService {
	return &GetStaffStoreAccessService{
		repo: repo,
	}
}

func (s *GetStaffStoreAccessService) GetStaffStoreAccess(ctx context.Context, staffID int64) (*adminStaffModel.GetStaffStoreAccessResponse, error) {
	// Verify staff user exists
	_, err := s.repo.Staff.GetStaffUserByID(ctx, staffID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get staff user", err)
	}

	// Get staff store access
	storeAccessList, err := s.repo.StaffUserStoreAccess.GetStaffUserStoreAccessByStaffId(ctx, staffID, nil)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get staff store access", err)
	}

	// Convert to response format
	var items []common.Store
	for _, access := range storeAccessList {
		items = append(items, common.Store{
			ID:   utils.FormatID(access.StoreID),
			Name: access.Name,
		})
	}

	return &adminStaffModel.GetStaffStoreAccessResponse{
		StoreList: items,
	}, nil
}
