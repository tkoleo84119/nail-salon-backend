package adminStore

type UpdateRequest struct {
	Name     *string `json:"name" binding:"omitempty,noBlank,max=100"`
	Address  *string `json:"address" binding:"omitempty,noBlank,max=255"`
	Phone    *string `json:"phone" binding:"omitempty,taiwanphone"`
	IsActive *bool   `json:"isActive" binding:"omitempty"`
}

type UpdateResponse struct {
	ID string `json:"id"`
}

func (r UpdateRequest) HasUpdates() bool {
	return r.Name != nil || r.Address != nil || r.Phone != nil || r.IsActive != nil
}
