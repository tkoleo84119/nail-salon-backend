package common

import "github.com/golang-jwt/jwt/v5"

type ApiResponse struct {
	Message string            `json:"message,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

type Store struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type StaffContext struct {
	UserID    int64   `json:"user_id"`
	Username  string  `json:"username"`
	Role      string  `json:"role"`
	StoreList []Store `json:"store_list"`
}

type JWTClaims struct {
	StaffContext
	jwt.RegisteredClaims
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