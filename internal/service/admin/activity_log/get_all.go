package adminActivityLog

import (
	"context"

	adminActivityLogModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/activity_log"
	"github.com/tkoleo84119/nail-salon-backend/internal/service/cache"
)

type GetAll struct {
	activityLog cache.ActivityLogCacheInterface
}

func NewGetAll(activityLog cache.ActivityLogCacheInterface) GetAllInterface {
	return &GetAll{
		activityLog: activityLog,
	}
}

func (s *GetAll) GetAll(ctx context.Context, req adminActivityLogModel.GetAllParsedRequest) (*adminActivityLogModel.GetAllResponse, error) {
	result, err := s.activityLog.GetRecentActivities(ctx, req.Limit)
	if err != nil {
		return nil, err
	}

	return &adminActivityLogModel.GetAllResponse{
		Total: result.Total,
		Items: result.Activities,
	}, nil
}
