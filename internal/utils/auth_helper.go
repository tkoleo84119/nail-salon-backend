package utils

import (
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

// CheckOneStoreAccess checks if the staff has access to the store
func CheckOneStoreAccess(storeID int64, staffContext common.StaffContext) (bool, error) {
	var storeAccess []int64
	for _, store := range staffContext.StoreList {
		storeAccess = append(storeAccess, store.ID)
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

// CheckStoreAccess checks if the staff has access to the store
func CheckStoreAccess(storeID int64, storeIDs []int64) (bool, error) {
	hasAccess := false
	for _, id := range storeIDs {
		if id == storeID {
			hasAccess = true
			break
		}
	}
	return hasAccess, nil
}
