package user

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var ErrUserNotFound = errors.New("user not found")
var ErrEmailAlreadyExists = errors.New("email already exists")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByEmail(email string) (*User, error) {
	user := &User{}
	err := r.db.QueryRow(`
		SELECT id, email, password_hash, is_active, created_at
		FROM users WHERE email = $1`, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.IsActive, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) Create(user *User) error {
	err := r.db.QueryRow(`
		INSERT INTO users (email, password_hash, is_active)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`, user.Email, user.PasswordHash, user.IsActive).
		Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}
