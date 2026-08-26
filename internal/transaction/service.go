package transaction

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/campaign"
	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/donor"
	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/pdf"
)

var (
	ErrCampaignInactive   = errors.New("campaign is inactive or not found")
	ErrAmountTooLow       = errors.New("amount is below campaign minimum amount")
	ErrTransactionAmount  = errors.New("amount must be greater than zero")
	ErrDonorDataRequired  = errors.New("donor data is required for a new donor")
	ErrDonorPhoneInvalid  = errors.New("donor phone number is invalid")
	ErrDonorGenderInvalid = errors.New("donor gender must be P or L")
	ErrDonorEmailInvalid  = errors.New("donor email is invalid")
)

type TransactionRequest struct {
	Donor         string `json:"donor"`
	Gender        string `json:"gender"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	IsAnonymous   bool   `json:"is_anonymous"`
	CampaignID    int64  `json:"campaign_id"`
	Amount        int64  `json:"amount"`
	PaymentMethod string `json:"payment_method"`
	Prayer        string `json:"prayer"`
}

type Service struct {
	repository      *Repository
	donorService    *donor.Service
	campaignService *campaign.Service
}

func NewService(repository *Repository, donorService *donor.Service, campaignService *campaign.Service) *Service {
	return &Service{repository: repository, donorService: donorService, campaignService: campaignService}
}

func (s *Service) validate(request *TransactionRequest) error {
	if request.Amount <= 0 {
		return ErrTransactionAmount
	}
	item, err := s.campaignService.GetByID(int(request.CampaignID))
	if err != nil || !item.IsActive {
		return ErrCampaignInactive
	}
	if request.Amount < int64(item.MinAmount) {
		return ErrAmountTooLow
	}

	return nil
}

func (s *Service) resolveDonor(request *TransactionRequest) (*donor.Donor, error) {
	phone := strings.TrimSpace(request.Phone)
	email := strings.TrimSpace(request.Email)
	item, err := s.donorService.GetDonorByEmail(email)
	if err == nil {
		return item, nil
	}
	if !errors.Is(err, donor.ErrDonorNotFound) {
		return nil, err
	}

	if strings.TrimSpace(request.Donor) == "" || strings.TrimSpace(request.Email) == "" || strings.TrimSpace(request.Gender) == "" {
		return nil, ErrDonorDataRequired
	}
	if request.Gender != "P" && request.Gender != "L" {
		return nil, ErrDonorGenderInvalid
	}
	if len(strings.TrimSpace(request.Phone)) < 8 || len(strings.TrimSpace(request.Phone)) > 15 {
		return nil, ErrDonorPhoneInvalid
	}
	address, err := mail.ParseAddress(request.Email)
	if err != nil || address.Address != request.Email {
		return nil, ErrDonorEmailInvalid
	}
	return s.donorService.CreateDonor(&donor.DonorRequest{
		Donor: request.Donor, Gender: request.Gender, Email: request.Email, Phone: phone,
	})
}

func (s *Service) Create(request *TransactionRequest) (*Transaction, error) {
	if err := s.validate(request); err != nil {
		return nil, err
	}
	item, err := s.resolveDonor(request)
	if err != nil {
		return nil, err
	}
	result := &Transaction{
		DonorID: item.ID, Donor: item.Donor, Email: item.Email, Phone: item.Phone, Gender: item.Gender,
		IsAnonymous: request.IsAnonymous, CampaignID: request.CampaignID, Amount: request.Amount,
		PaymentMethod: request.PaymentMethod, Prayer: request.Prayer,
	}
	if err := s.repository.Create(result); err != nil {
		return nil, err
	}
	if result.IsAnonymous {
		result.Donor = "Hamba Allah"
	}
	return result, nil
}

func (s *Service) GetAll() ([]*Transaction, error) {
	return s.repository.GetAll()
}

func (s *Service) GetReceipt(id int64) (*pdf.Receipt, error) {
	return s.repository.GetReceipt(id)
}

func (s *Service) GetHistory(campaignID int64) ([]*TransactionHistory, error) {
	return s.repository.GetHistory(campaignID)
}

func (s *Service) Update(id int64, request *TransactionRequest) (*Transaction, error) {
	if err := s.validate(request); err != nil {
		return nil, err
	}
	result := &Transaction{
		ID: id, IsAnonymous: request.IsAnonymous, CampaignID: request.CampaignID,
		Amount: request.Amount, PaymentMethod: request.PaymentMethod, Prayer: request.Prayer,
	}
	if err := s.repository.Update(id, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) Delete(id int64) error {
	return s.repository.Delete(id)
}
