package adminReport

import "time"

type GetPerformanceMeRequest struct {
	StartDate string `form:"startDate" binding:"required"`
	EndDate   string `form:"endDate" binding:"required"`
}

type GetPerformanceMeParsedRequest struct {
	StartDate time.Time
	EndDate   time.Time
}

type GetPerformanceMeResponse struct {
	StartDate         string                  `json:"startDate"`
	EndDate           string                  `json:"endDate"`
	TotalBookings     int                     `json:"totalBookings"`
	CompletedBookings int                     `json:"completedBookings"`
	CancelledBookings int                     `json:"cancelledBookings"`
	NoShowBookings    int                     `json:"noShowBookings"`
	LinePayRevenue    float64                 `json:"linePayRevenue"`
	CashRevenue       float64                 `json:"cashRevenue"`
	TransferRevenue   float64                 `json:"transferRevenue"`
	TotalPaidAmount   float64                 `json:"totalPaidAmount"`
	TotalServiceTime  int                     `json:"totalServiceTime"`
	Stores            []GetPerformanceMeStore `json:"stores"`
}

type GetPerformanceMeStore struct {
	StoreID           string  `json:"storeId"`
	StoreName         string  `json:"storeName"`
	TotalBookings     int     `json:"totalBookings"`
	CompletedBookings int     `json:"completedBookings"`
	CancelledBookings int     `json:"cancelledBookings"`
	NoShowBookings    int     `json:"noShowBookings"`
	LinePayRevenue    float64 `json:"linePayRevenue"`
	CashRevenue       float64 `json:"cashRevenue"`
	TransferRevenue   float64 `json:"transferRevenue"`
	TotalPaidAmount   float64 `json:"totalPaidAmount"`
	TotalServiceTime  int     `json:"totalServiceTime"`
}
