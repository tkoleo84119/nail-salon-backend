package adminSchedule

type CreateTimeSlotRequest struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

type CreateTimeSlotResponse struct {
	ID          string `json:"id"`
	ScheduleID  string `json:"scheduleId"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	IsAvailable bool   `json:"isAvailable"`
}

// -------------------------------------------------------------------------------------

type UpdateTimeSlotRequest struct {
	StartTime   *string `json:"startTime,omitempty" binding:"omitempty"`
	EndTime     *string `json:"endTime,omitempty" binding:"omitempty"`
	IsAvailable *bool   `json:"isAvailable,omitempty" binding:"omitempty"`
}

type UpdateTimeSlotResponse struct {
	ID          string `json:"id"`
	ScheduleID  string `json:"scheduleId"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	IsAvailable bool   `json:"isAvailable"`
}

func (r *UpdateTimeSlotRequest) HasUpdate() bool {
	return r.StartTime != nil || r.EndTime != nil || r.IsAvailable != nil
}

// -------------------------------------------------------------------------------------

type DeleteTimeSlotResponse struct {
	Deleted string `json:"deleted"`
}
