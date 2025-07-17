package schedule

// DeleteSchedulesBulkRequest represents the request to delete multiple schedules
type DeleteSchedulesBulkRequest struct {
	StylistID   string   `json:"stylistId" binding:"required"`
	StoreID     string   `json:"storeId" binding:"required"`
	ScheduleIDs []string `json:"scheduleIds" binding:"required,min=1,dive,required"`
}

// DeleteSchedulesBulkResponse represents the response after deleting multiple schedules
type DeleteSchedulesBulkResponse struct {
	Deleted []string `json:"deleted"`
}