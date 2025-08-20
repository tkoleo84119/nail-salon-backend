package adminBooking

// UpdateBookingByStaffRequest represents the request to update a booking by staff
type UpdateBookingByStaffRequest struct {
	TimeSlotID    *string  `json:"timeSlotId,omitempty"`
	MainServiceID *string  `json:"mainServiceId,omitempty"`
	SubServiceIDs []string `json:"subServiceIds,omitempty" binding:"omitempty,max=5"`
	IsChatEnabled *bool    `json:"isChatEnabled,omitempty"`
	Note          *string  `json:"note,omitempty" binding:"omitempty,max=200"`
}

// UpdateBookingByStaffResponse represents the response after updating a booking by staff
type UpdateBookingByStaffResponse struct {
	ID string `json:"id"`
}

func (r UpdateBookingByStaffRequest) HasUpdates() bool {
	return r.TimeSlotID != nil || r.MainServiceID != nil || r.SubServiceIDs != nil || r.IsChatEnabled != nil || r.Note != nil
}

func (r UpdateBookingByStaffRequest) HasTimeSlotUpdate() bool {
	return r.TimeSlotID != nil || r.MainServiceID != nil || r.SubServiceIDs != nil
}

func (r UpdateBookingByStaffRequest) IsTimeSlotUpdateComplete() bool {
	if !r.HasTimeSlotUpdate() {
		return false
	}

	if r.TimeSlotID == nil || r.MainServiceID == nil || r.SubServiceIDs == nil {
		return false
	}

	return true
}

// -------------------------------------------------------------------------------------

// CancelBookingRequest represents the request to cancel a booking by staff
type CancelBookingRequest struct {
	Status       string  `json:"status" binding:"required,oneof=CANCELLED NO_SHOW"`
	CancelReason *string `json:"cancelReason,omitempty" binding:"omitempty,max=100"`
}

// CancelBookingResponse represents the response after cancelling a booking by staff
type CancelBookingResponse struct {
	ID string `json:"id"`
}
