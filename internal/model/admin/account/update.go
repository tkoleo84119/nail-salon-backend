package adminAccount

type UpdateRequest struct {
	Name     *string `json:"name" binding:"omitempty,noBlank,max=100"`
	Note     *string `json:"note" binding:"omitempty,max=255"`
	IsActive *bool   `json:"isActive"`
}

type UpdateParsedRequest struct {
	Name     *string
	Note     *string
	IsActive *bool
}

type UpdateResponse struct {
	ID string `json:"id"`
}

func (r *UpdateRequest) HasUpdates() bool {
	return r.Name != nil || r.Note != nil || r.IsActive != nil
}