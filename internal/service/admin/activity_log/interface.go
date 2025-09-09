package adminActivityLog

import (
	"context"

	adminActivityLogModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/activity_log"
)

type GetAllInterface interface {
	GetAll(ctx context.Context, req adminActivityLogModel.GetAllParsedRequest) (*adminActivityLogModel.GetAllResponse, error)
}