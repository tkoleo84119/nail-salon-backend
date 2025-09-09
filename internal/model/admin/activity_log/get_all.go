package adminActivityLog

import "github.com/tkoleo84119/nail-salon-backend/internal/model/common"

type GetAllRequest struct {
	Limit *int `form:"limit" binding:"omitempty,min=1,max=50"`
}

type GetAllParsedRequest struct {
	Limit int
}

type GetAllResponse struct {
	Total int                       `json:"total"`
	Items []common.ActivityLogEntry `json:"items"`
}
