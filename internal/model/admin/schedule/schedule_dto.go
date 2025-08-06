package adminSchedule

import "time"

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
