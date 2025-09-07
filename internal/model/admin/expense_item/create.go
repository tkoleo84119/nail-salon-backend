package adminExpenseItem

import "time"

type CreateRequest struct {
	ProductID       string  `json:"productId" binding:"required"`
	Quantity        *int64  `json:"quantity" binding:"omitempty,min=0,max=1000000"`
	Price           *int64  `json:"price" binding:"omitempty,min=0,max=1000000"`
	ExpirationDate  *string `json:"expirationDate" binding:"omitempty"`
	IsArrived       *bool   `json:"isArrived" binding:"omitempty"`
	ArrivalDate     *string `json:"arrivalDate" binding:"omitempty"`
	StorageLocation *string `json:"storageLocation" binding:"omitempty,max=100"`
	Note            *string `json:"note" binding:"omitempty,max=255"`
}

type CreateParsedRequest struct {
	ProductID       int64
	Quantity        int64
	Price           int64
	ExpirationDate  *time.Time
	IsArrived       bool
	ArrivalDate     *time.Time
	StorageLocation *string
	Note            *string
}

type CreateResponse struct {
	ID string `json:"id"`
}
