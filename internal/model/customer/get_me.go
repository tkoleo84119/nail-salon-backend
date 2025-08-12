package customer

type GetMeResponse struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Phone          string   `json:"phone"`
	Birthday       string   `json:"birthday"`
	Email          string   `json:"email"`
	City           string   `json:"city"`
	FavoriteShapes []string `json:"favoriteShapes"`
	FavoriteColors []string `json:"favoriteColors"`
	FavoriteStyles []string `json:"favoriteStyles"`
	IsIntrovert    bool     `json:"isIntrovert"`
	CustomerNote   string   `json:"customerNote"`
}
