package common

import "github.com/golang-jwt/jwt/v5"

type ApiResponse struct {
	Message string            `json:"message,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

type Store struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type StaffContext struct {
	UserID    string  `json:"userId"`
	Username  string  `json:"username"`
	Role      string  `json:"role"`
	StoreList []Store `json:"storeList"`
}

type JWTClaims struct {
	StaffContext
	jwt.RegisteredClaims
}

type CustomerContext struct {
	CustomerID int64 `json:"customerId"`
}

type LineJWTClaims struct {
	jwt.RegisteredClaims
	CustomerID int64 `json:"customerId"`
}

func SuccessResponse(data interface{}) ApiResponse {
	return ApiResponse{
		Data: data,
	}
}

