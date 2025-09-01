package adminBrand

type UpdateRequest struct {
	Name     *string `json:"name" binding:"omitempty,noBlank,max=100"`
	IsActive *bool   `json:"isActive"`
}

type UpdateResponse struct {
	ID string `json:"id"`
}

func (r UpdateRequest) HasUpdates() bool {
	return r.Name != nil || r.IsActive != nil
}
