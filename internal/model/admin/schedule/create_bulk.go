package adminSchedule

// TimeSlotRequest represents a time slot in the request
type CreateBulkTimeSlotRequest struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

// CreateBulkScheduleRequest represents a single schedule in the request
type CreateBulkScheduleRequest struct {
	WorkDate  string                    `json:"workDate" binding:"required"`
	Note      *string                   `json:"note,omitempty" binding:"omitempty,max=100"`
	TimeSlots []CreateBulkTimeSlotRequest `json:"timeSlots" binding:"required,min=1,max=20"`
}

// CreateBulkRequest represents the request to create multiple schedules
type CreateBulkRequest struct {
	StylistID string                    `json:"stylistId" binding:"required"`
	Schedules []CreateBulkScheduleRequest `json:"schedules" binding:"required,min=1,max=31"`
}

// TimeSlotResponse represents a time slot in the response
type CreateBulkTimeSlotResponse struct {
	ID        string `json:"id"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// CreateBulkScheduleResponse represents a single schedule in the response
type CreateBulkScheduleResponse struct {
	ID        string                    `json:"id"`
	WorkDate  string                    `json:"workDate"`
	Note      string                    `json:"note"`
	TimeSlots []CreateBulkTimeSlotResponse `json:"timeSlots"`
}

// CreateBulkResponse represents the response after creating multiple schedules
type CreateBulkResponse struct {
	Schedules []CreateBulkScheduleResponse `json:"schedules"`
}
