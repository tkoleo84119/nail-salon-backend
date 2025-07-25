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

// BookingServiceInfo represents service information for booking
type BookingServiceInfo struct {
	ServiceId     int64
	ServiceName   string
	Price         float64
	IsMainService bool
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

// -------------------------------------------------------------------------------------

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

// -------------------------------------------------------------------------------------

// GetMyBookingsQueryParams represents query parameters for getting customer bookings
type GetMyBookingsQueryParams struct {
	Limit  int      `form:"limit,default=20" binding:"omitempty,min=1,max=100"`
	Offset int      `form:"offset,default=0" binding:"omitempty,min=0"`
	Status []string `form:"status" binding:"omitempty"`
}

// GetMyBookingsResponse represents the response for getting customer bookings
type GetMyBookingsResponse struct {
	Total int                      `json:"total"`
	Items []GetMyBookingsItemModel `json:"items"`
}

// GetMyBookingsItemModel represents a single booking item in the list
type GetMyBookingsItemModel struct {
	ID          string                     `json:"id"`
	StoreId     string                     `json:"storeId"`
	StoreName   string                     `json:"storeName"`
	StylistId   string                     `json:"stylistId"`
	StylistName string                     `json:"stylistName"`
	Date        string                     `json:"date"`
	TimeSlot    GetMyBookingsTimeSlotModel `json:"timeSlot"`
	Status      string                     `json:"status"`
}

// GetMyBookingsTimeSlotModel represents time slot information in booking list
type GetMyBookingsTimeSlotModel struct {
	ID        string `json:"id"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// -------------------------------------------------------------------------------------

// GetMyBookingResponse represents the response for getting a single booking detail
type GetMyBookingResponse struct {
	ID          string                      `json:"id"`
	StoreId     string                      `json:"storeId"`
	StoreName   string                      `json:"storeName"`
	StylistId   string                      `json:"stylistId"`
	StylistName string                      `json:"stylistName"`
	Date        string                      `json:"date"`
	TimeSlot    GetMyBookingTimeSlotModel   `json:"timeSlot"`
	Services    []GetMyBookingServiceModel  `json:"services"`
	Note        *string                     `json:"note,omitempty"`
	Status      string                      `json:"status"`
}

// GetMyBookingTimeSlotModel represents time slot information in single booking detail
type GetMyBookingTimeSlotModel struct {
	ID        string `json:"id"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// GetMyBookingServiceModel represents service information in single booking detail
type GetMyBookingServiceModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
