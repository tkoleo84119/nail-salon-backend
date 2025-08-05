package adminTimeSlotTemplate

// UpdateRequest represents the request to update a time slot template
type UpdateRequest struct {
	Name *string `json:"name,omitempty" binding:"omitempty,min=1,max=50"`
	Note *string `json:"note,omitempty" binding:"omitempty,max=100"`
}

// UpdateResponse represents the response after updating a time slot template
type UpdateResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Note      string `json:"note"`
	Updater   string `json:"updater"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (r *UpdateRequest) HasUpdate() bool {
	return r.Name != nil || r.Note != nil
}
