package adminExpense

type GetResponse struct {
	ID           string              `json:"id"`
	Supplier     *GetExpenseSupplier `json:"supplier,omitempty"`
	Payer        *GetExpensePayer    `json:"payer,omitempty"`
	Category     string              `json:"category"`
	Amount       int                 `json:"amount"`
	ExpenseDate  string              `json:"expenseDate"`
	Note         string              `json:"note"`
	IsReimbursed *bool               `json:"isReimbursed,omitempty"`
	ReimbursedAt *string             `json:"reimbursedAt,omitempty"`
	CreatedAt    string              `json:"createdAt"`
	UpdatedAt    string              `json:"updatedAt"`
	Items        []GetExpenseItem    `json:"items"`
}

type GetExpenseSupplier struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetExpensePayer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetExpenseItem struct {
	ID              string                `json:"id"`
	Product         GetExpenseItemProduct `json:"product"`
	Quantity        int                   `json:"quantity"`
	TotalPrice      int                   `json:"totalPrice"`
	ExpirationDate  *string               `json:"expirationDate,omitempty"`
	IsArrived       bool                  `json:"isArrived"`
	ArrivalDate     *string               `json:"arrivalDate,omitempty"`
	StorageLocation string                `json:"storageLocation"`
	Note            string                `json:"note"`
}

type GetExpenseItemProduct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
