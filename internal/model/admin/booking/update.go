package adminBooking

import "github.com/jackc/pgx/v5/pgtype"

type UpdateRequest struct {
	StylistID     *string   `json:"stylistId"`
	TimeSlotID    *string   `json:"timeSlotId"`
	MainServiceID *string   `json:"mainServiceId"`
	SubServiceIDs *[]string `json:"subServiceIds" binding:"omitempty,max=10"`
	IsChatEnabled *bool     `json:"isChatEnabled" binding:"omitempty"`
	Note          *string   `json:"note" binding:"omitempty,max=255"`
}

type UpdateParsedRequest struct {
	StylistID     *int64
	TimeSlotID    *int64
	MainServiceID *int64
	SubServiceIDs []int64
	IsChatEnabled *bool
	Note          *string
}

type UpdateBookingServiceInfo struct {
	ServiceId     int64
	ServiceName   string
	IsMainService bool
	Price         pgtype.Numeric
}

type UpdateResponse struct {
	ID string `json:"id"`
}

func (r UpdateRequest) HasUpdates() bool {
	return r.TimeSlotID != nil || r.MainServiceID != nil || r.SubServiceIDs != nil || r.IsChatEnabled != nil || r.Note != nil
}

func (r UpdateRequest) HasTimeSlotUpdate() bool {
	return r.StylistID != nil || r.TimeSlotID != nil || r.MainServiceID != nil || r.SubServiceIDs != nil
}

func (r UpdateRequest) IsTimeSlotUpdateComplete() bool {
	if !r.HasTimeSlotUpdate() {
		return false
	}

	if r.StylistID == nil || r.TimeSlotID == nil || r.MainServiceID == nil || r.SubServiceIDs == nil {
		return false
	}

	return true
}

func (r UpdateParsedRequest) HasUpdates() bool {
	return r.StylistID != nil || r.TimeSlotID != nil || r.MainServiceID != nil || r.SubServiceIDs != nil || r.IsChatEnabled != nil || r.Note != nil
}

func (r UpdateParsedRequest) HasTimeSlotUpdate() bool {
	return r.StylistID != nil || r.TimeSlotID != nil || r.MainServiceID != nil || r.SubServiceIDs != nil
}

func (r UpdateParsedRequest) IsTimeSlotUpdateComplete() bool {
	if !r.HasTimeSlotUpdate() {
		return false
	}

	if r.StylistID == nil || r.TimeSlotID == nil || r.MainServiceID == nil || r.SubServiceIDs == nil {
		return false
	}

	return true
}
