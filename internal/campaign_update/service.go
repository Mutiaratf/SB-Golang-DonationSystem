package campaign_update

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetByCampaignID(campaignID int) ([]*CampaignUpdate, error) {
	return s.repository.GetByCampaignID(campaignID)
}

func (s *Service) Create(campaignID int, request *CampaignUpdateRequest) (*CampaignUpdate, error) {
	if err := s.repository.ValidateCampaign(campaignID); err != nil {
		return nil, err
	}
	update := &CampaignUpdate{CampaignID: campaignID, Title: request.Title, Content: request.Content}
	if err := s.repository.Create(update); err != nil {
		return nil, err
	}
	return update, nil
}

func (s *Service) Update(updateID, campaignID int, request *CampaignUpdateRequest) (*CampaignUpdate, error) {
	if err := s.repository.ValidateCampaign(campaignID); err != nil {
		return nil, err
	}
	update := &CampaignUpdate{Title: request.Title, Content: request.Content}
	if err := s.repository.Update(updateID, campaignID, update); err != nil {
		return nil, err
	}
	return update, nil
}

func (s *Service) Delete(updateID, campaignID int) error {
	if err := s.repository.ValidateCampaign(campaignID); err != nil {
		return err
	}
	return s.repository.Delete(updateID, campaignID)
}
