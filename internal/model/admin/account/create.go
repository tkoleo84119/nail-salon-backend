package adminAccount

type CreateRequest struct {
	StoreID string  `json:"storeId" binding:"required"`
	Name    string  `json:"name" binding:"required,noBlank,max=100"`
	Note    *string `json:"note" binding:"omitempty,max=255"`
}

type CreateParsedRequest struct {
	StoreID int64
	Name    string
	Note    *string
}

type CreateResponse struct {
	ID string `json:"id"`
}
