package service

import (
	"context"

	serviceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/service"
)

type GetAllInterface interface {
	GetAll(ctx context.Context, queryParams serviceModel.GetAllParsedRequest) (*serviceModel.GetAllResponse, error)
}
