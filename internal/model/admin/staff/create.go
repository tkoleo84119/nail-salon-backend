package adminStaff

type CreateRequest struct {
	Username string   `json:"username" binding:"required,max=50"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,max=50"`
	Role     string   `json:"role" binding:"required,oneof=ADMIN MANAGER STYLIST"`
	StoreIDs []string `json:"storeIds" binding:"required,min=1,max=10"`
}

type CreateParsedRequest struct {
	Username string
	Email    string
	Password string
	Role     string
	StoreIDs []int64
}

type CreateResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
