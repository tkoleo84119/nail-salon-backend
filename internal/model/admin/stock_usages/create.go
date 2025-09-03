package adminStockUsages

import "time"

type CreateRequest struct {
	ProductID    string  `json:"productId" binding:"required"`
	Quantity     *int64  `json:"quantity" binding:"required,min=1,max=1000000"`
	Expiration   *string `json:"expiration" binding:"omitempty"`
	UsageStarted string  `json:"usageStarted" binding:"required"`
}

type CreateParsedRequest struct {
	ProductID    int64
	Quantity     int
	Expiration   *time.Time
	UsageStarted time.Time
}

type CreateResponse struct {
	ID string `json:"id"`
}
