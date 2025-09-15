package utils

import (
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

// CheckStoreAccess checks if the staff has access to the store
func CheckStoreAccess(storeID int64, storeIDs []int64, role string) error {
	if role == common.RoleSuperAdmin {
		return nil
	}

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
