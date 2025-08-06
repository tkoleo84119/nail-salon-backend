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
