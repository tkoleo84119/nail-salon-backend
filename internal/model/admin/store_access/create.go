package adminStoreAccess

import "github.com/tkoleo84119/nail-salon-backend/internal/model/common"

type CreateRequest struct {
	StoreID string `json:"storeId" binding:"required"`
}

type CreateResponse struct {
	StoreList []common.Store `json:"storeList"`
}
