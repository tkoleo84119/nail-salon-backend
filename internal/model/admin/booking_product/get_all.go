package adminBookingProduct

type GetAllRequest struct {
	Limit  *int    `form:"limit"`
	Offset *int    `form:"offset"`
	Sort   *string `form:"sort" binding:"omitempty"`
}

type GetAllParsedRequest struct {
	Limit  int
	Offset int
	Sort   []string
}

type GetAllResponse struct {
	Total int              `json:"total"`
	Items []GetAllItemData `json:"items"`
}

type GetAllItemData struct {
	ID        string                `json:"id"`
	Product   GetAllItemProductData `json:"product"`
	CreatedAt string                `json:"createdAt"`
}

type GetAllItemProductData struct {
	ID       string                        `json:"id"`
	Name     string                        `json:"name"`
	Brand    GetAllItemProductBrandData    `json:"brand"`
	Category GetAllItemProductCategoryData `json:"category"`
}

type GetAllItemProductBrandData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetAllItemProductCategoryData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
