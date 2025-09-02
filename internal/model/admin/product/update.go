package adminProduct

type UpdateRequest struct {
	BrandID         *string `json:"brandId" binding:"omitempty"`
	CategoryID      *string `json:"categoryId" binding:"omitempty"`
	Name            *string `json:"name" binding:"omitempty,noBlank,max=200"`
	CurrentStock    *int64  `json:"currentStock" binding:"omitempty,min=0,max=1000000"`
	SafetyStock     *int64  `json:"safetyStock" binding:"omitempty,min=-1,max=1000000"`
	Unit            *string `json:"unit" binding:"omitempty,max=50"`
	StorageLocation *string `json:"storageLocation" binding:"omitempty,max=100"`
	Note            *string `json:"note" binding:"omitempty,max=255"`
}

type UpdateParsedRequest struct {
	BrandID         *int64
	CategoryID      *int64
	Name            *string
	CurrentStock    *int
	SafetyStock     *int
	Unit            *string
	StorageLocation *string
	Note            *string
}

type UpdateResponse struct {
	ID string `json:"id"`
}

func (r *UpdateRequest) HasUpdates() bool {
	return r.BrandID != nil || r.CategoryID != nil || r.Name != nil || r.CurrentStock != nil || r.SafetyStock != nil || r.Unit != nil || r.StorageLocation != nil || r.Note != nil
}
