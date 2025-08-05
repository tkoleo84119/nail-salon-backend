package adminTimeSlotTemplate

import "github.com/tkoleo84119/nail-salon-backend/internal/model/common"

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
