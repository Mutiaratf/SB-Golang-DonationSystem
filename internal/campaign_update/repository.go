package campaign_update

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByCampaignID(campaignID int) ([]*CampaignUpdate, error) {
	rows, err := r.db.Query(`
		SELECT id, campaign_id, title, content, created_at, updated_at
		FROM campaign_updates
		WHERE campaign_id = $1
		ORDER BY id`, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	updates := make([]*CampaignUpdate, 0)
	for rows.Next() {
		update := new(CampaignUpdate)
		if err := rows.Scan(&update.ID, &update.CampaignID, &update.Title,
			&update.Content, &update.CreatedAt, &update.UpdatedAt); err != nil {
			return nil, err
		}
		updates = append(updates, update)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return updates, nil
}
