package adminCoupon

type UpdateRequest struct {
	Name     *string `json:"name" binding:"omitempty,noBlank,max=100"`
	IsActive *bool   `json:"isActive" binding:"omitempty"`
	Note     *string `json:"note" binding:"omitempty,max=255"`
}

type UpdateResponse struct {
	ID string `json:"id"`
}

func (r UpdateRequest) HasUpdates() bool {
	return r.Name != nil || r.IsActive != nil || r.Note != nil
}
