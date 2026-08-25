package user

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type UserRepository interface {
	GetByEmail(email string) (*User, error)
	Create(user *User) error
}

type Service struct {
	repository UserRepository
}

func NewService(repository UserRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Authenticate(email, password string) (*User, error) {
	user, err := s.repository.GetByEmail(strings.TrimSpace(strings.ToLower(email)))
	if err != nil || user == nil || !user.IsActive {
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}

func (s *Service) Register(email, password string) (*User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:        email,
		PasswordHash: string(passwordHash),
		IsActive:     true,
	}
	if err := s.repository.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}
