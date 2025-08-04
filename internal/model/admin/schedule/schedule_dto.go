package adminSchedule

import "time"

// TimeSlotRequest represents a time slot in the request
type TimeSlotRequest struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

// ScheduleRequest represents a single schedule in the request
type ScheduleRequest struct {
	WorkDate  string            `json:"workDate" binding:"required"`
	Note      *string           `json:"note,omitempty" binding:"omitempty,max=100"`
	TimeSlots []TimeSlotRequest `json:"timeSlots" binding:"required,min=1,max=30"`
}

// CreateSchedulesBulkRequest represents the request to create multiple schedules
type CreateSchedulesBulkRequest struct {
	StylistID string            `json:"stylistId" binding:"required"`
	StoreID   string            `json:"storeId" binding:"required"`
	Schedules []ScheduleRequest `json:"schedules" binding:"required,min=1,max=50"`
}

// TimeSlotResponse represents a time slot in the response
type TimeSlotResponse struct {
	ID        string `json:"id"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// ScheduleResponse represents a single schedule in the response
type ScheduleResponse struct {
	ScheduleID string             `json:"scheduleId"`
	StylistID  string             `json:"stylistId"`
	StoreID    string             `json:"storeId"`
	WorkDate   string             `json:"workDate"`
	Note       *string            `json:"note,omitempty"`
	TimeSlots  []TimeSlotResponse `json:"timeSlots"`
}

// CreateSchedulesBulkResponse represents the response after creating multiple schedules
type CreateSchedulesBulkResponse []ScheduleResponse

// -------------------------------------------------------------------------------------

// DeleteSchedulesBulkRequest represents the request to delete multiple schedules
type DeleteSchedulesBulkRequest struct {
	StylistID   string   `json:"stylistId" binding:"required"`
	StoreID     string   `json:"storeId" binding:"required"`
	ScheduleIDs []string `json:"scheduleIds" binding:"required,min=1,max=50"`
}

// DeleteSchedulesBulkResponse represents the response after deleting multiple schedules
type DeleteSchedulesBulkResponse struct {
	Deleted []string `json:"deleted"`
}

// -------------------------------------------------------------------------------------

// GetScheduleListRequest represents the request to get schedules list
type GetScheduleListRequest struct {
	StylistID   *string `form:"stylistId"`
	StartDate   string  `form:"startDate" binding:"required"`
	EndDate     string  `form:"endDate" binding:"required"`
	IsAvailable *bool   `form:"isAvailable" binding:"omitempty,boolean"`
}

type GetScheduleListParsedRequest struct {
	StylistID   *[]int64
	StartDate   time.Time
	EndDate     time.Time
	IsAvailable *bool
}

// GetScheduleListResponse represents the response with schedules list
type GetScheduleListResponse struct {
	StylistList []GetScheduleListStylistItem `json:"stylistList"`
}

// GetScheduleListItem represents a single schedule in the list
type GetScheduleListStylistItem struct {
	ID        string                        `json:"id"`
	Name      string                        `json:"name"`
	Schedules []GetScheduleListScheduleItem `json:"schedules"`
}

type GetScheduleListScheduleItem struct {
	ID        string                        `json:"id"`
	WorkDate  string                        `json:"workDate"`
	Note      string                        `json:"note"`
	TimeSlots []GetScheduleListTimeSlotInfo `json:"timeSlots"`
}

// GetScheduleListTimeSlotInfo represents time slot info in schedule list
type GetScheduleListTimeSlotInfo struct {
	ID          string `json:"id"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	IsAvailable bool   `json:"isAvailable"`
}

// -------------------------------------------------------------------------------------

// GetScheduleResponse represents the response for getting a single schedule
type GetScheduleResponse struct {
	ID        string                    `json:"id"`
	WorkDate  string                    `json:"workDate"`
	Note      string                    `json:"note"`
	TimeSlots []GetScheduleTimeSlotInfo `json:"timeSlots"`
}

// GetScheduleTimeSlotInfo represents time slot info in single schedule response
type GetScheduleTimeSlotInfo struct {
	ID          string `json:"id"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	IsAvailable bool   `json:"isAvailable"`
}
