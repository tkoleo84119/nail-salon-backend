package adminAccountTransaction

type GetAllRequest struct {
	Limit  *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
}

type GetAllParsedRequest struct {
	Limit  int
	Offset int
}

type GetAllResponse struct {
	Total int          `json:"total"`
	Items []GetAllItem `json:"items"`
}

type GetAllItem struct {
	ID              string `json:"id"`
	TransactionDate string `json:"transactionDate"`
	Type            string `json:"type"`
	Amount          int64  `json:"amount"`
	Balance         int64  `json:"balance"`
	Note            string `json:"note"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}
