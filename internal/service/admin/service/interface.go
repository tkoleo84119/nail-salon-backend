package adminService

import (
	"context"

	adminServiceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/service"
)

type CreateInterface interface {
	Create(ctx context.Context, req adminServiceModel.CreateRequest, creatorRole string) (*adminServiceModel.CreateResponse, error)
}

type GetAllInterface interface {
	GetAll(ctx context.Context, req adminServiceModel.GetAllParsedRequest) (*adminServiceModel.GetAllResponse, error)
}

type GetInterface interface {
	Get(ctx context.Context, serviceID int64) (*adminServiceModel.GetResponse, error)
}

type UpdateInterface interface {
	Update(ctx context.Context, serviceID int64, req adminServiceModel.UpdateRequest, updaterRole string) (*adminServiceModel.UpdateResponse, error)
}
