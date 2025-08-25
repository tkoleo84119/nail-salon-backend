package adminSchedule

type DeleteBulkRequest struct {
	StylistID   string   `json:"stylistId" binding:"required"`
	ScheduleIDs []string `json:"scheduleIds" binding:"required,min=1,max=31"`
}

type DeleteBulkParsedRequest struct {
	StylistID   int64
	ScheduleIDs []int64
}

type DeleteBulkResponse struct {
	Deleted []string `json:"deleted"`
}
