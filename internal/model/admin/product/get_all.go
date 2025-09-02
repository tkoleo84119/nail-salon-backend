package adminProduct

type GetAllRequest struct {
	BrandID             *string `form:"brandId" binding:"omitempty"`
	CategoryID          *string `form:"categoryId" binding:"omitempty"`
	Name                *string `form:"name" binding:"omitempty,noBlank,max=100"`
	LessThanSafetyStock *bool   `form:"lessThanSafetyStock" binding:"omitempty"`
	Limit               *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset              *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
	Sort                *string `form:"sort" binding:"omitempty"`
}

type GetAllParsedRequest struct {
	BrandID             *int64
	CategoryID          *int64
	Name                *string
	LessThanSafetyStock *bool
	Limit               int
	Offset              int
	Sort                []string
}

type GetAllResponse struct {
	Total int                 `json:"total"`
	Items []GetAllProductItem `json:"items"`
}

type GetAllProductItem struct {
	ID              string                    `json:"id"`
	Name            string                    `json:"name"`
	Brand           GetAllProductBrandItem    `json:"brand"`
	Category        GetAllProductCategoryItem `json:"category"`
	CurrentStock    int                       `json:"currentStock"`
	SafetyStock     int                       `json:"safetyStock"`
	Unit            string                    `json:"unit"`
	StorageLocation string                    `json:"storageLocation"`
	Note            string                    `json:"note"`
	IsActive        bool                      `json:"isActive"`
	CreatedAt       string                    `json:"createdAt"`
	UpdatedAt       string                    `json:"updatedAt"`
}

type GetAllProductBrandItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetAllProductCategoryItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
