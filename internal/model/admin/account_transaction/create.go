package adminAccountTransaction

import "time"

type CreateRequest struct {
	TransactionDate string  `json:"transactionDate" binding:"required"`
	Type            string  `json:"type" binding:"required,oneof=INCOME EXPENSE"`
	Amount          int     `json:"amount" binding:"required,min=1,max=1000000"`
	Note            *string `json:"note" binding:"omitempty,max=255"`
}

type CreateParsedRequest struct {
	TransactionDate time.Time
	Type            string
	Amount          int
	Note            *string
}

type CreateResponse struct {
	ID string `json:"id"`
}
