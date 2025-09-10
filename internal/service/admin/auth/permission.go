package adminAuth

import (
	"context"

	adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Permission struct{}

func NewPermission() *Permission {
	return &Permission{}
}

func (s *Permission) Permission(ctx context.Context, staffContext *common.StaffContext) (*adminAuthModel.PermissionResponse, error) {
	// Convert StoreList to StoreAccess format
	storeAccess := make([]adminAuthModel.StoreInfo, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		storeAccess[i] = adminAuthModel.StoreInfo{
			ID:   utils.FormatID(store.ID),
			Name: store.Name,
		}
	}

	response := &adminAuthModel.PermissionResponse{
		ID:          utils.FormatID(staffContext.UserID),
		Username:    staffContext.Username,
		Role:        staffContext.Role,
		StoreAccess: storeAccess,
	}

	return response, nil
}
