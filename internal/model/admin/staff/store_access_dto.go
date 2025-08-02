package adminStaff

import "github.com/tkoleo84119/nail-salon-backend/internal/model/common"

// CreateStoreAccessRequest represents the request for creating store access
type CreateStoreAccessRequest struct {
	StoreID string `json:"storeId" binding:"required"`
}

// CreateStoreAccessResponse represents the response for creating store access
type CreateStoreAccessResponse struct {
	StoreList []common.Store `json:"storeList"`
}

// -------------------------------------------------------------------------------------

// GetStaffStoreAccessResponse represents the response for staff store access list
type GetStaffStoreAccessResponse struct {
	StoreList []common.Store `json:"storeList"`
}

// -------------------------------------------------------------------------------------

// DeleteStoreAccessBulkRequest represents the request for deleting store access
type DeleteStoreAccessBulkRequest struct {
	StoreIDs []string `json:"storeIds" binding:"required,min=1,max=20"`
}

// DeleteStoreAccessBulkResponse represents the response for deleting store access
type DeleteStoreAccessBulkResponse struct {
	StoreList []common.Store `json:"storeList"`
}
