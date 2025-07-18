package staff

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
)

type CreateStaffServiceInterface interface {
	CreateStaff(ctx context.Context, req staff.CreateStaffRequest, creatorRole string, creatorStoreIDs []int64) (*staff.CreateStaffResponse, error)
}

type UpdateStaffServiceInterface interface {
	UpdateStaff(ctx context.Context, targetID string, req staff.UpdateStaffRequest, updaterID int64, updaterRole string) (*staff.UpdateStaffResponse, error)
}

type UpdateMyStaffServiceInterface interface {
	UpdateMyStaff(ctx context.Context, req staff.UpdateMyStaffRequest, staffUserID int64) (*staff.UpdateMyStaffResponse, error)
}

type CreateStoreAccessServiceInterface interface {
	CreateStoreAccess(ctx context.Context, targetID string, req staff.CreateStoreAccessRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*staff.CreateStoreAccessResponse, bool, error)
}

type DeleteStoreAccessServiceInterface interface {
	DeleteStoreAccess(ctx context.Context, targetID string, req staff.DeleteStoreAccessRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*staff.DeleteStoreAccessResponse, error)
}
