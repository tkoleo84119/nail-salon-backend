package common

type ApiResponse struct {
	Message string            `json:"message,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func SuccessResponse(data interface{}) ApiResponse {
	return ApiResponse{
		Data: data,
	}
}

func ErrorResponse(message string, errors map[string]string) ApiResponse {
	return ApiResponse{
		Message: message,
		Errors:  errors,
	}
}

func ValidationErrorResponse(errors map[string]string) ApiResponse {
	return ApiResponse{
		Message: "輸入驗證失敗",
		Errors:  errors,
	}
}