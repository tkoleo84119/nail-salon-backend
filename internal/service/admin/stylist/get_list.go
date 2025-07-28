package adminStylist

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	adminStylistModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stylist"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// GetStylistListServiceInterface defines the interface for getting stylist list
type GetStylistListServiceInterface interface {
	GetStylistList(ctx context.Context, storeID string, req adminStylistModel.GetStylistListRequest, staffContext common.StaffContext) (*adminStylistModel.GetStylistListResponse, error)
}

type GetStylistListService struct {
	queries     *dbgen.Queries
	stylistRepo sqlxRepo.StylistRepositoryInterface
}

func NewGetStylistListService(queries *dbgen.Queries, stylistRepo sqlxRepo.StylistRepositoryInterface) *GetStylistListService {
	return &GetStylistListService{
		queries:     queries,
		stylistRepo: stylistRepo,
	}
}

func (s *GetStylistListService) GetStylistList(ctx context.Context, storeID string, req adminStylistModel.GetStylistListRequest, staffContext common.StaffContext) (*adminStylistModel.GetStylistListResponse, error) {
	// Parse store ID
	storeIDInt, err := utils.ParseID(storeID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Invalid store ID", err)
	}

	// Verify store exists
	_, err = s.queries.GetStoreByID(ctx, storeIDInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get store", err)
	}

	// Check store access for the staff member (except SUPER_ADMIN)
	if staffContext.Role != adminStaffModel.RoleSuperAdmin {
		hasAccess, err := utils.CheckOneStoreAccess(storeIDInt, staffContext)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "Failed to check store access", err)
		}
		if !hasAccess {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Get stylists from repository with dynamic filtering
	stylists, total, err := s.stylistRepo.GetStoreStylistList(ctx, storeIDInt, sqlxRepo.GetStoreStylistListParams{
		Name:        req.Name,
		IsIntrovert: req.IsIntrovert,
		Limit:       req.Limit,
		Offset:      req.Offset,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get stylist list", err)
	}

	// Convert to response models
	response := &adminStylistModel.GetStylistListResponse{
		Total: total,
		Items: make([]adminStylistModel.GetStylistListItem, len(stylists)),
	}
	for i, stylist := range stylists {
		response.Items[i] = adminStylistModel.GetStylistListItem{
			ID:           utils.FormatID(stylist.ID),
			StaffUserID:  utils.FormatID(stylist.StaffUserID),
			Name:         stylist.Name,
			GoodAtShapes: stylist.GoodAtShapes,
			GoodAtColors: stylist.GoodAtColors,
			GoodAtStyles: stylist.GoodAtStyles,
			IsIntrovert:  stylist.IsIntrovert,
		}
	}

	return response, nil
}
