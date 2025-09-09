package auth

type AcceptTermRequest struct {
	TermsVersion string `json:"termsVersion" binding:"required,oneof=v1"`
}

type AcceptTermResponse struct {
	ID string `json:"id"`
}
