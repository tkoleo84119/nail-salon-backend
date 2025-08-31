package booking

type CancelRequest struct {
	HasChatPermission *bool   `json:"hasChatPermission" binding:"omitempty"`
	CancelReason      *string `json:"cancelReason,omitempty" binding:"omitempty,max=255"`
}

type CancelResponse struct {
	ID              string   `json:"id"`
	StoreId         string   `json:"storeId"`
	StoreName       string   `json:"storeName"`
	StylistId       string   `json:"stylistId"`
	StylistName     string   `json:"stylistName"`
	CustomerName    string   `json:"customerName"`
	CustomerPhone   string   `json:"customerPhone"`
	Date            string   `json:"date"`
	TimeSlotId      string   `json:"timeSlotId"`
	StartTime       string   `json:"startTime"`
	EndTime         string   `json:"endTime"`
	MainServiceName string   `json:"mainServiceName"`
	SubServiceNames []string `json:"subServiceNames"`
	Status          string   `json:"status"`
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
}
