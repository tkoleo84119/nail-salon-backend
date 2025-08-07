package adminStoreAccess

type CreateRequest struct {
	StoreID string `json:"storeId" binding:"required"`
}

type CreateResponse struct {
	StoreList []CreateStore `json:"storeList"`
}

type CreateStore struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
