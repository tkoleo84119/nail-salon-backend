package adminSchedule

import "time"

type UpdateRequest struct {
	StylistID string  `json:"stylistId" binding:"required"`
	WorkDate  *string `json:"workDate"`
	Note      *string `json:"note" binding:"omitempty,max=100"`
}

type UpdateParsedRequest struct {
	StylistID int64
	WorkDate  *time.Time
	Note      *string
}

type UpdateResponse struct {
	ID        string               `json:"id"`
}

func (r *UpdateRequest) HasUpdates() bool {
	return r.WorkDate != nil || r.Note != nil
}
