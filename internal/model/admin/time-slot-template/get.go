package adminTimeSlotTemplate

type GetResponse struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Note      string        `json:"note"`
	Updater   string        `json:"updater"`
	CreatedAt string        `json:"createdAt"`
	UpdatedAt string        `json:"updatedAt"`
	Items     []GetItemInfo `json:"items"`
}

// GetItemInfo represents a time slot template item in the response
type GetItemInfo struct {
	ID        string `json:"id"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
