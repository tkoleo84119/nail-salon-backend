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
	ID          string                     `json:"id"`
	StoreId     string                     `json:"storeId"`
	StoreName   string                     `json:"storeName"`
	StylistId   string                     `json:"stylistId"`
	StylistName string                     `json:"stylistName"`
	Date        string                     `json:"date"`
	TimeSlot    GetMyBookingTimeSlotModel  `json:"timeSlot"`
	Services    []GetMyBookingServiceModel `json:"services"`
	Note        *string                    `json:"note,omitempty"`
	Status      string                     `json:"status"`
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
