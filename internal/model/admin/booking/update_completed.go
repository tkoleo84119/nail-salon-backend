package adminBooking

type UpdateCompletedRequest struct {
	ActualDuration *int32 `json:"actualDuration" binding:"omitempty,min=0,max=1440"`
}

type UpdateCompletedResponse struct {
	ID string `json:"id"`
}

func (r UpdateCompletedRequest) HasUpdates() bool {
	return r.ActualDuration != nil
}
