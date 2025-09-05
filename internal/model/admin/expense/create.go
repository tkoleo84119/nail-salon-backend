package adminExpense

import "time"

type CreateRequest struct {
	SupplierID  string  `json:"supplierId" binding:"required"`
	Category    string  `json:"category" binding:"required,noBlank,max=100"`
	Amount      *int64  `json:"amount" binding:"omitempty,min=0,max=1000000"`
	ExpenseDate string  `json:"expenseDate" binding:"required"`
	Note        *string `json:"note" binding:"omitempty,max=255"`
	PayerID     *string `json:"payerId" binding:"omitempty"`
}

type CreateParsedRequest struct {
	SupplierID  int64
	Category    string
	Amount      int64
	ExpenseDate time.Time
	Note        *string
	PayerID     *int64
}

type CreateResponse struct {
	ID string `json:"id"`
}
