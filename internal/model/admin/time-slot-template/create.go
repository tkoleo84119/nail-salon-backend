package adminTimeSlotTemplate

// TimeSlotItem represents a time slot in a template
type CreateTimeSlotItem struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

// CreateTimeSlotTemplateRequest represents the request to create a time slot template
type CreateRequest struct {
	Name      string               `json:"name" binding:"required,noBlank,max=50"`
	Note      *string              `json:"note" binding:"max=100"`
	TimeSlots []CreateTimeSlotItem `json:"timeSlots" binding:"required,min=1,max=50"`
}

// TimeSlotItemResponse represents a time slot in the response
type CreateTimeSlotItemResponse struct {
	ID        string `json:"id"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// CreateTimeSlotTemplateResponse represents the response after creating a time slot template
type CreateResponse struct {
	ID        string                       `json:"id"`
	Name      string                       `json:"name"`
	Note      string                       `json:"note"`
	TimeSlots []CreateTimeSlotItemResponse `json:"timeSlots"`
}
