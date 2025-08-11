package schedule

// GetTimeSlotsByScheduleResponse represents the response for getting time slots by schedule
type GetTimeSlotsByScheduleResponse struct {
	Items []TimeSlotResponseItem `json:"items"`
}

// TimeSlotResponseItem represents a single time slot item
type TimeSlotResponseItem struct {
	ID              string `json:"id"`
	StartTime       string `json:"startTime"`
	EndTime         string `json:"endTime"`
	DurationMinutes int    `json:"durationMinutes"`
}
