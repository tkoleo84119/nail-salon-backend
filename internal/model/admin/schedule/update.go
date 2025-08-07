package adminSchedule

type UpdateRequest struct {
	StylistID string  `json:"stylistId" binding:"required"`
	WorkDate  *string `json:"workDate,omitempty"`
	Note      *string `json:"note,omitempty" binding:"omitempty,max=100"`
}

type UpdateParsedRequest struct {
	StylistID int64
	WorkDate  *string
	Note      *string
}

type UpdateResponse struct {
	ID        string               `json:"id"`
	WorkDate  string               `json:"workDate"`
	Note      string               `json:"note"`
	TimeSlots []UpdateTimeSlotInfo `json:"timeSlots"`
}

type UpdateTimeSlotInfo struct {
	ID          string `json:"id"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	IsAvailable bool   `json:"isAvailable"`
}

func (r *UpdateParsedRequest) HasUpdates() bool {
	return r.WorkDate != nil || r.Note != nil
}
