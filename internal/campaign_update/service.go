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
