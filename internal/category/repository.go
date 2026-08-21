package category

import (
	"database/sql"
	"errors"
)

var ErrCategoryNotFound = errors.New("category not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(category *Model) error {
	query := `INSERT INTO campaign_categories (category, is_active)
		VALUES ($1, $2) RETURNING id`

	err := r.db.QueryRow(
		query,
		category.Category,
		category.IsActive,
	).Scan(&category.ID)
	return err
}

func (r *Repository) GetAll() ([]*Model, error) {
	query := `SELECT * FROM campaign_categories order by id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Model
	for rows.Next() {
		var category Model
		if err := rows.Scan(&category.ID, &category.Category, &category.IsActive); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *Repository) Update(id int64, category *Model) error {
	query := `UPDATE campaign_categories SET category = $1, is_active = $2 WHERE id = $3`
	result, err := r.db.Exec(query, category.Category, category.IsActive, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrCategoryNotFound
	}
	return nil
}

func (r *Repository) Delete(id int64) error {
	query := `DELETE FROM campaign_categories WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrCategoryNotFound
	}
	return nil
}
