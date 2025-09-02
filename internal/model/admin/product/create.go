package adminProduct

type CreateRequest struct {
	Name            string  `json:"name" binding:"required,noBlank,max=200"`
	BrandID         string  `json:"brandId" binding:"required"`
	CategoryID      string  `json:"categoryId" binding:"required"`
	CurrentStock    *int64  `json:"currentStock" binding:"required,min=0,max=1000000"`
	SafetyStock     *int64  `json:"safetyStock" binding:"omitempty,min=-1,max=1000000"`
	Unit            *string `json:"unit" binding:"omitempty,max=50"`
	StorageLocation *string `json:"storageLocation" binding:"omitempty,max=100"`
	Note            *string `json:"note" binding:"omitempty,max=255"`
}

type CreateParsedRequest struct {
	Name            string
	BrandID         int64
	CategoryID      int64
	CurrentStock    int
	SafetyStock     int
	Unit            *string
	StorageLocation *string
	Note            *string
}

type CreateResponse struct {
	ID string `json:"id"`
}
