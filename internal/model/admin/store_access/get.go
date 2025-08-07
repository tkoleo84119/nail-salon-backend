package adminStoreAccess

type GetResponse struct {
	StoreList []GetStore `json:"storeList"`
}

type GetStore struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
