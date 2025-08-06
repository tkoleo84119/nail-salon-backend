package adminSchedule

import (
	"context"

	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
)

type CreateBulkServiceInterface interface {
	CreateBulk(ctx context.Context, storeID int64, req adminScheduleModel.CreateBulkRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminScheduleModel.CreateBulkResponse, error)
}

type DeleteBulkInterface interface {
	DeleteBulk(ctx context.Context, storeID int64, req adminScheduleModel.DeleteBulkParsedRequest, updaterID int64, updaterRole string, updaterStoreIDs []int64) (*adminScheduleModel.DeleteBulkResponse, error)
}

type GetScheduleListServiceInterface interface {
	GetScheduleList(ctx context.Context, storeID int64, req adminScheduleModel.GetScheduleListParsedRequest, role string, storeIDs []int64) (*adminScheduleModel.GetScheduleListResponse, error)
}

type GetScheduleServiceInterface interface {
	GetSchedule(ctx context.Context, storeID int64, scheduleID int64, role string, storeIDs []int64) (*adminScheduleModel.GetScheduleResponse, error)
}
