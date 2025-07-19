package timeSlotTemplate

// CreateTimeSlotTemplateItemRequest represents the request to create a time slot template item
type CreateTimeSlotTemplateItemRequest struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

// CreateTimeSlotTemplateItemResponse represents the response after creating a time slot template item
type CreateTimeSlotTemplateItemResponse struct {
	ID         string `json:"id"`
	TemplateID string `json:"templateId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

type UpdateTimeSlotTemplateItemRequest struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

type UpdateTimeSlotTemplateItemResponse struct {
	ID         string `json:"id"`
	TemplateID string `json:"templateId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

type DeleteTimeSlotTemplateItemResponse struct {
	Deleted []string `json:"deleted"`
}
