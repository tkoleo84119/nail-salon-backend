package adminSchedule

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateSchedulesBulkService struct {
	db   *sqlx.DB
	repo *sqlxRepo.Repositories
}

func NewCreateSchedulesBulkService(db *sqlx.DB, repo *sqlxRepo.Repositories) *CreateSchedulesBulkService {
	return &CreateSchedulesBulkService{
		db:   db,
		repo: repo,
	}
}

func (s *CreateSchedulesBulkService) CreateSchedulesBulk(ctx context.Context, storeID int64, req adminScheduleModel.CreateSchedulesBulkRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminScheduleModel.CreateSchedulesBulkResponse, error) {
	parsedStylistID, err := utils.ParseID(req.StylistID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid stylist ID", err)
	}

	// Check if stylist exists
	stylist, err := s.repo.Stylist.GetStylistByID(ctx, parsedStylistID)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotFound)
	}

	// Check permission: STYLIST can only create schedules for themselves
	if creatorRole == common.RoleStylist {
		if stylist.ID != creatorID {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Check if store exists and is active
	store, err := s.repo.Store.GetStoreByID(ctx, storeID, nil)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
	}
	if !store.IsActive.Bool {
		return nil, errorCodes.NewServiceError(errorCodes.StoreNotActive, "store is not active", err)
	}

	// Check if staff has access to this store
	hasAccess, err := utils.CheckStoreAccess(storeID, creatorStoreIDs)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to check store access", err)
	}
	if !hasAccess {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Validate time slots and check for conflicts
	if err := s.validateSchedules(req.Schedules); err != nil {
		return nil, err
	}

	// Check for existing schedules
	for _, scheduleReq := range req.Schedules {
		workDate, err := utils.DateStringToTime(scheduleReq.WorkDate)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValFieldDateFormat, "invalid work date format", err)
		}

		exists, err := s.repo.Schedule.CheckScheduleExists(ctx, storeID, parsedStylistID, workDate)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check schedule existence", err)
		}
		if exists {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleAlreadyExists)
		}
	}

	// Prepare batch data for schedules and time slots
	scheduleRows, timeSlotRows, err := s.prepareBatchData(req.Schedules, storeID, stylist.ID)
	if err != nil {
		return nil, err
	}

	// Create schedules in transaction using batch insert
	var response adminScheduleModel.CreateSchedulesBulkResponse

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback()

	// Batch insert schedules
	if err := s.repo.Schedule.BatchCreateSchedulesTx(ctx, tx, scheduleRows); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to batch create schedules", err)
	}

	// Batch insert time slots
	if err := s.repo.TimeSlot.BatchCreateTimeSlotsTx(ctx, tx, timeSlotRows); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to batch create time slots", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	// Build response
	response, err = s.buildResponseFromScheduleRows(ctx, scheduleRows, timeSlotRows)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// prepareBatchData prepares data for batch insertion and returns schedule rows, time slot rows, a mapping, and created schedule IDs
func (s *CreateSchedulesBulkService) prepareBatchData(schedules []adminScheduleModel.ScheduleRequest, storeID, stylistID int64) ([]sqlxRepo.BatchCreateSchedulesTxParams, []sqlxRepo.BatchCreateTimeSlotsTxParams, error) {
	var scheduleRows []sqlxRepo.BatchCreateSchedulesTxParams
	var timeSlotRows []sqlxRepo.BatchCreateTimeSlotsTxParams

	for _, scheduleReq := range schedules {
		// Parse work date
		workDate, err := utils.DateStringToPgDate(scheduleReq.WorkDate)
		if err != nil {
			return nil, nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid work date format", err)
		}

		// Generate schedule ID
		scheduleID := utils.GenerateID()

		// Prepare schedule row
		noteValue := utils.StringPtrToPgText(scheduleReq.Note, false)

		scheduleRow := sqlxRepo.BatchCreateSchedulesTxParams{
			ID:        scheduleID,
			StoreID:   storeID,
			StylistID: stylistID,
			WorkDate:  workDate,
			Note:      noteValue,
		}
		scheduleRows = append(scheduleRows, scheduleRow)

		// Prepare time slot rows
		for _, timeSlotReq := range scheduleReq.TimeSlots {
			startTime, err := utils.TimeStringToTime(timeSlotReq.StartTime)
			if err != nil {
				return nil, nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid start time format", err)
			}

			endTime, err := utils.TimeStringToTime(timeSlotReq.EndTime)
			if err != nil {
				return nil, nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid end time format", err)
			}

			timeSlotID := utils.GenerateID()

			timeSlotRow := sqlxRepo.BatchCreateTimeSlotsTxParams{
				ID:         timeSlotID,
				ScheduleID: scheduleID,
				StartTime:  utils.TimeToPgTime(startTime),
				EndTime:    utils.TimeToPgTime(endTime),
			}
			timeSlotRows = append(timeSlotRows, timeSlotRow)
		}
	}

	return scheduleRows, timeSlotRows, nil
}

// buildResponseFromScheduleRows builds the response from joined schedule and time slot data
func (s *CreateSchedulesBulkService) buildResponseFromScheduleRows(ctx context.Context, scheduleRows []sqlxRepo.BatchCreateSchedulesTxParams, timeSlotRows []sqlxRepo.BatchCreateTimeSlotsTxParams) (adminScheduleModel.CreateSchedulesBulkResponse, error) {
	scheduleGroups := make(map[int64]adminScheduleModel.ScheduleResponse, len(scheduleRows))

	for _, scheduleRow := range scheduleRows {
		scheduleGroups[scheduleRow.ID] = adminScheduleModel.ScheduleResponse{
			ID:        utils.FormatID(scheduleRow.ID),
			WorkDate:  utils.PgDateToDateString(scheduleRow.WorkDate),
			Note:      utils.PgTextToString(scheduleRow.Note),
			TimeSlots: []adminScheduleModel.TimeSlotResponse{},
		}
	}

	for _, timeSlotRow := range timeSlotRows {
		timeSlot := adminScheduleModel.TimeSlotResponse{
			ID:        utils.FormatID(timeSlotRow.ID),
			StartTime: utils.PgTimeToTimeString(timeSlotRow.StartTime),
			EndTime:   utils.PgTimeToTimeString(timeSlotRow.EndTime),
		}

		if schedule, ok := scheduleGroups[timeSlotRow.ScheduleID]; ok {
			schedule.TimeSlots = append(schedule.TimeSlots, timeSlot)
			scheduleGroups[timeSlotRow.ScheduleID] = schedule
		}
	}

	var response adminScheduleModel.CreateSchedulesBulkResponse
	for _, schedule := range scheduleGroups {
		response.Schedules = append(response.Schedules, schedule)
	}

	return response, nil
}

func (s *CreateSchedulesBulkService) validateSchedules(schedules []adminScheduleModel.ScheduleRequest) error {
	workDates := make(map[string]bool)

	for _, scheduleReq := range schedules {
		// can't pass date before today
		if scheduleReq.WorkDate < time.Now().Format("2006-01-02") {
			return errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleCannotCreateBeforeToday)
		}

		// Check for duplicate work dates
		if workDates[scheduleReq.WorkDate] {
			return errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleDuplicateWorkDateInput)
		}
		workDates[scheduleReq.WorkDate] = true

		// Validate time slots for each schedule
		if err := s.validateTimeSlots(scheduleReq.TimeSlots); err != nil {
			return err
		}
	}

	return nil
}

func (s *CreateSchedulesBulkService) validateTimeSlots(timeSlots []adminScheduleModel.TimeSlotRequest) error {
	// Parse and sort time slots
	type timeSlotParsed struct {
		StartTime time.Time
		EndTime   time.Time
	}

	var parsedSlots []timeSlotParsed
	for _, slot := range timeSlots {
		startTime, err := utils.TimeStringToTime(slot.StartTime)
		if err != nil {
			return errorCodes.NewServiceError(errorCodes.ValFieldTimeFormat, fmt.Sprintf("invalid start time format: %s", slot.StartTime), err)
		}

		endTime, err := utils.TimeStringToTime(slot.EndTime)
		if err != nil {
			return errorCodes.NewServiceError(errorCodes.ValFieldTimeFormat, fmt.Sprintf("invalid end time format: %s", slot.EndTime), err)
		}

		if !endTime.After(startTime) {
			return errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotEndBeforeStart)
		}

		parsedSlots = append(parsedSlots, timeSlotParsed{
			StartTime: startTime,
			EndTime:   endTime,
		})
	}

	// Sort by start time
	sort.Slice(parsedSlots, func(i, j int) bool {
		return parsedSlots[i].StartTime.Before(parsedSlots[j].StartTime)
	})

	// Check for overlapping time slots
	for i := 1; i < len(parsedSlots); i++ {
		if parsedSlots[i].StartTime.Before(parsedSlots[i-1].EndTime) {
			return errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotConflict)
		}
	}

	return nil
}
