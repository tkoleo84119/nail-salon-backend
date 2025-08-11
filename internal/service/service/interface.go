package service

import (
	"context"

	serviceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/service"
)

type GetAllInterface interface {
	GetAll(ctx context.Context, storeID int64, queryParams serviceModel.GetAllParsedRequest) (*serviceModel.GetAllResponse, error)
}
