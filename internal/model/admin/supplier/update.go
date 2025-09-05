package adminSupplier

type UpdateRequest struct {
	Name     *string `json:"name" binding:"omitempty,noBlank,max=100"`
	IsActive *bool   `json:"isActive" binding:"omitempty"`
}

type UpdateResponse struct {
	ID string `json:"id"`
}

func (r *UpdateRequest) HasUpdates() bool {
	return r.Name != nil || r.IsActive != nil
}
