package booking

// CreateMyBookingRequest represents the request for creating a new booking
type CreateMyBookingRequest struct {
	StoreId       string   `json:"storeId" binding:"required"`
	StylistId     string   `json:"stylistId" binding:"required"`
	TimeSlotId    string   `json:"timeSlotId" binding:"required"`
	MainServiceId string   `json:"mainServiceId" binding:"required"`
	SubServiceIds []string `json:"subServiceIds,omitempty" binding:"omitempty,max=5"`
	IsChatEnabled *bool    `json:"isChatEnabled,omitempty"`
	Note          *string  `json:"note,omitempty" binding:"omitempty,max=500"`
}

// CreateMyBookingResponse represents the response for creating a new booking
type CreateMyBookingResponse struct {
	ID              string   `json:"id"`
	StoreId         string   `json:"storeId"`
	StoreName       string   `json:"storeName"`
	StylistId       string   `json:"stylistId"`
	StylistName     string   `json:"stylistName"`
	Date            string   `json:"date"`
	TimeSlotId      string   `json:"timeSlotId"`
	StartTime       string   `json:"startTime"`
	EndTime         string   `json:"endTime"`
	MainServiceName string   `json:"mainServiceName"`
	SubServiceNames []string `json:"subServiceNames"`
	IsChatEnabled   bool     `json:"isChatEnabled"`
	Note            *string  `json:"note,omitempty"`
	Status          string   `json:"status"`
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
}

// -------------------------------------------------------------------------------------

// UpdateMyBookingRequest represents the request for updating my booking
type UpdateMyBookingRequest struct {
	StoreId       *string  `json:"storeId,omitempty"`
	StylistId     *string  `json:"stylistId,omitempty"`
	TimeSlotId    *string  `json:"timeSlotId,omitempty"`
	MainServiceId *string  `json:"mainServiceId,omitempty"`
	SubServiceIds []string `json:"subServiceIds,omitempty" binding:"omitempty,max=5"`
	IsChatEnabled *bool    `json:"isChatEnabled,omitempty"`
	Note          *string  `json:"note,omitempty" binding:"omitempty,max=500"`
}

// UpdateMyBookingResponse represents the response for updating my booking
type UpdateMyBookingResponse struct {
	ID              string   `json:"id"`
	StoreId         string   `json:"storeId"`
	StoreName       string   `json:"storeName"`
	StylistId       string   `json:"stylistId"`
	StylistName     string   `json:"stylistName"`
	Date            string   `json:"date"`
	TimeSlotId      string   `json:"timeSlotId"`
	StartTime       string   `json:"startTime"`
	EndTime         string   `json:"endTime"`
	MainServiceName string   `json:"mainServiceName"`
	SubServiceNames []string `json:"subServiceNames"`
	IsChatEnabled   bool     `json:"isChatEnabled"`
	Note            *string  `json:"note,omitempty"`
	Status          string   `json:"status"`
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
}

// BookingServiceInfo represents service information for booking
type BookingServiceInfo struct {
	ServiceId     int64
	ServiceName   string
	Price         float64
	IsMainService bool
}

// HasUpdates checks if the request has at least one field to update
func (r UpdateMyBookingRequest) HasUpdates() bool {
	return r.StoreId != nil || r.StylistId != nil || r.TimeSlotId != nil ||
		r.MainServiceId != nil || r.SubServiceIds != nil || r.IsChatEnabled != nil || r.Note != nil
}

// HasTimeSlotUpdate checks if time slot related fields are being updated
func (r UpdateMyBookingRequest) HasTimeSlotUpdate() bool {
	return r.StoreId != nil || r.StylistId != nil || r.TimeSlotId != nil || r.MainServiceId != nil || r.SubServiceIds != nil
}

// IsTimeSlotUpdateComplete checks if all time slot related fields are provided together
func (r UpdateMyBookingRequest) IsTimeSlotUpdateComplete() bool {
	if !r.HasTimeSlotUpdate() {
		return true
	}
	return r.StoreId != nil && r.StylistId != nil && r.TimeSlotId != nil && r.MainServiceId != nil && r.SubServiceIds != nil
}
