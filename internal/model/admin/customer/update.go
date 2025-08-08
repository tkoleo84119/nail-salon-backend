package adminCustomer

type UpdateRequest struct {
	StoreNote     *string `json:"storeNote" binding:"omitempty,max=255"`
	Level         *string `json:"level" binding:"omitempty,oneof=NORMAL VIP VVIP"`
	IsBlacklisted *bool   `json:"isBlacklisted" binding:"omitempty,boolean"`
}

type UpdateResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Birthday      string `json:"birthday"`
	City          string `json:"city"`
	Level         string `json:"level"`
	IsBlacklisted bool   `json:"isBlacklisted"`
	LastVisitAt   string `json:"lastVisitAt,omitempty"`
	UpdatedAt     string `json:"updatedAt"`
}

func (r *UpdateRequest) HasUpdates() bool {
	return r.StoreNote != nil || r.Level != nil || r.IsBlacklisted != nil
}
