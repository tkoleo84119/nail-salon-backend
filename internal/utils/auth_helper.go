package utils

import (
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

// CheckOneStoreAccess checks if the staff has access to the store
func CheckOneStoreAccess(storeID int64, staffContext common.StaffContext) (bool, error) {
	var storeAccess []int64
	for _, store := range staffContext.StoreList {
		storeId, err := ParseID(store.ID)
		if err != nil {
			return false, err
		}
		storeAccess = append(storeAccess, storeId)
	}

	hasAccess := false
	for _, storeId := range storeAccess {
		if storeId == storeID {
			hasAccess = true
			break
		}
	}
	return hasAccess, nil
}
