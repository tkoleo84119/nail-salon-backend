package adminBookingProduct

type BulkDeleteRequest struct {
	ProductIds []string `json:"productIds" binding:"required,min=1,max=50"`
}

type BulkDeleteParsedRequest struct {
	ProductIds []int64
}

type BulkDeleteResponse struct {
	Deleted []string `json:"deleted"`
}
