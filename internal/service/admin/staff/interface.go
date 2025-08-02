package adminStaff

import (
	"context"

	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
)

type CreateStaffServiceInterface interface {
	CreateStaff(ctx context.Context, req adminStaffModel.CreateStaffParsedRequest, creatorRole string, creatorStoreIDs []int64) (*adminStaffModel.CreateStaffResponse, error)
}

type UpdateStaffServiceInterface interface {
	UpdateStaff(ctx context.Context, targetID int64, req adminStaffModel.UpdateStaffRequest, updaterID int64, updaterRole string) (*adminStaffModel.UpdateStaffResponse, error)
}

type UpdateMyStaffServiceInterface interface {
	UpdateMyStaff(ctx context.Context, req adminStaffModel.UpdateMyStaffRequest, staffUserID int64) (*adminStaffModel.UpdateMyStaffResponse, error)
}

type CreateStoreAccessServiceInterface interface {
	CreateStoreAccess(ctx context.Context, staffID int64, req adminStaffModel.CreateStoreAccessRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminStaffModel.CreateStoreAccessResponse, bool, error)
}

type DeleteStoreAccessBulkServiceInterface interface {
	DeleteStoreAccessBulk(ctx context.Context, targetID string, req adminStaffModel.DeleteStoreAccessBulkRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminStaffModel.DeleteStoreAccessBulkResponse, error)
}

type GetStaffListServiceInterface interface {
	GetStaffList(ctx context.Context, req adminStaffModel.GetStaffListParsedRequest) (*adminStaffModel.GetStaffListResponse, error)
}

type GetMyStaffServiceInterface interface {
	GetMyStaff(ctx context.Context, staffUserID int64) (*adminStaffModel.GetMyStaffResponse, error)
}

type GetStaffServiceInterface interface {
	GetStaff(ctx context.Context, staffID int64) (*adminStaffModel.GetStaffResponse, error)
}

type GetStaffStoreAccessServiceInterface interface {
	GetStaffStoreAccess(ctx context.Context, staffID int64) (*adminStaffModel.GetStaffStoreAccessResponse, error)
}
