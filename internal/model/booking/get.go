package booking

type GetResponse struct {
	ID            string           `json:"id"`
	StoreId       string           `json:"storeId"`
	StoreName     string           `json:"storeName"`
	StylistId     string           `json:"stylistId"`
	StylistName   string           `json:"stylistName"`
	Date          string           `json:"date"`
	TimeSlotId    string           `json:"timeSlotId"`
	StartTime     string           `json:"startTime"`
	EndTime       string           `json:"endTime"`
	MainService   GetServiceItem   `json:"mainService"`
	SubServices   []GetServiceItem `json:"subServices"`
	IsChatEnabled bool             `json:"isChatEnabled"`
	Note          string           `json:"note"`
	Status        string           `json:"status"`
	CreatedAt     string           `json:"createdAt"`
	UpdatedAt     string           `json:"updatedAt"`
}

type GetServiceItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
