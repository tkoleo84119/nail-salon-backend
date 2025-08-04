package adminStylist

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStylistModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stylist"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetStylistListService struct {
	repo *sqlxRepo.Repositories
}

func NewGetStylistListService(repo *sqlxRepo.Repositories) *GetStylistListService {
	return &GetStylistListService{
		repo: repo,
	}
}

func (s *GetStylistListService) GetStylistList(ctx context.Context, storeID int64, req adminStylistModel.GetStylistListParsedRequest, role string, storeIds []int64) (*adminStylistModel.GetStylistListResponse, error) {
	// Verify store exists
	_, err := s.repo.Store.GetStoreByID(ctx, storeID, nil)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get store", err)
	}

	// Check store access for the staff member (except SUPER_ADMIN)
	if role != common.RoleSuperAdmin {
		hasAccess, err := utils.CheckStoreAccess(storeID, storeIds)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "Failed to check store access", err)
		}
		if !hasAccess {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Get stylists from repository with dynamic filtering
	total, stylists, err := s.repo.Stylist.GetStoreAllStylistByFilter(ctx, storeID, sqlxRepo.GetStoreAllStylistByFilterParams{
		Name:        req.Name,
		IsIntrovert: req.IsIntrovert,
		IsActive:    req.IsActive,
		Limit:       &req.Limit,
		Offset:      &req.Offset,
		Sort:        &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get stylist list", err)
	}

	// Convert to response models
	itemDTOs := make([]adminStylistModel.GetStylistListItem, len(stylists))
	for i, stylist := range stylists {
		itemDTOs[i] = adminStylistModel.GetStylistListItem{
			ID:           utils.FormatID(stylist.ID),
			StaffUserID:  utils.FormatID(stylist.StaffUserID),
			Name:         utils.PgTextToString(stylist.Name),
			GoodAtShapes: stylist.GoodAtShapes,
			GoodAtColors: stylist.GoodAtColors,
			GoodAtStyles: stylist.GoodAtStyles,
			IsIntrovert:  utils.PgBoolToBool(stylist.IsIntrovert),
			IsActive:     utils.PgBoolToBool(stylist.IsActive),
		}
	}

	return &adminStylistModel.GetStylistListResponse{
		Total: total,
		Items: itemDTOs,
	}, nil
}
