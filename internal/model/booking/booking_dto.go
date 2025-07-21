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
}

// BookingServiceInfo represents service information for booking
type BookingServiceInfo struct {
	ServiceId     int64
	ServiceName   string
	Price         float64
	IsMainService bool
}
