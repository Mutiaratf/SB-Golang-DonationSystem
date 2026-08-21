package campaign

type Campaign struct {
	ID           int    `json:"id"`
	CampaignName string `json:"campaign"`
	CategoryID   int    `json:"category_id"`
	Description  string `json:"description"`
	IsActive     bool   `json:"is_active"`
	MinAmount    int    `json:"min_amount"`
	TargetAmount int    `json:"target_amount"`
	Thumbnail    string `json:"thumbnail"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
