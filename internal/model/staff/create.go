package staff

import "github.com/tkoleo84119/nail-salon-backend/internal/model/common"

// CreateStaffRequest represents the request to create a new staff member
type CreateStaffRequest struct {
	Username string   `json:"username" binding:"required,min=1,max=30"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=1,max=50"`
	Role     string   `json:"role" binding:"required,oneof=ADMIN MANAGER STYLIST"`
	StoreIDs []string `json:"storeIds" binding:"required,min=1,max=100"`
}

// CreateStaffResponse represents the response after creating a staff member
type CreateStaffResponse struct {
	ID        string         `json:"id"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	Role      string         `json:"role"`
	StoreList []common.Store `json:"storeList"`
}
