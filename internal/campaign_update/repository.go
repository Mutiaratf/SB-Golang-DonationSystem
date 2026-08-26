package campaign_update

import (
	"database/sql"
	"errors"
)

var ErrCampaignNotFound = errors.New("campaign not found")
var ErrCampaignInactive = errors.New("campaign is inactive")
var ErrCampaignUpdateNotFound = errors.New("campaign update not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ValidateCampaign(campaignID int) error {
	var isActive bool
	err := r.db.QueryRow(`
		SELECT is_active
		FROM campaigns
		WHERE id = $1`, campaignID).Scan(&isActive)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrCampaignNotFound
	}
	if err != nil {
		return err
	}
	if !isActive {
		return ErrCampaignInactive
	}
	return nil
}

func (r *Repository) Create(update *CampaignUpdate) error {
	return r.db.QueryRow(`
		INSERT INTO campaign_updates (campaign_id, title, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`, update.CampaignID, update.Title, update.Content).
		Scan(&update.ID, &update.CreatedAt, &update.UpdatedAt)
}

func (r *Repository) Update(updateID, campaignID int, update *CampaignUpdate) error {
	result, err := r.db.Exec(`
		UPDATE campaign_updates
		SET title = $1, content = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND campaign_id = $4`, update.Title, update.Content, updateID, campaignID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrCampaignUpdateNotFound
	}
	update.ID = updateID
	update.CampaignID = campaignID
	return nil
}

func (r *Repository) Delete(updateID, campaignID int) error {
	result, err := r.db.Exec(`
		DELETE FROM campaign_updates
		WHERE id = $1 AND campaign_id = $2`, updateID, campaignID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrCampaignUpdateNotFound
	}
	return nil
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
