package booking

type GetAllRequest struct {
	Limit  *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
	Sort   *string `form:"sort" binding:"omitempty"`
	Status *string `form:"status" binding:"omitempty"`
}

type GetAllParsedRequest struct {
	Limit  int
	Offset int
	Sort   []string
	Status *[]string
}

type GetAllResponse struct {
	Total int          `json:"total"`
	Items []GetAllItem `json:"items"`
}

type GetAllItem struct {
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
}
