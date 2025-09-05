package adminExpense

type GetAllRequest struct {
	Category     *string `form:"category" binding:"omitempty,noBlank,max=100"`
	SupplierID   *string `form:"supplierId" binding:"omitempty"`
	PayerID      *string `form:"payerId" binding:"omitempty"`
	IsReimbursed *bool   `form:"isReimbursed" binding:"omitempty"`
	Limit        *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset       *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
	Sort         *string `form:"sort" binding:"omitempty"`
}

type GetAllParsedRequest struct {
	Category     *string
	SupplierID   *int64
	PayerID      *int64
	IsReimbursed *bool
	Limit        int
	Offset       int
	Sort         []string
}

type GetAllResponse struct {
	Total int                 `json:"total"`
	Items []GetAllExpenseItem `json:"items"`
}

type GetAllExpenseItem struct {
	ID           string                    `json:"id"`
	Supplier     GetAllExpenseSupplierItem `json:"supplier"`
	Payer        *GetAllExpensePayerItem   `json:"payer,omitempty"`
	Category     string                    `json:"category"`
	Amount       int                       `json:"amount"`
	ExpenseDate  string                    `json:"expenseDate"`
	Note         string                    `json:"note"`
	IsReimbursed *bool                     `json:"isReimbursed,omitempty"`
	ReimbursedAt *string                   `json:"reimbursedAt,omitempty"`
	CreatedAt    string                    `json:"createdAt"`
	UpdatedAt    string                    `json:"updatedAt"`
}

type GetAllExpenseSupplierItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetAllExpensePayerItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
