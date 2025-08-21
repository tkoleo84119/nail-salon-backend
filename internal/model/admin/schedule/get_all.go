package adminSchedule

import "time"

type GetAllRequest struct {
	StylistID   *string `form:"stylistId"`
	StartDate   string  `form:"startDate" binding:"required"`
	EndDate     string  `form:"endDate" binding:"required"`
	IsAvailable *bool   `form:"isAvailable" binding:"omitempty"`
}

type GetAllParsedRequest struct {
	StylistID   *[]int64
	StartDate   time.Time
	EndDate     time.Time
	IsAvailable *bool
}

type GetAllResponse struct {
	StylistList []GetAllStylistItem `json:"stylistList"`
}

type GetAllStylistItem struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	Schedules []GetAllScheduleItem `json:"schedules"`
}

type GetAllScheduleItem struct {
	ID        string               `json:"id"`
	WorkDate  string               `json:"workDate"`
	Note      string               `json:"note"`
	TimeSlots []GetAllTimeSlotInfo `json:"timeSlots"`
}

type GetAllTimeSlotInfo struct {
	ID          string `json:"id"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	IsAvailable bool   `json:"isAvailable"`
}
