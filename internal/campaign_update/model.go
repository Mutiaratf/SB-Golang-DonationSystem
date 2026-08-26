package campaign_update

type CampaignUpdate struct {
	ID         int    `json:"id"`
	CampaignID int    `json:"campaign_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type CampaignUpdateRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}
