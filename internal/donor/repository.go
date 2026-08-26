package donor

import (
	"database/sql"
	"errors"
)

var ErrDonorNotFound = errors.New("donor not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(donor *Donor) error {
	return r.db.QueryRow(`INSERT INTO donors (donor, gender, email, phone)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`, donor.Donor, donor.Gender,
		donor.Email, donor.Phone).Scan(&donor.ID, &donor.CreatedAt, &donor.UpdatedAt)
}

func (r *Repository) ExistsByPhone(phone string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM donors WHERE phone = $1)`, phone).Scan(&exists)
	return exists, err
}

func (r *Repository) GetByPhone(phone string) (*Donor, error) {
	donor := new(Donor)
	err := r.db.QueryRow(`SELECT id, donor, gender, email, phone, created_at, updated_at
		FROM donors WHERE phone = $1`, phone).Scan(&donor.ID, &donor.Donor, &donor.Gender,
		&donor.Email, &donor.Phone, &donor.CreatedAt, &donor.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrDonorNotFound
	}
	return donor, err
}

func (r *Repository) GetByEmail(email string) (*Donor, error) {
	donor := new(Donor)
	err := r.db.QueryRow(`SELECT id, donor, gender, email, phone, created_at, updated_at
		FROM donors WHERE email = $1`, email).Scan(&donor.ID, &donor.Donor, &donor.Gender,
		&donor.Email, &donor.Phone, &donor.CreatedAt, &donor.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrDonorNotFound
	}
	return donor, err
}

func (r *Repository) GetAll() ([]*Donor, error) {
	rows, err := r.db.Query(`SELECT id, donor, gender, email, phone, created_at, updated_at
		FROM donors ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var donors []*Donor
	for rows.Next() {
		donor := new(Donor)
		if err := rows.Scan(&donor.ID, &donor.Donor, &donor.Gender, &donor.Email,
			&donor.Phone, &donor.CreatedAt, &donor.UpdatedAt); err != nil {
			return nil, err
		}
		donors = append(donors, donor)
	}
	return donors, rows.Err()
}

func (r *Repository) GetByID(id int64) (*Donor, error) {
	donor := new(Donor)
	err := r.db.QueryRow(`SELECT id, donor, gender, email, phone, created_at, updated_at
		FROM donors WHERE id = $1`, id).Scan(&donor.ID, &donor.Donor, &donor.Gender,
		&donor.Email, &donor.Phone, &donor.CreatedAt, &donor.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrDonorNotFound
	}
	return donor, err
}

func (r *Repository) Update(id int64, donor *Donor) error {
	result, err := r.db.Exec(`UPDATE donors SET donor = $1, gender = $2, email = $3,
		phone = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5`, donor.Donor,
		donor.Gender, donor.Email, donor.Phone, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrDonorNotFound
	}
	return nil
}

func (r *Repository) Delete(id int64) error {
	result, err := r.db.Exec(`DELETE FROM donors WHERE id = $1`, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrDonorNotFound
	}
	return nil
}
