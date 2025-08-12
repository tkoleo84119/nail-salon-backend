package schedule

import "time"

type GetAllRequest struct {
	StartDate string `form:"startDate" binding:"required"`
	EndDate   string `form:"endDate" binding:"required"`
}

type GetAllParsedRequest struct {
	StartDate time.Time
	EndDate   time.Time
}

type GetAllResponse struct {
	Schedules []GetAllItem `json:"schedules"`
}

type GetAllItem struct {
	ID             string `json:"id"`
	Date           string `json:"date"`
	AvailableSlots int    `json:"availableSlots"`
}
