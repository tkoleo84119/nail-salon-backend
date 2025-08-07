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
	UserID    int64   `json:"userId"`
	Username  string  `json:"username"`
	Role      string  `json:"role"`
	StoreList []Store `json:"storeList"`
}

type StaffJWTClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

type CustomerContext struct {
	CustomerID int64 `json:"customerId"`
}

type LineJWTClaims struct {
	CustomerID int64 `json:"customerId"`
	jwt.RegisteredClaims
}

type ListResponse[T any] struct {
	Total int `json:"total"`
	Items []T `json:"items"`
}

func SuccessResponse(data interface{}) ApiResponse {
	return ApiResponse{
		Data: data,
	}
}
