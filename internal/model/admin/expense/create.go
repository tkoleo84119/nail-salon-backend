package adminExpense

import "time"

type CreateRequest struct {
	SupplierID  string              `json:"supplierId" binding:"required"`
	Category    string              `json:"category" binding:"required,noBlank,max=100"`
	Amount      *int64              `json:"amount" binding:"omitempty,min=0,max=1000000"`
	ExpenseDate string              `json:"expenseDate" binding:"required"`
	Note        *string             `json:"note" binding:"omitempty,max=255"`
	PayerID     *string             `json:"payerId" binding:"omitempty"`
	Items       []CreateItemRequest `json:"items" binding:"omitempty,min=1,max=100"`
}

type CreateItemRequest struct {
	ProductID       string  `json:"productId" binding:"required"`
	Quantity        *int64  `json:"quantity" binding:"omitempty,min=0,max=1000000"`
	TotalPrice      *int64  `json:"totalPrice" binding:"omitempty,min=0,max=1000000"`
	ExpirationDate  *string `json:"expirationDate" binding:"omitempty"`
	IsArrived       *bool   `json:"isArrived" binding:"omitempty"`
	ArrivalDate     *string `json:"arrivalDate" binding:"omitempty"`
	StorageLocation *string `json:"storageLocation" binding:"omitempty,max=100"`
	Note            *string `json:"note" binding:"omitempty,max=255"`
}

type CreateParsedRequest struct {
	SupplierID  int64
	Category    string
	Amount      int64
	ExpenseDate time.Time
	Note        *string
	PayerID     *int64
	Items       []CreateItemParsedRequest
}

type CreateItemParsedRequest struct {
	ProductID       int64
	Quantity        int64
	TotalPrice      int64
	ExpirationDate  *time.Time
	IsArrived       bool
	ArrivalDate     *time.Time
	StorageLocation *string
	Note            *string
}

type CreateResponse struct {
	ID string `json:"id"`
}
