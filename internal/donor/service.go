package donor

import (
	"errors"
	"net/mail"
	"strings"
	"time"
)

var (
	ErrDonorNameEmpty     = errors.New("donor name is empty")
	ErrDonorEmailEmpty    = errors.New("donor email is empty")
	ErrDonorEmailInvalid  = errors.New("donor email is invalid")
	ErrDonorPhoneEmpty    = errors.New("donor phone is empty")
	ErrDonorPhoneInvalid  = errors.New("donor phone is invalid")
	ErrDonorGenderInvalid = errors.New("gender must be P or L")
	ErrDonorPhoneExists   = errors.New("donor with this phone number already exists")
)

type DonorRequest struct {
	Donor  string `json:"donor"`
	Gender string `json:"gender"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
}

type Service struct {
	donorRepo *Repository
}

func NewService(donorRepo *Repository) *Service {
	return &Service{donorRepo: donorRepo}
}

func validateDonorRequest(req *DonorRequest) error {
	if strings.TrimSpace(req.Donor) == "" {
		return ErrDonorNameEmpty
	}
	if strings.TrimSpace(req.Email) == "" {
		return ErrDonorEmailEmpty
	}
	address, err := mail.ParseAddress(req.Email)
	if err != nil || address.Address != req.Email {
		return ErrDonorEmailInvalid
	}
	if strings.TrimSpace(req.Phone) == "" {
		return ErrDonorPhoneEmpty
	}
	if len(strings.TrimSpace(req.Phone)) < 8 || len(strings.TrimSpace(req.Phone)) > 15 {
		return ErrDonorPhoneInvalid
	}
	return nil
}

func donorFromRequest(req *DonorRequest) *Donor {
	return &Donor{
		Donor:     strings.TrimSpace(req.Donor),
		Gender:    strings.TrimSpace(req.Gender),
		Email:     strings.TrimSpace(req.Email),
		Phone:     strings.TrimSpace(req.Phone),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (s *Service) CreateDonor(req *DonorRequest) (*Donor, error) {
	if err := validateDonorRequest(req); err != nil {
		return nil, err
	}
	if req.Gender != "P" && req.Gender != "L" {
		return nil, ErrDonorGenderInvalid
	}
	phone := strings.TrimSpace(req.Phone)
	exists, err := s.donorRepo.ExistsByPhone(phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDonorPhoneExists
	}
	donor := donorFromRequest(req)
	if err := s.donorRepo.Create(donor); err != nil {
		return nil, err
	}
	return donor, nil
}

func (s *Service) GetAllDonors() ([]*Donor, error) {
	return s.donorRepo.GetAll()
}

func (s *Service) GetDonorByID(id int64) (*Donor, error) {
	return s.donorRepo.GetByID(id)
}

func (s *Service) GetDonorByPhone(phone string) (*Donor, error) {
	return s.donorRepo.GetByPhone(phone)
}

func (s *Service) UpdateDonor(id int64, req *DonorRequest) (*Donor, error) {
	if err := validateDonorRequest(req); err != nil {
		return nil, err
	}
	donor := donorFromRequest(req)
	donor.ID = id
	if err := s.donorRepo.Update(id, donor); err != nil {
		return nil, err
	}
	return donor, nil
}

func (s *Service) DeleteDonor(id int64) error {
	return s.donorRepo.Delete(id)
}
