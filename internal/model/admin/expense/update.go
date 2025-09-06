package adminExpense

import "time"

type UpdateRequest struct {
	SupplierID   *string `json:"supplierId" binding:"omitempty"`
	Category     *string `json:"category" binding:"omitempty,noBlank,max=100"`
	Amount       *int64  `json:"amount" binding:"omitempty,min=0,max=1000000"`
	OtherFee     *int64  `json:"otherFee" binding:"omitempty,min=0,max=1000000"`
	ExpenseDate  *string `json:"expenseDate" binding:"omitempty"`
	Note         *string `json:"note" binding:"omitempty,max=255"`
	PayerID      *string `json:"payerId" binding:"omitempty"`
	IsReimbursed *bool   `json:"isReimbursed" binding:"omitempty"`
	ReimbursedAt *string `json:"reimbursedAt" binding:"omitempty"`
}

type UpdateParsedRequest struct {
	SupplierID    *int64
	Category      *string
	Amount        *int64
	OtherFee      *int64
	ExpenseDate   *time.Time
	Note          *string
	PayerID       *int64
	PayerIDIsNone *bool
	IsReimbursed  *bool
	ReimbursedAt  *time.Time
}

type UpdateResponse struct {
	ID string `json:"id"`
}

func (r *UpdateRequest) HasUpdates() bool {
	return r.SupplierID != nil || r.Category != nil || r.Amount != nil || r.OtherFee != nil || r.ExpenseDate != nil || r.Note != nil || r.PayerID != nil || r.IsReimbursed != nil || r.ReimbursedAt != nil
}
