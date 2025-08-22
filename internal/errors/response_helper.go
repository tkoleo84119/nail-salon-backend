package errors

import "github.com/gin-gonic/gin"

// RespondWithError sends a JSON error response using the error manager
func RespondWithError(c *gin.Context, errors []ErrorItem) {
	errorManager := GetManager()

	status := 400
	if len(errors) > 0 {
		status = errorManager.GetStatusByCode(errors[0].Code)
	}

	errorResponse := errorManager.GetErrorResponse(errors)
	c.JSON(status, errorResponse)
}

// RespondWithErrorDetails sends a JSON error response with optional development details
func RespondWithErrorDetails(c *gin.Context, constantName string, fieldErrors map[string]string, details string) {
	errorManager := GetManager()
	var errors []ErrorItem

	if len(fieldErrors) > 0 {
		// if there are field errors, create error items for each field
		for field, message := range fieldErrors {
			item := ErrorItem{
				Code:    errorManager.GetCode(constantName),
				Message: message,
				Field:   field,
			}
			errors = append(errors, item)
		}
	} else {
		// if there are no field errors, create a single error item
		errorItem := errorManager.CreateErrorItem(constantName, "", nil)
		errors = append(errors, errorItem)
	}

	status := errorManager.GetStatus(constantName)
	errorResponse := errorManager.GetErrorResponse(errors)

	// add details in debug mode
	if details != "" && errorManager.isDebugMode() {
		errorResponse["dev_details"] = details
	}

	c.JSON(status, errorResponse)
}

// RespondWithServiceError handles service errors and sends appropriate JSON response
func RespondWithServiceError(c *gin.Context, err error) {
	if code, ok := IsServiceError(err); ok {
		RespondWithErrorDetails(c, code, nil, err.Error())
	} else {
		RespondWithErrorDetails(c, SysInternalError, nil, err.Error())
	}
}

// RespondWithValidationErrors sends a JSON error response for validation errors
func RespondWithValidationErrors(c *gin.Context, validationErrors []ErrorItem) {
	RespondWithError(c, validationErrors)
}

// AbortWithError sends a JSON error response and aborts the request (for middleware)
func AbortWithError(c *gin.Context, constantName string, fieldErrors map[string]string) {
	errorManager := GetManager()
	var errors []ErrorItem

	if len(fieldErrors) > 0 {
		// if there are field errors, create error items for each field
		for field, message := range fieldErrors {
			item := ErrorItem{
				Code:    errorManager.GetCode(constantName),
				Message: message,
				Field:   field,
			}
			errors = append(errors, item)
		}
	} else {
		// if there are no field errors, create a single error item
		errors = append(errors, errorManager.CreateErrorItem(constantName, "", nil))
	}

	RespondWithError(c, errors)
	c.Abort()
}
