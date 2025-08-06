package adminSchedule

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
