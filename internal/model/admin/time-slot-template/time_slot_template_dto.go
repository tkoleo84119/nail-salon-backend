package adminTimeSlotTemplate

// UpdateTimeSlotTemplateRequest represents the request to update a time slot template
type UpdateTimeSlotTemplateRequest struct {
	Name *string `json:"name,omitempty" binding:"omitempty,min=1,max=50"`
	Note *string `json:"note,omitempty" binding:"omitempty,max=100"`
}

// UpdateTimeSlotTemplateResponse represents the response after updating a time slot template
type UpdateTimeSlotTemplateResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Note string `json:"note"`
}

func (r *UpdateTimeSlotTemplateRequest) HasUpdate() bool {
	return r.Name != nil || r.Note != nil
}
