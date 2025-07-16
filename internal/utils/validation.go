package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ExtractValidationErrors extracts validation errors from Gin binding errors
func ExtractValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			fieldName := PascalToCamel(fieldError.Field())

			switch fieldError.Tag() {
			case "required":
				errors[fieldName] = getFieldDisplayName(fieldError.Field()) + "為必填項目"
			case "min":
				if fieldError.Kind().String() == "string" {
					errors[fieldName] = fmt.Sprintf("%s長度至少需要%s個字元", getFieldDisplayName(fieldError.Field()), fieldError.Param())
				} else {
					errors[fieldName] = fmt.Sprintf("%s最小值為%s", getFieldDisplayName(fieldError.Field()), fieldError.Param())
				}
			case "max":
				if fieldError.Kind().String() == "string" {
					errors[fieldName] = fmt.Sprintf("%s長度最多%s個字元", getFieldDisplayName(fieldError.Field()), fieldError.Param())
				} else {
					errors[fieldName] = fmt.Sprintf("%s最大值為%s", getFieldDisplayName(fieldError.Field()), fieldError.Param())
				}
			case "email":
				errors[fieldName] = getFieldDisplayName(fieldError.Field()) + "格式不正確"
			case "numeric":
				errors[fieldName] = getFieldDisplayName(fieldError.Field()) + "必須為數字"
			default:
				errors[fieldName] = getFieldDisplayName(fieldError.Field()) + "格式不正確"
			}
		}
	} else {
		// Handle other types of binding errors (e.g., JSON syntax errors)
		errors["request"] = "請求格式錯誤"
	}

	return errors
}

// getFieldDisplayName returns Chinese display name for field
func getFieldDisplayName(fieldName string) string {
	fieldNames := map[string]string{
		"Username": "帳號",
		"Password": "密碼",
		"Email":    "Email",
		"Role":     "角色",
		"StoreIDs": "門市清單",
		"StoreID":  "門市ID",
	}

	if displayName, exists := fieldNames[fieldName]; exists {
		return displayName
	}
	return fieldName
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	_, ok := err.(validator.ValidationErrors)
	return ok
}
