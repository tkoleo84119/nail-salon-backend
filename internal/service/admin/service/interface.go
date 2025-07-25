package adminService

import (
	"context"

	adminServiceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/service"
)

type CreateServiceInterface interface {
	CreateService(ctx context.Context, req adminServiceModel.CreateServiceRequest, creatorRole string) (*adminServiceModel.CreateServiceResponse, error)
}

type UpdateServiceInterface interface {
	UpdateService(ctx context.Context, serviceID string, req adminServiceModel.UpdateServiceRequest, updaterRole string) (*adminServiceModel.UpdateServiceResponse, error)
}