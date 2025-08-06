package adminSchedule

import (
	"context"

	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

type CreateSchedulesBulkServiceInterface interface {
	CreateSchedulesBulk(ctx context.Context, storeID int64, req adminScheduleModel.CreateSchedulesBulkRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminScheduleModel.CreateSchedulesBulkResponse, error)
}

type DeleteBulkInterface interface {
	DeleteBulk(ctx context.Context, storeID int64, req adminScheduleModel.DeleteBulkParsedRequest, updaterID int64, updaterRole string, updaterStoreIDs []int64) (*adminScheduleModel.DeleteBulkResponse, error)
}

type CreateTimeSlotServiceInterface interface {
	CreateTimeSlot(ctx context.Context, scheduleID string, req adminScheduleModel.CreateTimeSlotRequest, staffContext common.StaffContext) (*adminScheduleModel.CreateTimeSlotResponse, error)
}

type UpdateTimeSlotServiceInterface interface {
	UpdateTimeSlot(ctx context.Context, scheduleID string, timeSlotID string, req adminScheduleModel.UpdateTimeSlotRequest, staffContext common.StaffContext) (*adminScheduleModel.UpdateTimeSlotResponse, error)
}

type DeleteTimeSlotServiceInterface interface {
	DeleteTimeSlot(ctx context.Context, scheduleID string, timeSlotID string, staffContext common.StaffContext) (*adminScheduleModel.DeleteTimeSlotResponse, error)
}

type GetScheduleListServiceInterface interface {
	GetScheduleList(ctx context.Context, storeID int64, req adminScheduleModel.GetScheduleListParsedRequest, role string, storeIDs []int64) (*adminScheduleModel.GetScheduleListResponse, error)
}

type GetScheduleServiceInterface interface {
	GetSchedule(ctx context.Context, storeID int64, scheduleID int64, role string, storeIDs []int64) (*adminScheduleModel.GetScheduleResponse, error)
}
