package adminBooking

type CreateRequest struct {
	CustomerID    string    `json:"customerId" binding:"required"`
	TimeSlotID    string    `json:"timeSlotId" binding:"required"`
	MainServiceID string    `json:"mainServiceId" binding:"required"`
	SubServiceIDs *[]string `json:"subServiceIds" binding:"omitempty,max=10"`
	IsChatEnabled *bool     `json:"isChatEnabled" binding:"omitempty"`
	StoreNote     *string   `json:"storeNote" binding:"omitempty,max=255"`
}

type CreateParsedRequest struct {
	CustomerID    int64
	TimeSlotID    int64
	MainServiceID int64
	SubServiceIDs []int64
	IsChatEnabled bool
	StoreNote     *string
}

type CreateResponse struct {
	ID string `json:"id"`
}
