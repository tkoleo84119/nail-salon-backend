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
				errors[fieldName] = fieldName + "為必填項目"
			case "min":
				if fieldError.Kind().String() == "string" {
					errors[fieldName] = fmt.Sprintf("%s長度至少需要%s個字元", fieldName, fieldError.Param())
				} else if fieldError.Kind().String() == "slice" {
					errors[fieldName] = fmt.Sprintf("%s至少需要%s個項目", fieldName, fieldError.Param())
				} else {
					errors[fieldName] = fmt.Sprintf("%s最小值為%s", fieldName, fieldError.Param())
				}
			case "max":
				if fieldError.Kind().String() == "string" {
					errors[fieldName] = fmt.Sprintf("%s長度最多%s個字元", fieldName, fieldError.Param())
				} else if fieldError.Kind().String() == "slice" {
					errors[fieldName] = fmt.Sprintf("%s最多只能有%s個項目", fieldName, fieldError.Param())
				} else {
					errors[fieldName] = fmt.Sprintf("%s最大值為%s", fieldName, fieldError.Param())
				}
			case "email":
				errors[fieldName] = fieldName + "格式不是有效的Email"
			case "numeric":
				errors[fieldName] = fieldName + "必須為數字"
			case "boolean":
				errors[fieldName] = fieldName + "必須為布林值"
			case "oneof":
				errors[fieldName] = fieldName + "只可以傳入特定值"
			default:
				errors[fieldName] = fieldName + "格式不正確"
			}
		}
	} else {
		// Handle other types of binding errors (e.g., JSON syntax errors)
		errors["request"] = "JSON格式錯誤"
	}

	return errors
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	_, ok := err.(validator.ValidationErrors)
	return ok
}
