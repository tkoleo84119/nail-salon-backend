package adminBooking

type CancelRequest struct {
	Status       string  `json:"status" binding:"required,oneof=CANCELLED NO_SHOW"`
	CancelReason *string `json:"cancelReason" binding:"omitempty,max=255"`
}

type CancelResponse struct {
	ID string `json:"id"`
}
