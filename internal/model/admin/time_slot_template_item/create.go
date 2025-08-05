package adminTimeSlotTemplateItem

// CreateRequest represents the request to create a time slot template item
type CreateRequest struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

// CreateResponse represents the response after creating a time slot template item
type CreateResponse struct {
	ID         string `json:"id"`
	TemplateID string `json:"templateId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}
