package adminReport

import "time"

type GetStoreExpenseRequest struct {
	StartDate string `form:"startDate" binding:"required"`
	EndDate   string `form:"endDate" binding:"required"`
}

type GetStoreExpenseParsedRequest struct {
	StartDate time.Time
	EndDate   time.Time
}

type GetStoreExpenseResponse struct {
	StartDate string                `json:"startDate"`
	EndDate   string                `json:"endDate"`
	Summary   ExpenseReportSummary  `json:"summary"`
	Category  []CategoryExpenseStat `json:"category"`
	Supplier  []SupplierExpenseStat `json:"supplier"`
	Payer     []PayerExpenseStat    `json:"payer"`
}

type ExpenseReportSummary struct {
	TotalCount       int   `json:"totalCount"`       // 總支出筆數
	TotalAmount      int64 `json:"totalAmount"`      // 總支出金額(含其他費用)
	AdvanceAmount    int64 `json:"advanceAmount"`    // 代墊總金額(含其他費用)
	ReimbursedAmount int64 `json:"reimbursedAmount"` // 已結清代墊金額
	PendingAmount    int64 `json:"pendingAmount"`    // 待結清代墊金額
}

type CategoryExpenseStat struct {
	Category string `json:"category"` // 類別名稱
	Count    int    `json:"count"`    // 該類別支出筆數
	Amount   int64  `json:"amount"`   // 該類別總金額(含其他費用)
}

type SupplierExpenseStat struct {
	SupplierID   string `json:"supplierId"`   // 供應商ID
	SupplierName string `json:"supplierName"` // 供應商名稱
	Count        int    `json:"count"`        // 該供應商支出筆數
	Amount       int64  `json:"amount"`       // 該供應商總金額(含其他費用)
}

type PayerExpenseStat struct {
	PayerID          string `json:"payerId"`          // 代墊人ID
	PayerName        string `json:"payerName"`        // 代墊人姓名
	AdvanceCount     int    `json:"advanceCount"`     // 代墊筆數
	AdvanceAmount    int64  `json:"advanceAmount"`    // 代墊總金額(含其他費用)
	ReimbursedAmount int64  `json:"reimbursedAmount"` // 已結清金額
	PendingAmount    int64  `json:"pendingAmount"`    // 待結清金額
}
