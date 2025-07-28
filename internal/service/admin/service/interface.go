package adminService

import (
	"context"

	adminServiceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

type CreateServiceInterface interface {
	CreateService(ctx context.Context, req adminServiceModel.CreateServiceRequest, creatorRole string) (*adminServiceModel.CreateServiceResponse, error)
}

type UpdateServiceInterface interface {
	UpdateService(ctx context.Context, serviceID string, req adminServiceModel.UpdateServiceRequest, updaterRole string) (*adminServiceModel.UpdateServiceResponse, error)
}

type GetServiceListServiceInterface interface {
	GetServiceList(ctx context.Context, storeID string, req adminServiceModel.GetServiceListRequest, staffContext common.StaffContext) (*adminServiceModel.GetServiceListResponse, error)
}

type GetServiceServiceInterface interface {
	GetService(ctx context.Context, storeID, serviceID string, staffContext common.StaffContext) (*adminServiceModel.GetServiceResponse, error)
}