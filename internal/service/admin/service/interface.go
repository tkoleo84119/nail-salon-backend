package adminService

import (
	"context"

	adminServiceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/service"
)

type CreateServiceInterface interface {
	CreateService(ctx context.Context, req adminServiceModel.CreateServiceRequest, creatorRole string) (*adminServiceModel.CreateServiceResponse, error)
}

type UpdateServiceInterface interface {
	UpdateService(ctx context.Context, serviceID int64, req adminServiceModel.UpdateServiceRequest, updaterRole string) (*adminServiceModel.UpdateServiceResponse, error)
}

type GetServiceListServiceInterface interface {
	GetServiceList(ctx context.Context, req adminServiceModel.GetServiceListParsedRequest) (*adminServiceModel.GetServiceListResponse, error)
}

type GetServiceServiceInterface interface {
	GetService(ctx context.Context, serviceID int64) (*adminServiceModel.GetServiceResponse, error)
}
