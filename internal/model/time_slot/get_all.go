package timeSlot

type GetAllResponse struct {
	TimeSlots []GetAllResponseItem `json:"timeSlots"`
}

type GetAllResponseItem struct {
	ID              string `json:"id"`
	StartTime       string `json:"startTime"`
	EndTime         string `json:"endTime"`
	DurationMinutes int    `json:"durationMinutes"`
}
