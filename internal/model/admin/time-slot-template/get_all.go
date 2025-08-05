package adminTimeSlotTemplate

import "github.com/tkoleo84119/nail-salon-backend/internal/model/common"

type GetAllRequest struct {
	Name   *string `form:"name" binding:"omitempty,max=100"`
	Limit  *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
	Sort   *string `form:"sort" binding:"omitempty"`
}

type GetAllParsedRequest struct {
	Name   *string
	Limit  int
	Offset int
	Sort   []string
}

// GetAllResponse represents the response with time slot template list
type GetAllResponse common.ListResponse[GetAllItem]

// GetAllItem represents a single time slot template in the list
type GetAllItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Note      string `json:"note"`
	Updater   string `json:"updater"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
