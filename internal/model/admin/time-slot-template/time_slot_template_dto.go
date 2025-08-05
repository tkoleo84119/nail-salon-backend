package adminTimeSlotTemplate

import "github.com/tkoleo84119/nail-salon-backend/internal/model/common"

// TimeSlotItem represents a time slot in a template
type TimeSlotItem struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

// CreateTimeSlotTemplateRequest represents the request to create a time slot template
type CreateTimeSlotTemplateRequest struct {
	Name      string         `json:"name" binding:"required,min=1,max=50"`
	Note      string         `json:"note" binding:"max=100"`
	TimeSlots []TimeSlotItem `json:"timeSlots" binding:"required,min=1,max=50"`
}

// TimeSlotItemResponse represents a time slot in the response
type TimeSlotItemResponse struct {
	ID        string `json:"id"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// CreateTimeSlotTemplateResponse represents the response after creating a time slot template
type CreateTimeSlotTemplateResponse struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Note      string                 `json:"note"`
	TimeSlots []TimeSlotItemResponse `json:"timeSlots"`
}

// -------------------------------------------------------------------------------------

// UpdateTimeSlotTemplateRequest represents the request to update a time slot template
type UpdateTimeSlotTemplateRequest struct {
	Name *string `json:"name,omitempty" binding:"omitempty,min=1,max=50"`
	Note *string `json:"note,omitempty" binding:"omitempty,max=100"`
}

// UpdateTimeSlotTemplateResponse represents the response after updating a time slot template
type UpdateTimeSlotTemplateResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Note string `json:"note"`
}

func (r *UpdateTimeSlotTemplateRequest) HasUpdate() bool {
	return r.Name != nil || r.Note != nil
}

// -------------------------------------------------------------------------------------

// DeleteTimeSlotTemplateResponse represents the response after deleting a time slot template
type DeleteTimeSlotTemplateResponse struct {
	Deleted string `json:"deleted"`
}

// -------------------------------------------------------------------------------------

// GetTimeSlotTemplateListRequest represents the request to get time slot template list
type GetTimeSlotTemplateListRequest struct {
	Name   *string `form:"name" binding:"omitempty,max=100"`
	Limit  *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
	Sort   *string `form:"sort" binding:"omitempty"`
}

type GetTimeSlotTemplateListParsedRequest struct {
	Name   *string
	Limit  int
	Offset int
	Sort   []string
}

// GetTimeSlotTemplateListResponse represents the response with time slot template list
type GetTimeSlotTemplateListResponse common.ListResponse[GetTimeSlotTemplateListItem]

// GetTimeSlotTemplateListItem represents a single time slot template in the list
type GetTimeSlotTemplateListItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Note      string `json:"note"`
	Updater   string `json:"updater"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// -------------------------------------------------------------------------------------

// GetTimeSlotTemplateResponse represents the response for getting a single time slot template
type GetTimeSlotTemplateResponse struct {
	ID        string                        `json:"id"`
	Name      string                        `json:"name"`
	Note      string                        `json:"note"`
	Updater   string                        `json:"updater"`
	CreatedAt string                        `json:"createdAt"`
	UpdatedAt string                        `json:"updatedAt"`
	Items     []GetTimeSlotTemplateItemInfo `json:"items"`
}

// GetTimeSlotTemplateItemInfo represents a time slot template item in the response
type GetTimeSlotTemplateItemInfo struct {
	ID        string `json:"id"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
