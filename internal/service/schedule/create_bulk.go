package schedule

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateSchedulesBulkService struct {
	queries dbgen.Querier
	db      *pgxpool.Pool
}

func NewCreateSchedulesBulkService(queries dbgen.Querier, db *pgxpool.Pool) *CreateSchedulesBulkService {
	return &CreateSchedulesBulkService{
		queries: queries,
		db:      db,
	}
}

func (s *CreateSchedulesBulkService) CreateSchedulesBulk(ctx context.Context, req schedule.CreateSchedulesBulkRequest, staffContext common.StaffContext) (*schedule.CreateSchedulesBulkResponse, error) {
	// Parse IDs
	stylistID, err := utils.ParseID(req.StylistID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid stylist ID", err)
	}

	storeID, err := utils.ParseID(req.StoreID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid store ID", err)
	}

	// Check if stylist exists
	stylist, err := s.queries.GetStylistByID(ctx, stylistID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.StylistNotFound, "stylist not found", err)
	}

	// Check permission: STYLIST can only create schedules for themselves
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

	// Validate time slots and check for conflicts
	if err := s.validateSchedules(req.Schedules); err != nil {
		return nil, err
	}

	// Check for existing schedules
	for _, scheduleReq := range req.Schedules {
		workDate, err := utils.DateStringToPgDate(scheduleReq.WorkDate)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid work date format", err)
		}

		exists, err := s.queries.CheckScheduleExists(ctx, dbgen.CheckScheduleExistsParams{
			StoreID:   storeID,
			StylistID: stylistID,
			WorkDate:  workDate,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check schedule existence", err)
		}
		if exists {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleAlreadyExists)
		}
	}

	// Prepare batch data for schedules and time slots
	scheduleRows, timeSlotRows, createdScheduleIDs, err := s.prepareBatchData(req.Schedules, storeID, stylistID)
	if err != nil {
		return nil, err
	}

	// Create schedules in transaction using batch insert
	var response schedule.CreateSchedulesBulkResponse

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

	// Get the created schedules with time slots using efficient join query
	createdScheduleWithTimeSlots, err := s.queries.GetSchedulesWithTimeSlotsByIDs(ctx, createdScheduleIDs)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get created schedules", err)
	}

	// Build response
	response, err = s.buildResponseFromScheduleRows(ctx, createdScheduleWithTimeSlots, storeID, stylistID)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// prepareBatchData prepares data for batch insertion and returns schedule rows, time slot rows, a mapping, and created schedule IDs
func (s *CreateSchedulesBulkService) prepareBatchData(schedules []schedule.ScheduleRequest, storeID, stylistID int64) ([]dbgen.BatchCreateSchedulesParams, []dbgen.BatchCreateTimeSlotsParams, []int64, error) {
	now := time.Now()
	var scheduleRows []dbgen.BatchCreateSchedulesParams
	var timeSlotRows []dbgen.BatchCreateTimeSlotsParams
	var createdScheduleIDs []int64

	for _, scheduleReq := range schedules {
		// Parse work date
		workDate, err := utils.DateStringToPgDate(scheduleReq.WorkDate)
		if err != nil {
			return nil, nil, nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid work date format", err)
		}

		// Generate schedule ID
		scheduleID := utils.GenerateID()
		createdScheduleIDs = append(createdScheduleIDs, scheduleID)

		// Prepare schedule row
		noteValue := utils.StringPtrToPgText(scheduleReq.Note, false)

		scheduleRow := dbgen.BatchCreateSchedulesParams{
			ID:        scheduleID,
			StoreID:   storeID,
			StylistID: stylistID,
			WorkDate:  workDate,
			Note:      noteValue,
			CreatedAt: utils.TimeToPgTimestamptz(now),
			UpdatedAt: utils.TimeToPgTimestamptz(now),
		}
		scheduleRows = append(scheduleRows, scheduleRow)

		// Prepare time slot rows
		for _, timeSlotReq := range scheduleReq.TimeSlots {
			startTime, err := utils.TimeStringToTime(timeSlotReq.StartTime)
			if err != nil {
				return nil, nil, nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid start time format", err)
			}

			endTime, err := utils.TimeStringToTime(timeSlotReq.EndTime)
			if err != nil {
				return nil, nil, nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid end time format", err)
			}

			timeSlotID := utils.GenerateID()
			isAvailable := true

			timeSlotRow := dbgen.BatchCreateTimeSlotsParams{
				ID:          timeSlotID,
				ScheduleID:  scheduleID,
				StartTime:   utils.TimeToPgTime(startTime),
				EndTime:     utils.TimeToPgTime(endTime),
				IsAvailable: utils.BoolPtrToPgBool(&isAvailable),
				CreatedAt:   utils.TimeToPgTimestamptz(now),
				UpdatedAt:   utils.TimeToPgTimestamptz(now),
			}
			timeSlotRows = append(timeSlotRows, timeSlotRow)
		}
	}

	return scheduleRows, timeSlotRows, createdScheduleIDs, nil
}

// buildResponseFromScheduleRows builds the response from joined schedule and time slot data
func (s *CreateSchedulesBulkService) buildResponseFromScheduleRows(ctx context.Context, scheduleRows []dbgen.GetSchedulesWithTimeSlotsByIDsRow, storeID, stylistID int64) (schedule.CreateSchedulesBulkResponse, error) {
	var response schedule.CreateSchedulesBulkResponse

	// Group time slots by schedule ID
	scheduleGroups := make(map[int64][]dbgen.GetSchedulesWithTimeSlotsByIDsRow)
	scheduleData := make(map[int64]dbgen.GetSchedulesWithTimeSlotsByIDsRow)

	for _, row := range scheduleRows {
		scheduleGroups[row.ID] = append(scheduleGroups[row.ID], row)
		scheduleData[row.ID] = row
	}

	// Build response for each schedule
	for scheduleID, rows := range scheduleGroups {
		scheduleIDStr := utils.FormatID(scheduleID)
		scheduleInfo := scheduleData[scheduleID]

		// Convert time slots to response format
		var timeSlotResponses []schedule.TimeSlotResponse
		for _, row := range rows {
			// Skip rows without time slot data (can happen with LEFT JOIN if no time slots exist)
			if !row.TimeSlotID.Valid {
				continue
			}

			timeSlotResponses = append(timeSlotResponses, schedule.TimeSlotResponse{
				ID:        utils.PgInt8ToIDString(row.TimeSlotID),
				StartTime: utils.PgTimeToTimeString(row.StartTime),
				EndTime:   utils.PgTimeToTimeString(row.EndTime),
			})
		}

		// Build schedule response
		var notePtr *string
		if scheduleInfo.Note.Valid {
			notePtr = &scheduleInfo.Note.String
		}

		scheduleResponse := schedule.ScheduleResponse{
			ScheduleID: scheduleIDStr,
			StylistID:  utils.FormatID(stylistID),
			StoreID:    utils.FormatID(storeID),
			WorkDate:   scheduleInfo.WorkDate.Time.Format("2006-01-02"),
			Note:       notePtr,
			TimeSlots:  timeSlotResponses,
		}

		response = append(response, scheduleResponse)
	}

	return response, nil
}

func (s *CreateSchedulesBulkService) validateSchedules(schedules []schedule.ScheduleRequest) error {
	workDates := make(map[string]bool)

	for _, scheduleReq := range schedules {
		// Check for duplicate work dates
		if workDates[scheduleReq.WorkDate] {
			return errorCodes.NewServiceErrorWithCode(errorCodes.ValDuplicateWorkDate)
		}
		workDates[scheduleReq.WorkDate] = true

		// Validate time slots for each schedule
		if err := s.validateTimeSlots(scheduleReq.TimeSlots); err != nil {
			return err
		}
	}

	return nil
}

func (s *CreateSchedulesBulkService) validateTimeSlots(timeSlots []schedule.TimeSlotRequest) error {
	if len(timeSlots) == 0 {
		return errorCodes.NewServiceErrorWithCode(errorCodes.ValTimeSlotRequired)
	}

	// Parse and sort time slots
	type timeSlotParsed struct {
		StartTime time.Time
		EndTime   time.Time
	}

	var parsedSlots []timeSlotParsed
	for _, slot := range timeSlots {
		startTime, err := utils.TimeStringToTime(slot.StartTime)
		if err != nil {
			return errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, fmt.Sprintf("invalid start time format: %s", slot.StartTime), err)
		}

		endTime, err := utils.TimeStringToTime(slot.EndTime)
		if err != nil {
			return errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, fmt.Sprintf("invalid end time format: %s", slot.EndTime), err)
		}

		if !endTime.After(startTime) {
			return errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotInvalidTimeRange)
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
