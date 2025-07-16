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