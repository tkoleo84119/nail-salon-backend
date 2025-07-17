package schedule

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
