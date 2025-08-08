package adminStylist

type GetAllRequest struct {
	Name        *string `form:"name" binding:"omitempty,max=100"`
	IsIntrovert *bool   `form:"isIntrovert" binding:"omitempty,boolean"`
	IsActive    *bool   `form:"isActive" binding:"omitempty,boolean"`
	Limit       *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset      *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
	Sort        *string `form:"sort" binding:"omitempty"`
}

type GetAllParsedRequest struct {
	Name        *string
	IsIntrovert *bool
	IsActive    *bool
	Limit       int
	Offset      int
	Sort        []string
}

type GetAllResponse struct {
	Total int          `json:"total"`
	Items []GetAllItem `json:"items"`
}

type GetAllItem struct {
	ID           string   `json:"id"`
	StaffUserID  string   `json:"staffUserId"`
	Name         string   `json:"name"`
	GoodAtShapes []string `json:"goodAtShapes"`
	GoodAtColors []string `json:"goodAtColors"`
	GoodAtStyles []string `json:"goodAtStyles"`
	IsIntrovert  bool     `json:"isIntrovert"`
	IsActive     bool     `json:"isActive"`
}
