package adminBooking

// CreateBookingRequest represents the request to create a booking for admin
type CreateBookingRequest struct {
	CustomerID     string   `json:"customerId" binding:"required"`
	TimeSlotID     string   `json:"timeSlotId" binding:"required"`
	MainServiceID  string   `json:"mainServiceId" binding:"required"`
	SubServiceIDs  []string `json:"subServiceIds" binding:"omitempty,max=10"`
	IsChatEnabled  bool     `json:"isChatEnabled"`
	Note           *string  `json:"note,omitempty" binding:"omitempty,max=200"`
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