package adminCustomer

type GetResponse struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	LineName       string   `json:"lineName"`
	Phone          string   `json:"phone"`
	Birthday       string   `json:"birthday"`
	Email          string   `json:"email"`
	City           string   `json:"city"`
	FavoriteShapes []string `json:"favoriteShapes"`
	FavoriteColors []string `json:"favoriteColors"`
	FavoriteStyles []string `json:"favoriteStyles"`
	IsIntrovert    bool     `json:"isIntrovert"`
	ReferralSource []string `json:"referralSource"`
	Referrer       string   `json:"referrer"`
	CustomerNote   string   `json:"customerNote"`
	StoreNote      string   `json:"storeNote"`
	Level          string   `json:"level"`
	IsBlacklisted  bool     `json:"isBlacklisted"`
	LastVisitAt    string   `json:"lastVisitAt"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
}
