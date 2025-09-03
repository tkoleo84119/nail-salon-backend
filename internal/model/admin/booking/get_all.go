package adminBooking

import "time"

type GetAllRequest struct {
	StylistID *string `form:"stylistId" binding:"omitempty"`
	StartDate *string `form:"startDate" binding:"omitempty"`
	EndDate   *string `form:"endDate" binding:"omitempty"`
	Status    *string `form:"status" binding:"omitempty,oneof=SCHEDULED CANCELLED COMPLETED NO_SHOW"`
	Limit     *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset    *int    `form:"offset" binding:"omitempty,min=0"`
	Sort      *string `form:"sort" binding:"omitempty"`
}

type GetAllParsedRequest struct {
	StylistID *int64
	StartDate *time.Time
	EndDate   *time.Time
	Status    *string
	Limit     int
	Offset    int
	Sort      []string
}

type GetAllResponse struct {
	Total int          `json:"total"`
	Items []GetAllItem `json:"items"`
}

type GetAllItem struct {
	ID             string             `json:"id"`
	Customer       GetAllCustomer     `json:"customer"`
	Stylist        GetAllStylist      `json:"stylist"`
	TimeSlot       GetAllTimeSlot     `json:"timeSlot"`
	MainService    GetAllMainService  `json:"mainService"`
	SubServices    []GetAllSubService `json:"subServices"`
	ActualDuration *int32             `json:"actualDuration,omitempty"`
	Status         string             `json:"status"`
}

type GetAllCustomer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetAllStylist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetAllTimeSlot struct {
	ID        string `json:"id"`
	WorkDate  string `json:"workDate"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type GetAllMainService struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetAllSubService struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
