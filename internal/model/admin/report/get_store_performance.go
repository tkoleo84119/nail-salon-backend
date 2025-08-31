package adminReport

import "time"

type GetStorePerformanceRequest struct {
	StartDate string `form:"startDate" binding:"required"`
	EndDate   string `form:"endDate" binding:"required"`
}

type GetStorePerformanceParsedRequest struct {
	StartDate time.Time
	EndDate   time.Time
}

type GetStorePerformanceResponse struct {
	StartDate         string                       `json:"startDate"`
	EndDate           string                       `json:"endDate"`
	TotalBookings     int                          `json:"totalBookings"`
	CompletedBookings int                          `json:"completedBookings"`
	CancelledBookings int                          `json:"cancelledBookings"`
	NoShowBookings    int                          `json:"noShowBookings"`
	LinePayRevenue    float64                      `json:"linePayRevenue"`
	CashRevenue       float64                      `json:"cashRevenue"`
	TotalPaidAmount   float64                      `json:"totalPaidAmount"`
	TotalServiceTime  int                          `json:"totalServiceTime"`
	Stylists          []GetStorePerformanceStylist `json:"stylists"`
}

type GetStorePerformanceStylist struct {
	StylistID         string  `json:"stylistId"`
	StylistName       string  `json:"stylistName"`
	TotalBookings     int     `json:"totalBookings"`
	CompletedBookings int     `json:"completedBookings"`
	CancelledBookings int     `json:"cancelledBookings"`
	NoShowBookings    int     `json:"noShowBookings"`
	LinePayRevenue    float64 `json:"linePayRevenue"`
	CashRevenue       float64 `json:"cashRevenue"`
	TotalPaidAmount   float64 `json:"totalPaidAmount"`
	TotalServiceTime  int     `json:"totalServiceTime"`
}
