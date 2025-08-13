package booking

type BookingServiceInfo struct {
	ServiceId     int64
	ServiceName   string
	Price         float64
	IsMainService bool
}

// CancelMyBookingRequest represents the request for canceling my booking
type CancelMyBookingRequest struct {
	CancelReason *string `json:"cancelReason,omitempty" binding:"omitempty,max=100"`
}

// CancelMyBookingResponse represents the response for canceling my booking
type CancelMyBookingResponse struct {
	ID           string  `json:"id"`
	Status       string  `json:"status"`
	CancelReason *string `json:"cancelReason,omitempty"`
}
