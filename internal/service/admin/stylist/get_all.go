package adminStylist

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStylistModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stylist"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	repo *sqlxRepo.Repositories
}

func NewGetAll(repo *sqlxRepo.Repositories) GetAllInterface {
	return &GetAll{
		repo: repo,
	}
}

func (s *GetAll) GetAll(ctx context.Context, storeID int64, req adminStylistModel.GetAllParsedRequest, role string, storeIds []int64) (*adminStylistModel.GetAllResponse, error) {
	// Check store access for the staff member
	if err := utils.CheckStoreAccess(storeID, storeIds, role); err != nil {
		return nil, err
	}

	// Get stylists from repository with dynamic filtering
	total, stylists, err := s.repo.Stylist.GetAllStoreStylistsByFilter(ctx, storeID, sqlxRepo.GetAllStoreStylistsByFilterParams{
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
	itemDTOs := make([]adminStylistModel.GetAllItem, len(stylists))
	for i, stylist := range stylists {
		itemDTOs[i] = adminStylistModel.GetAllItem{
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

	return &adminStylistModel.GetAllResponse{
		Total: total,
		Items: itemDTOs,
	}, nil
}
