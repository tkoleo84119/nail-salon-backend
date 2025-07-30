package utils

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	ownErrors "github.com/tkoleo84119/nail-salon-backend/internal/errors"
)

// ExtractValidationErrors extracts validation errors from Gin binding errors
func ExtractValidationErrors(err error) []ownErrors.ErrorItem {
	var errorItems []ownErrors.ErrorItem
	errorManager := ownErrors.GetManager()

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fieldError := range ve {
			fieldName := PascalToCamel(fieldError.Field())
			var constantName string
			params := map[string]string{
				"field": fieldName,
			}

			switch fieldError.Tag() {
			case "required":
				constantName = ownErrors.ValFieldRequired
			case "min":
				if fieldError.Kind().String() == "string" {
					constantName = ownErrors.ValFieldStringMinLength
					params["param"] = fieldError.Param()
				} else if fieldError.Kind().String() == "slice" {
					constantName = ownErrors.ValFieldArrayMinLength
					params["param"] = fieldError.Param()
				} else {
					constantName = ownErrors.ValFieldMinNumber
					params["param"] = fieldError.Param()
				}
			case "max":
				if fieldError.Kind().String() == "string" {
					constantName = ownErrors.ValFieldStringMaxLength
					params["param"] = fieldError.Param()
				} else if fieldError.Kind().String() == "slice" {
					constantName = ownErrors.ValFieldArrayMaxLength
					params["param"] = fieldError.Param()
				} else {
					constantName = ownErrors.ValFieldMaxNumber
					params["param"] = fieldError.Param()
				}
			case "email":
				constantName = ownErrors.ValFieldInvalidEmail
			case "numeric":
				constantName = ownErrors.ValFieldNumeric
			case "boolean":
				constantName = ownErrors.ValFieldBoolean
			case "oneof":
				constantName = ownErrors.ValFieldOneof
				params["param"] = fieldError.Param()
			case "taiwanlandline":
				constantName = ownErrors.ValFieldTaiwanLandline
			case "taiwanmobile":
				constantName = ownErrors.ValFieldTaiwanMobile
			default:
				constantName = ownErrors.ValInputValidationFailed
			}

			errorItem := errorManager.CreateErrorItem(constantName, fieldName, params)
			errorItems = append(errorItems, errorItem)
		}

		return errorItems
	}

	var synErr *json.SyntaxError
	if errors.As(err, &synErr) {
		errorItem := errorManager.CreateErrorItem(
			ownErrors.ValJsonFormat,
			"body",
			nil,
		)
		errorItems = append(errorItems, errorItem)
		return errorItems
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) || strings.Contains(err.Error(), "parsing time") {
		fieldName := typeErr.Field
		errorItem := errorManager.CreateErrorItem(
			ownErrors.ValTypeConversionFailed,
			fieldName,
			nil,
		)
		errorItems = append(errorItems, errorItem)
		return errorItems
	}

	if len(errorItems) == 0 {
		errorItem := errorManager.CreateErrorItem(
			ownErrors.ValInputValidationFailed,
			"",
			nil,
		)
		errorItems = append(errorItems, errorItem)
	}

	return errorItems
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
