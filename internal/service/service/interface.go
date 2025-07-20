package service

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/service"
)

type CreateServiceInterface interface {
	CreateService(ctx context.Context, req service.CreateServiceRequest, creatorRole string) (*service.CreateServiceResponse, error)
}

type UpdateServiceInterface interface {
	UpdateService(ctx context.Context, serviceID string, req service.UpdateServiceRequest, updaterRole string) (*service.UpdateServiceResponse, error)
}