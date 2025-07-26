package schedule

// GetStoreSchedulesRequest represents the request for getting store schedules
type GetStoreSchedulesRequest struct {
	StartDate string `form:"startDate" binding:"required"`
	EndDate   string `form:"endDate" binding:"required"`
}

// GetStoreSchedulesResponse represents the response for getting store schedules
type GetStoreSchedulesResponse struct {
	Total int                         `json:"total"`
	Items []StoreScheduleResponseItem `json:"items"`
}

// StoreScheduleResponseItem represents a single schedule item
type StoreScheduleResponseItem struct {
	Date           string `json:"date"`
	AvailableSlots int    `json:"available_slots"`
}
