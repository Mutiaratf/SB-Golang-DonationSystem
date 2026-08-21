package campaign

import "errors"

var (
	ErrCampaignNameEmpty = errors.New("campaign name is empty")
	ErrCategoryIDInvalid = errors.New("category id is invalid")
	ErrAmountInvalid     = errors.New("amount is invalid")
)

type Service struct {
	repository *Repository
}

type CampaignRequest struct {
	CampaignName string `json:"campaign"`
	CategoryID   int    `json:"category_id"`
	Description  string `json:"description"`
	IsActive     bool   `json:"is_active"`
	MinAmount    int    `json:"min_amount"`
	TargetAmount int    `json:"target_amount"`
	Thumbnail    string `json:"thumbnail"`
}

func NewService(repository *Repository) *Service { return &Service{repository: repository} }

func validateRequest(request *CampaignRequest) error {
	if request.CampaignName == "" {
		return ErrCampaignNameEmpty
	}
	if request.CategoryID <= 0 {
		return ErrCategoryIDInvalid
	}
	if request.MinAmount < 0 || request.TargetAmount <= 0 || request.MinAmount > request.TargetAmount {
		return ErrAmountInvalid
	}
	return nil
}

func (s *Service) Create(request *CampaignRequest) (*Campaign, error) {
	if err := validateRequest(request); err != nil {
		return nil, err
	}
	campaign := &Campaign{CampaignName: request.CampaignName, CategoryID: request.CategoryID,
		Description: request.Description, IsActive: request.IsActive, MinAmount: request.MinAmount,
		TargetAmount: request.TargetAmount, Thumbnail: request.Thumbnail}
	if err := s.repository.Create(campaign); err != nil {
		return nil, err
	}
	return campaign, nil
}

func (s *Service) GetAll() ([]*Campaign, error) { return s.repository.GetAll() }

func (s *Service) GetByID(id int) (*Campaign, error) { return s.repository.GetByID(id) }

func (s *Service) Update(id int, request *CampaignRequest) (*Campaign, error) {
	if err := validateRequest(request); err != nil {
		return nil, err
	}
	campaign := &Campaign{ID: id, CampaignName: request.CampaignName, CategoryID: request.CategoryID,
		Description: request.Description, IsActive: request.IsActive, MinAmount: request.MinAmount,
		TargetAmount: request.TargetAmount, Thumbnail: request.Thumbnail}
	if err := s.repository.Update(id, campaign); err != nil {
		return nil, err
	}
	return campaign, nil
}

func (s *Service) Delete(id int) error { return s.repository.Delete(id) }
