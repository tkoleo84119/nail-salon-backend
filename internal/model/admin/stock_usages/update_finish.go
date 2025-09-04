package adminStockUsages

import "time"

type UpdateFinishRequest struct {
	UsageEndedAt string `json:"usageEndedAt" binding:"required"`
}

type UpdateFinishParsedRequest struct {
	UsageEndedAt time.Time
}

type UpdateFinishResponse struct {
	ID string `json:"id"`
}
