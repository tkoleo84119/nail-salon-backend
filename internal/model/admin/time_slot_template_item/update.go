package adminTimeSlotTemplateItem

type UpdateRequest struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

type UpdateResponse struct {
	ID         string `json:"id"`
	TemplateID string `json:"templateId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}
