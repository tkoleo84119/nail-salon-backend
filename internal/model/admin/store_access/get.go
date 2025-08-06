package adminStoreAccess

import "github.com/tkoleo84119/nail-salon-backend/internal/model/common"

type GetResponse struct {
	StoreList []common.Store `json:"storeList"`
}
