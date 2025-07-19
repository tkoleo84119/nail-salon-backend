package timeSlotTemplate

// TimeSlotItem represents a time slot in a template
type TimeSlotItem struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

// CreateTimeSlotTemplateRequest represents the request to create a time slot template
type CreateTimeSlotTemplateRequest struct {
	Name      string         `json:"name" binding:"required,min=1,max=50"`
	Note      string         `json:"note" binding:"max=100"`
	TimeSlots []TimeSlotItem `json:"timeSlots" binding:"required,min=1,max=50"`
}

// TimeSlotItemResponse represents a time slot in the response
type TimeSlotItemResponse struct {
	ID        string `json:"id"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// CreateTimeSlotTemplateResponse represents the response after creating a time slot template
type CreateTimeSlotTemplateResponse struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Note      string                 `json:"note"`
	TimeSlots []TimeSlotItemResponse `json:"timeSlots"`
}

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
