package booking

type CancelRequest struct {
	CancelReason *string `json:"cancelReason,omitempty" binding:"omitempty,max=255"`
}

type CancelResponse struct {
	ID          string `json:"id"`
	StoreId     string `json:"storeId"`
	StoreName   string `json:"storeName"`
	StylistId   string `json:"stylistId"`
	StylistName string `json:"stylistName"`
	Date        string `json:"date"`
	TimeSlotId  string `json:"timeSlotId"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
