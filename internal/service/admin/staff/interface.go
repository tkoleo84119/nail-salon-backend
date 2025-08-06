package adminStaff

import (
	"context"

	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
)

type CreateInterface interface {
	Create(ctx context.Context, req adminStaffModel.CreateParsedRequest, creatorRole string, creatorStoreIDs []int64) (*adminStaffModel.CreateResponse, error)
}

type GetAllInterface interface {
	GetAll(ctx context.Context, req adminStaffModel.GetAllParsedRequest) (*adminStaffModel.GetAllResponse, error)
}

type GetInterface interface {
	Get(ctx context.Context, staffID int64) (*adminStaffModel.GetResponse, error)
}

type GetMeInterface interface {
	GetMe(ctx context.Context, staffUserID int64) (*adminStaffModel.GetMeResponse, error)
}

type UpdateInterface interface {
	Update(ctx context.Context, targetID int64, req adminStaffModel.UpdateRequest, updaterID int64, updaterRole string) (*adminStaffModel.UpdateResponse, error)
}

type UpdateMeInterface interface {
	UpdateMe(ctx context.Context, req adminStaffModel.UpdateMeRequest, staffUserID int64) (*adminStaffModel.UpdateMeResponse, error)
}

type CreateStoreAccessServiceInterface interface {
	CreateStoreAccess(ctx context.Context, staffID int64, storeID int64, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminStaffModel.CreateStoreAccessResponse, bool, error)
}

type DeleteStoreAccessBulkServiceInterface interface {
	DeleteStoreAccessBulk(ctx context.Context, targetID int64, storeIDs []int64, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminStaffModel.DeleteStoreAccessBulkResponse, error)
}

type GetStaffStoreAccessServiceInterface interface {
	GetStaffStoreAccess(ctx context.Context, staffID int64) (*adminStaffModel.GetStaffStoreAccessResponse, error)
}
