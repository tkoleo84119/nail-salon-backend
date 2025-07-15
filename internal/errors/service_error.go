package errors

import "errors"

// ServiceError represents a service layer error with an error code
type ServiceError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface
func (e *ServiceError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *ServiceError) Unwrap() error {
	return e.Err
}

// NewServiceError creates a new service error
func NewServiceError(code string, message string, err error) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewServiceErrorWithCode creates a new service error with just a code
func NewServiceErrorWithCode(code string) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: code,
	}
}

// IsServiceError checks if an error is a ServiceError and returns its code
func IsServiceError(err error) (string, bool) {
	var serviceErr *ServiceError
	if errors.As(err, &serviceErr) {
		return serviceErr.Code, true
	}
	return "", false
}