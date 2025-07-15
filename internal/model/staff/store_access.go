package staff

import "github.com/tkoleo84119/nail-salon-backend/internal/model/common"

// CreateStoreAccessRequest represents the request for creating store access
type CreateStoreAccessRequest struct {
	StoreID int64 `json:"store_id" binding:"required"`
}

// CreateStoreAccessResponse represents the response for creating store access
type CreateStoreAccessResponse struct {
	StaffUserID string         `json:"staff_user_id"`
	StoreList   []common.Store `json:"store_list"`
}