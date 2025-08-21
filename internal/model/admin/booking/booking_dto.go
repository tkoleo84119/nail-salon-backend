package adminBooking

// CancelBookingRequest represents the request to cancel a booking by staff
type CancelBookingRequest struct {
	Status       string  `json:"status" binding:"required,oneof=CANCELLED NO_SHOW"`
	CancelReason *string `json:"cancelReason,omitempty" binding:"omitempty,max=100"`
}

// CancelBookingResponse represents the response after cancelling a booking by staff
type CancelBookingResponse struct {
	ID string `json:"id"`
}
