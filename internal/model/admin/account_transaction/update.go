package adminAccountTransaction

type UpdateRequest struct {
	Note *string `json:"note" binding:"omitempty,max=255"`
}

type UpdateParsedRequest struct {
	Note *string
}

type UpdateResponse struct {
	ID string `json:"id"`
}

func (r *UpdateRequest) HasUpdates() bool {
	return r.Note != nil
}
