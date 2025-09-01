package adminBrand

type CreateRequest struct {
	Name string `json:"name" binding:"required,noBlank,max=100"`
}

type CreateResponse struct {
	ID string `json:"id"`
}
