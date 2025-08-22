package utils

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
