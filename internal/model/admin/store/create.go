package adminStore

type CreateRequest struct {
	Name    string  `json:"name" binding:"required,noBlank,max=100"`
	Address *string `json:"address" binding:"omitempty,max=255"`
	Phone   *string `json:"phone" binding:"omitempty,taiwanphone"`
}

type CreateResponse struct {
	ID string `json:"id"`
}
