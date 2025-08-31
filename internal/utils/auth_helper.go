package utils

import (
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
)

// CheckStoreAccess checks if the staff has access to the store
func CheckStoreAccess(storeID int64, storeIDs []int64) error {
	hasAccess := false
	for _, id := range storeIDs {
		if id == storeID {
			hasAccess = true
			break
		}
	}

	if !hasAccess {
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	return nil
}
