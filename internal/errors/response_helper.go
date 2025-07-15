package errors

import "github.com/gin-gonic/gin"

// RespondWithError sends a JSON error response using the error manager
func RespondWithError(c *gin.Context, code string, fieldErrors map[string]string) {
	errorManager := GetManager()
	status := errorManager.GetStatus(code)
	errorResponse := errorManager.GetErrorResponse(code, fieldErrors)
	c.JSON(status, errorResponse)
}

// RespondWithErrorDetails sends a JSON error response with optional development details
func RespondWithErrorDetails(c *gin.Context, code string, fieldErrors map[string]string, details string) {
	errorManager := GetManager()
	status := errorManager.GetStatus(code)
	errorResponse := errorManager.GetErrorResponse(code, fieldErrors, details)
	c.JSON(status, errorResponse)
}

// RespondWithServiceError handles service errors and sends appropriate JSON response
func RespondWithServiceError(c *gin.Context, err error) {
	if code, ok := IsServiceError(err); ok {
		// Extract detailed error message for development
		var details string
		if serviceErr, ok := err.(*ServiceError); ok {
			details = serviceErr.Error()
		}
		RespondWithErrorDetails(c, code, nil, details)
	} else {
		// For non-service errors, provide the error details in development
		RespondWithErrorDetails(c, SysInternalError, nil, err.Error())
	}
}

// AbortWithError sends a JSON error response and aborts the request (for middleware)
func AbortWithError(c *gin.Context, code string, fieldErrors map[string]string) {
	RespondWithError(c, code, fieldErrors)
	c.Abort()
}