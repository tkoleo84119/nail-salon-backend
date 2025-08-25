package booking

import "github.com/jackc/pgx/v5/pgtype"

type CreateRequest struct {
	StoreId       string   `json:"storeId" binding:"required"`
	StylistId     string   `json:"stylistId" binding:"required"`
	TimeSlotId    string   `json:"timeSlotId" binding:"required"`
	MainServiceId string   `json:"mainServiceId" binding:"required"`
	SubServiceIds []string `json:"subServiceIds" binding:"omitempty,max=5"`
	IsChatEnabled *bool    `json:"isChatEnabled" binding:"omitempty"`
	Note          *string  `json:"note" binding:"omitempty,max=255"`
}

type CreateParsedRequest struct {
	StoreId       int64
	StylistId     int64
	TimeSlotId    int64
	MainServiceId int64
	SubServiceIds []int64
	IsChatEnabled *bool
	Note          *string
}

type CreateResponse struct {
	ID              string   `json:"id"`
	StoreId         string   `json:"storeId"`
	StoreName       string   `json:"storeName"`
	StylistId       string   `json:"stylistId"`
	StylistName     string   `json:"stylistName"`
	Date            string   `json:"date"`
	TimeSlotId      string   `json:"timeSlotId"`
	StartTime       string   `json:"startTime"`
	EndTime         string   `json:"endTime"`
	MainServiceName string   `json:"mainServiceName"`
	SubServiceNames []string `json:"subServiceNames"`
	IsChatEnabled   bool     `json:"isChatEnabled"`
	Note            string   `json:"note"`
	Status          string   `json:"status"`
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
}

type CreateBookingServiceInfo struct {
	ServiceId     int64
	ServiceName   string
	IsMainService bool
	Price         pgtype.Numeric
}
