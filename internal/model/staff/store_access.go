package staff

import "github.com/tkoleo84119/nail-salon-backend/internal/model/common"

// CreateStoreAccessRequest represents the request for creating store access
type CreateStoreAccessRequest struct {
	StoreID string `json:"storeId" binding:"required"`
}

// CreateStoreAccessResponse represents the response for creating store access
type CreateStoreAccessResponse struct {
	StaffUserID string         `json:"staffUserId"`
	StoreList   []common.Store `json:"storeList"`
}

// DeleteStoreAccessBulkRequest represents the request for deleting store access
type DeleteStoreAccessBulkRequest struct {
	StoreIDs []string `json:"storeIds" binding:"required,min=1,max=100"`
}

// DeleteStoreAccessBulkResponse represents the response for deleting store access
type DeleteStoreAccessBulkResponse struct {
	StaffUserID string         `json:"staffUserId"`
	StoreList   []common.Store `json:"storeList"`
}
