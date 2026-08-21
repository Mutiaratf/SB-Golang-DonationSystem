package campaign

import (
	"database/sql"
	"errors"
)

var ErrCampaignNotFound = errors.New("campaign not found")
var ErrCategoryInactive = errors.New("category not found or inactive")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(campaign *Campaign) error {
	query := `INSERT INTO campaigns
		(campaign, category_id, description, is_active, min_amount, target_amount, thumbnail)
		SELECT $1, id, $3, $4, $5, $6, $7
		FROM campaign_categories
		WHERE id = $2 AND is_active = true
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(query, campaign.CampaignName, campaign.CategoryID,
		campaign.Description, campaign.IsActive, campaign.MinAmount,
		campaign.TargetAmount, campaign.Thumbnail).Scan(&campaign.ID,
		&campaign.CreatedAt, &campaign.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrCategoryInactive
	}
	return err
}

func (r *Repository) GetAll() ([]*Campaign, error) {
	rows, err := r.db.Query(`SELECT id, campaign, category_id, description, is_active,
		min_amount, target_amount, thumbnail, created_at, updated_at
		FROM campaigns ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []*Campaign
	for rows.Next() {
		campaign := new(Campaign)
		if err := rows.Scan(&campaign.ID, &campaign.CampaignName, &campaign.CategoryID,
			&campaign.Description, &campaign.IsActive, &campaign.MinAmount,
			&campaign.TargetAmount, &campaign.Thumbnail, &campaign.CreatedAt,
			&campaign.UpdatedAt); err != nil {
			return nil, err
		}
		campaigns = append(campaigns, campaign)
	}
	return campaigns, rows.Err()
}

func (r *Repository) GetByID(id int) (*Campaign, error) {
	campaign := new(Campaign)
	err := r.db.QueryRow(`SELECT id, campaign, category_id, description, is_active,
		min_amount, target_amount, thumbnail, created_at, updated_at
		FROM campaigns WHERE id = $1`, id).Scan(
		&campaign.ID, &campaign.CampaignName, &campaign.CategoryID, &campaign.Description,
		&campaign.IsActive, &campaign.MinAmount, &campaign.TargetAmount,
		&campaign.Thumbnail, &campaign.CreatedAt, &campaign.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrCampaignNotFound
	}
	return campaign, err
}

func (r *Repository) Update(id int, campaign *Campaign) error {
	result, err := r.db.Exec(`UPDATE campaigns SET campaign = $1, category_id = $2,
		description = $3, is_active = $4, min_amount = $5, target_amount = $6,
		thumbnail = $7, updated_at = CURRENT_TIMESTAMP WHERE id = $8`,
		campaign.CampaignName, campaign.CategoryID, campaign.Description, campaign.IsActive,
		campaign.MinAmount, campaign.TargetAmount, campaign.Thumbnail, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrCampaignNotFound
	}
	return nil
}

func (r *Repository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM campaigns WHERE id = $1`, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrCampaignNotFound
	}
	return nil
}
