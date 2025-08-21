package stylist

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	stylistModel "github.com/tkoleo84119/nail-salon-backend/internal/model/stylist"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	queries dbgen.Querier
	repo    *sqlxRepo.Repositories
}

func NewGetAll(queries dbgen.Querier, repo *sqlxRepo.Repositories) GetAllInterface {
	return &GetAll{
		queries: queries,
		repo:    repo,
	}
}

func (s *GetAll) GetAll(ctx context.Context, storeID int64, queryParams stylistModel.GetAllParsedRequest) (*stylistModel.GetAllResponse, error) {
	// Validate store exists and is active
	exists, err := s.queries.CheckStoreExistAndActive(ctx, storeID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check store exist and active", err)
	}
	if !exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
	}

	activeCondition := true
	// Get stylists from repository
	total, stylists, err := s.repo.Stylist.GetAllStoreStylistsByFilter(ctx, storeID, sqlxRepo.GetAllStoreStylistsByFilterParams{
		IsActive: &activeCondition,
		Limit:    &queryParams.Limit,
		Offset:   &queryParams.Offset,
		Sort:     &queryParams.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get store stylists", err)
	}

	items := make([]stylistModel.GetAllStylistItem, len(stylists))
	for i, stylist := range stylists {
		items[i] = stylistModel.GetAllStylistItem{
			ID:           utils.FormatID(stylist.ID),
			Name:         utils.PgTextToString(stylist.Name),
			GoodAtShapes: stylist.GoodAtShapes,
			GoodAtColors: stylist.GoodAtColors,
			GoodAtStyles: stylist.GoodAtStyles,
			IsIntrovert:  utils.PgBoolToBool(stylist.IsIntrovert),
		}
	}

	return &stylistModel.GetAllResponse{
		Total: total,
		Items: items,
	}, nil
}
