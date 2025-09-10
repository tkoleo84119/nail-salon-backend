package adminAuth

type PermissionResponse struct {
	ID          string      `json:"id"`
	Username    string      `json:"username"`
	Role        string      `json:"role"`
	StoreAccess []StoreInfo `json:"storeAccess"`
}

type StoreInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
