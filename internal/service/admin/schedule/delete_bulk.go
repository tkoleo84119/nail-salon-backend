package adminSchedule

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type DeleteBulk struct {
	queries dbgen.Querier
}

func NewDeleteBulk(queries dbgen.Querier) *DeleteBulk {
	return &DeleteBulk{
		queries: queries,
	}
}

func (s *DeleteBulk) DeleteBulk(ctx context.Context, storeID int64, req adminScheduleModel.DeleteBulkParsedRequest, updaterID int64, updaterRole string, updaterStoreIDs []int64) (*adminScheduleModel.DeleteBulkResponse, error) {
	// Check if stylist exists
	stylist, err := s.queries.GetStylistByID(ctx, req.StylistID)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotFound)
	}

	// Check permission: STYLIST can only delete their own schedules
	if updaterRole == common.RoleStylist {
		if stylist.StaffUserID != updaterID {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Check if staff has access to this store
	hasAccess, err := utils.CheckStoreAccess(storeID, updaterStoreIDs)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to check store access", err)
	}
	if !hasAccess {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Get schedules with time slots
	schedules, err := s.queries.GetSchedulesWithTimeSlotsByIDs(ctx, req.ScheduleIDs)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get schedules with time slots", err)
	}

	// Check if all schedules exist and belong to the stylist/store
	if uniqueSchedulesCounts(schedules) != len(req.ScheduleIDs) {
		return nil, errorCodes.NewServiceError(errorCodes.ScheduleNotFound, "some schedules not found", nil)
	}
	for _, schedule := range schedules {
		if schedule.StoreID != storeID {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleNotBelongToStore)
		}

		if schedule.StylistID != req.StylistID {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleNotBelongToStylist)
		}

		// Check if time slots are not available => mean this schedule is already booked
		if !schedule.IsAvailable.Bool {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleAlreadyBookedDoNotDelete)
		}
	}

	// Delete schedules
	if err := s.queries.DeleteSchedulesByIDs(ctx, req.ScheduleIDs); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete schedules", err)
	}

	// Build response with deleted schedule IDs
	response := &adminScheduleModel.DeleteBulkResponse{
		Deleted: make([]string, len(req.ScheduleIDs)),
	}

	for i, id := range req.ScheduleIDs {
		response.Deleted[i] = utils.FormatID(id)
	}

	return response, nil
}

// schedules include many time slots, so we need to count unique schedules
func uniqueSchedulesCounts(schedules []dbgen.GetSchedulesWithTimeSlotsByIDsRow) int {
	uniqueMap := make(map[int64]int)

	for _, schedule := range schedules {
		uniqueMap[schedule.ID]++
	}

	return len(uniqueMap)
}
