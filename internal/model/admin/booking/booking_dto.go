package adminBooking

// CreateBookingRequest represents the request to create a booking for admin
type CreateBookingRequest struct {
	CustomerID    string   `json:"customerId" binding:"required"`
	TimeSlotID    string   `json:"timeSlotId" binding:"required"`
	MainServiceID string   `json:"mainServiceId" binding:"required"`
	SubServiceIDs []string `json:"subServiceIds" binding:"omitempty,max=10"`
	IsChatEnabled bool     `json:"isChatEnabled"`
	Note          *string  `json:"note,omitempty" binding:"omitempty,max=200"`
}

// CreateBookingResponse represents the response after creating a booking for admin
type CreateBookingResponse struct {
	ID              string   `json:"id"`
	StoreID         string   `json:"storeId"`
	StoreName       string   `json:"storeName"`
	StylistID       string   `json:"stylistId"`
	StylistName     string   `json:"stylistName"`
	Date            string   `json:"date"`
	TimeSlotID      string   `json:"timeSlotId"`
	StartTime       string   `json:"startTime"`
	EndTime         string   `json:"endTime"`
	MainServiceName string   `json:"mainServiceName"`
	SubServiceNames []string `json:"subServiceNames"`
	IsChatEnabled   bool     `json:"isChatEnabled"`
	Note            string   `json:"note"`
	Status          string   `json:"status"`
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
}

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
