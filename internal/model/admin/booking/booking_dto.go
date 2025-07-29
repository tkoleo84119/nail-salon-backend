package adminBooking

import "github.com/tkoleo84119/nail-salon-backend/internal/model/common"

// CreateBookingRequest represents the request to create a booking for admin
type CreateBookingRequest struct {
	CustomerID    string   `json:"customerId" binding:"required"`
	TimeSlotID    string   `json:"timeSlotId" binding:"required"`
	MainServiceID string   `json:"mainServiceId" binding:"required"`
	SubServiceIDs []string `json:"subServiceIds" binding:"omitempty,max=10"`
	IsChatEnabled bool     `json:"isChatEnabled"`
	Note          *string  `json:"note,omitempty" binding:"omitempty,max=200"`
}

// CreateBookingResponse represents the response after creating a booking for admin
type CreateBookingResponse struct {
	ID              string   `json:"id"`
	StoreID         string   `json:"storeId"`
	StoreName       string   `json:"storeName"`
	StylistID       string   `json:"stylistId"`
	StylistName     string   `json:"stylistName"`
	Date            string   `json:"date"`
	TimeSlotID      string   `json:"timeSlotId"`
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

// GetBookingListRequest represents the request for getting booking list
type GetBookingListRequest struct {
	StylistID *string `form:"stylistId" binding:"omitempty"`
	StartDate *string `form:"startDate" binding:"omitempty"`
	EndDate   *string `form:"endDate" binding:"omitempty"`
	Limit     *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset    *int    `form:"offset" binding:"omitempty,min=0"`
}

// GetBookingListResponse represents the response for booking list
type GetBookingListResponse common.ListResponse[BookingListItemDTO]

// BookingListItemDTO represents a booking item in the list
type BookingListItemDTO struct {
	ID          string                 `json:"id"`
	Customer    BookingCustomerDTO     `json:"customer"`
	Stylist     BookingStylistDTO      `json:"stylist"`
	TimeSlot    BookingTimeSlotDTO     `json:"timeSlot"`
	MainService BookingMainServiceDTO  `json:"mainService"`
	SubServices []BookingSubServiceDTO `json:"subServices"`
	Status      string                 `json:"status"`
}

// BookingCustomerDTO represents customer information in booking
type BookingCustomerDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// BookingStylistDTO represents stylist information in booking
type BookingStylistDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// BookingTimeSlotDTO represents time slot information in booking
type BookingTimeSlotDTO struct {
	ID        string `json:"id"`
	WorkDate  string `json:"workDate"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// BookingMainServiceDTO represents main service information in booking
type BookingMainServiceDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// BookingSubServiceDTO represents sub service information in booking
type BookingSubServiceDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UpdateBookingByStaffRequest represents the request to update a booking by staff
type UpdateBookingByStaffRequest struct {
	TimeSlotID    *string  `json:"timeSlotId,omitempty"`
	MainServiceID *string  `json:"mainServiceId,omitempty"`
	SubServiceIDs []string `json:"subServiceIds,omitempty" binding:"omitempty,max=5"`
	IsChatEnabled *bool    `json:"isChatEnabled,omitempty"`
	Note          *string  `json:"note,omitempty" binding:"omitempty,max=200"`
}

// UpdateBookingByStaffResponse represents the response after updating a booking by staff
type UpdateBookingByStaffResponse struct {
	ID string `json:"id"`
}

func (r UpdateBookingByStaffRequest) HasUpdates() bool {
	return r.TimeSlotID != nil || r.MainServiceID != nil || r.SubServiceIDs != nil || r.IsChatEnabled != nil || r.Note != nil
}

func (r UpdateBookingByStaffRequest) HasTimeSlotUpdate() bool {
	return r.TimeSlotID != nil || r.MainServiceID != nil || r.SubServiceIDs != nil
}

func (r UpdateBookingByStaffRequest) IsTimeSlotUpdateComplete() bool {
	if !r.HasTimeSlotUpdate() {
		return false
	}

	if r.TimeSlotID == nil || r.MainServiceID == nil || r.SubServiceIDs == nil {
		return false
	}

	return true
}

// -------------------------------------------------------------------------------------

// CancelBookingRequest represents the request to cancel a booking by staff
type CancelBookingRequest struct {
	Status       string  `json:"status" binding:"required,oneof=CANCELLED NO_SHOW"`
	CancelReason *string `json:"cancelReason,omitempty" binding:"omitempty,max=100"`
}

// CancelBookingResponse represents the response after cancelling a booking by staff
type CancelBookingResponse struct {
	ID string `json:"id"`
}
