package adminSchedule

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateBulk struct {
	queries *dbgen.Queries
	db      *pgxpool.Pool
}

func NewCreateBulk(queries *dbgen.Queries, db *pgxpool.Pool) *CreateBulk {
	return &CreateBulk{
		queries: queries,
		db:      db,
	}
}

func (s *CreateBulk) CreateBulk(ctx context.Context, storeID int64, req adminScheduleModel.CreateBulkRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminScheduleModel.CreateBulkResponse, error) {
	parsedStylistID, err := utils.ParseID(req.StylistID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid stylist ID", err)
	}

	// Check if stylist exists
	stylist, err := s.queries.GetStylistByID(ctx, parsedStylistID)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotFound)
	}

	// Check permission: STYLIST can only create schedules for themselves
	if creatorRole == common.RoleStylist {
		if stylist.StaffUserID != creatorID {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Check if staff has access to this store
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to check store access", err)
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

		exists, err := s.queries.CheckScheduleDateExists(ctx, dbgen.CheckScheduleDateExistsParams{
			StoreID:   storeID,
			StylistID: parsedStylistID,
			WorkDate:  utils.TimeToPgDate(workDate),
		})
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
	var response adminScheduleModel.CreateBulkResponse

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := dbgen.New(tx)

	// Batch insert schedules
	if _, err := qtx.BatchCreateSchedules(ctx, scheduleRows); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to batch create schedules", err)
	}

	// Batch insert time slots
	if _, err := qtx.BatchCreateTimeSlots(ctx, timeSlotRows); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to batch create time slots", err)
	}

	if err := tx.Commit(ctx); err != nil {
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
func (s *CreateBulk) prepareBatchData(schedules []adminScheduleModel.CreateBulkScheduleRequest, storeID, stylistID int64) ([]dbgen.BatchCreateSchedulesParams, []dbgen.BatchCreateTimeSlotsParams, error) {
	var scheduleRows []dbgen.BatchCreateSchedulesParams
	var timeSlotRows []dbgen.BatchCreateTimeSlotsParams

	for _, scheduleReq := range schedules {
		// Parse work date
		workDate, err := utils.DateStringToPgDate(scheduleReq.WorkDate)
		if err != nil {
			return nil, nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid work date format", err)
		}

		// Generate schedule ID
		scheduleID := utils.GenerateID()
		now := utils.TimeToPgTimestamptz(time.Now())

		scheduleRow := dbgen.BatchCreateSchedulesParams{
			ID:        scheduleID,
			StoreID:   storeID,
			StylistID: stylistID,
			WorkDate:  workDate,
			Note:      utils.StringPtrToPgText(scheduleReq.Note, true),
			CreatedAt: now,
			UpdatedAt: now,
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
			isAvailable := true

			timeSlotRow := dbgen.BatchCreateTimeSlotsParams{
				ID:          timeSlotID,
				ScheduleID:  scheduleID,
				StartTime:   utils.TimeToPgTime(startTime),
				EndTime:     utils.TimeToPgTime(endTime),
				IsAvailable: utils.BoolPtrToPgBool(&isAvailable),
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			timeSlotRows = append(timeSlotRows, timeSlotRow)
		}
	}

	return scheduleRows, timeSlotRows, nil
}

// buildResponseFromScheduleRows builds the response from joined schedule and time slot data
func (s *CreateBulk) buildResponseFromScheduleRows(ctx context.Context, scheduleRows []dbgen.BatchCreateSchedulesParams, timeSlotRows []dbgen.BatchCreateTimeSlotsParams) (adminScheduleModel.CreateBulkResponse, error) {
	scheduleGroups := make(map[int64]adminScheduleModel.CreateBulkScheduleResponse, len(scheduleRows))

	for _, scheduleRow := range scheduleRows {
		scheduleGroups[scheduleRow.ID] = adminScheduleModel.CreateBulkScheduleResponse{
			ID:        utils.FormatID(scheduleRow.ID),
			WorkDate:  utils.PgDateToDateString(scheduleRow.WorkDate),
			Note:      utils.PgTextToString(scheduleRow.Note),
			TimeSlots: []adminScheduleModel.CreateBulkTimeSlotResponse{},
		}
	}

	for _, timeSlotRow := range timeSlotRows {
		timeSlot := adminScheduleModel.CreateBulkTimeSlotResponse{
			ID:        utils.FormatID(timeSlotRow.ID),
			StartTime: utils.PgTimeToTimeString(timeSlotRow.StartTime),
			EndTime:   utils.PgTimeToTimeString(timeSlotRow.EndTime),
		}

		if schedule, ok := scheduleGroups[timeSlotRow.ScheduleID]; ok {
			schedule.TimeSlots = append(schedule.TimeSlots, timeSlot)
			scheduleGroups[timeSlotRow.ScheduleID] = schedule
		}
	}

	var response adminScheduleModel.CreateBulkResponse
	for _, schedule := range scheduleGroups {
		response.Schedules = append(response.Schedules, schedule)
	}

	return response, nil
}

func (s *CreateBulk) validateSchedules(schedules []adminScheduleModel.CreateBulkScheduleRequest) error {
	workDates := make(map[string]bool)

	for _, scheduleReq := range schedules {
		// can't pass date before today
		if err := IsValidWorkDate(scheduleReq.WorkDate); err != nil {
			return err
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

func (s *CreateBulk) validateTimeSlots(timeSlots []adminScheduleModel.CreateBulkTimeSlotRequest) error {
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

func IsValidWorkDate(workDateStr string) error {
	workDate, err := utils.DateStringToTime(workDateStr)
	if err != nil {
		return errorCodes.NewServiceError(errorCodes.ValFieldDateFormat, "invalid date format", err)
	}

	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to load location", err)
	}

	workDateInLoc := time.Date(
		workDate.Year(),
		workDate.Month(),
		workDate.Day(),
		0, 0, 0, 0,
		loc,
	)

	today := time.Now().In(loc).Truncate(24 * time.Hour)

	if workDateInLoc.Before(today) {
		return errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleCannotCreateBeforeToday)
	}

	return nil
}
