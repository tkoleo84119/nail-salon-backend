package storeAccess

import "github.com/tkoleo84119/nail-salon-backend/internal/model/common"

// DeleteStoreAccessRequest represents the request for deleting store access
type DeleteStoreAccessRequest struct {
	StoreIDs []string `json:"storeIds" binding:"required,min=1,max=100"`
}

// DeleteStoreAccessResponse represents the response for deleting store access
type DeleteStoreAccessResponse struct {
	StaffUserID string         `json:"staffUserId"`
	StoreList   []common.Store `json:"storeList"`
}
