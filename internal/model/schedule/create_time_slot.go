package schedule

type CreateTimeSlotRequest struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

type CreateTimeSlotResponse struct {
	ID         string `json:"id"`
	ScheduleID string `json:"scheduleId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}
