package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// ExtractValidationErrors extracts validation errors from Gin binding errors
func ExtractValidationErrors(err error) map[string]string {
	errs := make(map[string]string)

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fieldError := range ve {
			fieldName := PascalToCamel(fieldError.Field())

			switch fieldError.Tag() {
			case "required":
				errs[fieldName] = fieldName + "為必填項目"
			case "min":
				if fieldError.Kind().String() == "string" {
					errs[fieldName] = fmt.Sprintf("%s長度至少需要%s個字元", fieldName, fieldError.Param())
				} else if fieldError.Kind().String() == "slice" {
					errs[fieldName] = fmt.Sprintf("%s至少需要%s個項目", fieldName, fieldError.Param())
				} else {
					errs[fieldName] = fmt.Sprintf("%s最小值為%s", fieldName, fieldError.Param())
				}
			case "max":
				if fieldError.Kind().String() == "string" {
					errs[fieldName] = fmt.Sprintf("%s長度最多%s個字元", fieldName, fieldError.Param())
				} else if fieldError.Kind().String() == "slice" {
					errs[fieldName] = fmt.Sprintf("%s最多只能有%s個項目", fieldName, fieldError.Param())
				} else {
					errs[fieldName] = fmt.Sprintf("%s最大值為%s", fieldName, fieldError.Param())
				}
			case "email":
				errs[fieldName] = fieldName + "格式不是有效的Email"
			case "numeric":
				errs[fieldName] = fieldName + "必須為數字"
			case "boolean":
				errs[fieldName] = fieldName + "必須為布林值"
			case "oneof":
				errs[fieldName] = fieldName + "只可以傳入特定值"
			case "taiwanlandline":
				errs[fieldName] = fieldName + "必須為有效的台灣市話號碼格式 (例: 02-12345678)"
			case "taiwanmobile":
				errs[fieldName] = fieldName + "必須為有效的台灣手機號碼格式 (例: 09xxxxxxxx)"
			default:
				errs[fieldName] = fieldName + "格式不正確"
			}
		}
	}

	var synErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError
	switch {
	case errors.As(err, &synErr), errors.As(err, &typeErr):
		errs["body"] = "JSON 格式錯誤"
	default:
		errs["params"] = "參數格式錯誤"
	}
	return errs
}

// ValidateTaiwanLandline validates Taiwan landline phone numbers
// Format: 0X-XXXXXXXX or 0X-XXXXXXX where X is area code (2-8) and phone number
func ValidateTaiwanLandline(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// Allow empty string (use omitempty in binding tag if optional)
	if phone == "" {
		return true
	}

	// Taiwan landline format: 0X-XXXXXXXX or 0X-XXXXXXX
	// Area codes: 02, 03, 04, 05, 06, 07, 08, 089
	// Phone number: 7-8 digits
	pattern := `^0[2-8]-\d{7,8}$|^089-\d{6}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// ValidateTaiwanMobile validates Taiwan mobile phone numbers
// Format: 09XXXXXXXX where X is any digit
func ValidateTaiwanMobile(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// Allow empty string (use omitempty in binding tag if optional)
	if phone == "" {
		return true
	}

	// Taiwan mobile format: 09XXXXXXXX (10 digits total)
	pattern := `^09\d{8}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}
