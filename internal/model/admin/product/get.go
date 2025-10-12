package adminProduct

type GetResponse struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Brand           GetProductBrandItem    `json:"brand"`
	Category        GetProductCategoryItem `json:"category"`
	CurrentStock    int                    `json:"currentStock"`
	SafetyStock     int                    `json:"safetyStock"`
	Unit            string                 `json:"unit"`
	StorageLocation string                 `json:"storageLocation"`
	Note            string                 `json:"note"`
	IsActive        bool                   `json:"isActive"`
	CreatedAt       string                 `json:"createdAt"`
	UpdatedAt       string                 `json:"updatedAt"`
}

type GetProductBrandItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetProductCategoryItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
