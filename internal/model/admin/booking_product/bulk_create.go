package adminBookingProduct

type BulkCreateRequest struct {
	ProductIds []string `json:"productIds" binding:"required,min=1,max=50"`
}

type BulkCreateParsedRequest struct {
	ProductIds []int64
}

type BulkCreateResponse struct {
	Created []string `json:"created"`
}
