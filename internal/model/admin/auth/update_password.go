package adminAuth

type UpdatePasswordRequest struct {
	StaffId     string  `json:"staffId" binding:"required"`
	OldPassword *string `json:"oldPassword" binding:"omitempty,noBlank,max=100"`
	NewPassword string  `json:"newPassword" binding:"required,noBlank,max=100"`
}

type UpdatePasswordParsedRequest struct {
	StaffId     int64
	OldPassword *string
	NewPassword string
}

type UpdatePasswordResponse struct {
	ID string `json:"id"`
}
