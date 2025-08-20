package booking

import "github.com/jackc/pgx/v5/pgtype"

type UpdateRequest struct {
	StoreId       *string   `json:"storeId,omitempty"`
	StylistId     *string   `json:"stylistId,omitempty"`
	TimeSlotId    *string   `json:"timeSlotId,omitempty"`
	MainServiceId *string   `json:"mainServiceId,omitempty"`
	SubServiceIds *[]string `json:"subServiceIds" binding:"omitempty,max=5"`
	IsChatEnabled *bool     `json:"isChatEnabled,omitempty"`
	Note          *string   `json:"note" binding:"omitempty,max=500"`
}

type UpdateParsedRequest struct {
	StoreId       *int64
	StylistId     *int64
	TimeSlotId    *int64
	MainServiceId *int64
	SubServiceIds *[]int64
	IsChatEnabled *bool
	Note          *string
}

type UpdateResponse struct {
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

type UpdateBookingServiceInfo struct {
	ServiceId     int64
	ServiceName   string
	IsMainService bool
	Price         pgtype.Numeric
}

func (r UpdateRequest) HasUpdates() bool {
	return r.StoreId != nil || r.StylistId != nil || r.TimeSlotId != nil ||
		r.MainServiceId != nil || r.SubServiceIds != nil || r.IsChatEnabled != nil || r.Note != nil
}

func (r UpdateRequest) HasTimeSlotUpdate() bool {
	return r.StoreId != nil || r.StylistId != nil || r.TimeSlotId != nil || r.MainServiceId != nil || r.SubServiceIds != nil
}

func (r UpdateRequest) IsTimeSlotUpdateComplete() bool {
	if !r.HasTimeSlotUpdate() {
		return true
	}
	return r.StoreId != nil && r.StylistId != nil && r.TimeSlotId != nil && r.MainServiceId != nil && r.SubServiceIds != nil
}

func (r UpdateParsedRequest) HasUpdates() bool {
	return r.StoreId != nil || r.StylistId != nil || r.TimeSlotId != nil ||
		r.MainServiceId != nil || r.SubServiceIds != nil || r.IsChatEnabled != nil || r.Note != nil
}

func (r UpdateParsedRequest) HasTimeSlotUpdate() bool {
	return r.StoreId != nil || r.StylistId != nil || r.TimeSlotId != nil || r.MainServiceId != nil || r.SubServiceIds != nil
}

func (r UpdateParsedRequest) IsTimeSlotUpdateComplete() bool {
	if !r.HasTimeSlotUpdate() {
		return true
	}
	return r.StoreId != nil && r.StylistId != nil && r.TimeSlotId != nil && r.MainServiceId != nil && r.SubServiceIds != nil
}
