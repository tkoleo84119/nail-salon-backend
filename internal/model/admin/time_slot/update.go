package adminTimeSlot

type UpdateRequest struct {
	StartTime   *string `json:"startTime,omitempty" binding:"omitempty"`
	EndTime     *string `json:"endTime,omitempty" binding:"omitempty"`
	IsAvailable *bool   `json:"isAvailable,omitempty" binding:"omitempty"`
}

type UpdateResponse struct {
	ID          string `json:"id"`
	ScheduleID  string `json:"scheduleId"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	IsAvailable bool   `json:"isAvailable"`
}

func (r *UpdateRequest) HasUpdate() bool {
	return r.StartTime != nil || r.EndTime != nil || r.IsAvailable != nil
}
