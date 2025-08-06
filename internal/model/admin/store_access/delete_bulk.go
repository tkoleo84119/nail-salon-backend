package adminStoreAccess

type DeleteBulkRequest struct {
	StoreIDs []string `json:"storeIds" binding:"required,min=1,max=20"`
}

type DeleteBulkResponse struct {
	Deleted []string `json:"deleted"`
}
