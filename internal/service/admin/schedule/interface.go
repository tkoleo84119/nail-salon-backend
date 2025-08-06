package adminSchedule

import (
	"context"

	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
)

type CreateBulkInterface interface {
	CreateBulk(ctx context.Context, storeID int64, req adminScheduleModel.CreateBulkRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminScheduleModel.CreateBulkResponse, error)
}

type DeleteBulkInterface interface {
	DeleteBulk(ctx context.Context, storeID int64, req adminScheduleModel.DeleteBulkParsedRequest, updaterID int64, updaterRole string, updaterStoreIDs []int64) (*adminScheduleModel.DeleteBulkResponse, error)
}

type GetAllInterface interface {
	GetAll(ctx context.Context, storeID int64, req adminScheduleModel.GetAllParsedRequest, role string, storeIDs []int64) (*adminScheduleModel.GetAllResponse, error)
}

type GetInterface interface {
	Get(ctx context.Context, storeID int64, scheduleID int64, role string, storeIDs []int64) (*adminScheduleModel.GetResponse, error)
}
