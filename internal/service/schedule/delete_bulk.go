package schedule

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type DeleteSchedulesBulkService struct {
	queries dbgen.Querier
	db      *pgxpool.Pool
}

func NewDeleteSchedulesBulkService(queries dbgen.Querier, db *pgxpool.Pool) *DeleteSchedulesBulkService {
	return &DeleteSchedulesBulkService{
		queries: queries,
		db:      db,
	}
}

func (s *DeleteSchedulesBulkService) DeleteSchedulesBulk(ctx context.Context, req schedule.DeleteSchedulesBulkRequest, staffContext common.StaffContext) (*schedule.DeleteSchedulesBulkResponse, error) {
	// Parse IDs
	stylistID, err := utils.ParseID(req.StylistID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid stylist ID", err)
	}

	storeID, err := utils.ParseID(req.StoreID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid store ID", err)
	}

	scheduleIDs, err := utils.ParseIDSlice(req.ScheduleIDs)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid schedule IDs", err)
	}

	// Check if stylist exists
	stylist, err := s.queries.GetStylistByID(ctx, stylistID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.StylistNotFound, "stylist not found", err)
	}

	// Check permission: STYLIST can only delete their own schedules
	if staffContext.Role == staff.RoleStylist {
		staffUserID, err := utils.ParseID(staffContext.UserID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.AuthStaffFailed, "invalid staff user ID", err)
		}
		if stylist.StaffUserID != staffUserID {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Check if store exists and is active
	store, err := s.queries.GetStoreByID(ctx, storeID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.UserStoreNotFound, "store not found", err)
	}
	if !store.IsActive.Bool {
		return nil, errorCodes.NewServiceError(errorCodes.UserStoreNotActive, "store is not active", err)
	}

	// Check if staff has access to this store
	hasAccess := false
	for _, storeAccess := range staffContext.StoreList {
		if storeAccess.ID == req.StoreID {
			hasAccess = true
			break
		}
	}
	if !hasAccess {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Get schedules with time slots
	schedules, err := s.queries.GetSchedulesWithTimeSlotsByIDs(ctx, scheduleIDs)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get schedules with time slots", err)
	}

	// Check if all schedules exist and belong to the stylist/store
	if len(schedules) != len(scheduleIDs) {
		return nil, errorCodes.NewServiceError(errorCodes.ScheduleNotFound, "some schedules not found", nil)
	}
	for _, schedule := range schedules {
		if schedule.StoreID != storeID {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleNotBelongToStore)
		}

		if schedule.StylistID != stylistID {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleNotBelongToStylist)
		}

		// Check if time slots are not available => mean this schedule is already booked
		if !schedule.IsAvailable.Bool {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleAlreadyBookedDoNotDelete)
		}
	}

	// Delete schedules and time slots in transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := dbgen.New(tx)

	// Delete time slots first (foreign key constraint)
	if err := qtx.DeleteTimeSlotsByScheduleIDs(ctx, scheduleIDs); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete time slots", err)
	}

	// Delete schedules
	if err := qtx.DeleteSchedulesByIDs(ctx, scheduleIDs); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete schedules", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	// Build response with deleted schedule IDs
	response := &schedule.DeleteSchedulesBulkResponse{
		Deleted: req.ScheduleIDs,
	}

	return response, nil
}
