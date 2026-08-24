package transaction

import "time"

type Transaction struct {
	ID            int64     `json:"id"`
	DonorID       int64     `json:"donor_id"`
	Donor         string    `json:"donor"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Gender        string    `json:"gender"`
	IsAnonymous   bool      `json:"is_anonymous"`
	CampaignID    int64     `json:"campaign_id"`
	Amount        int64     `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	Prayer        string    `json:"prayer"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type TransactionHistory struct {
	Donor      string `json:"donor"`
	Prayer     string `json:"prayer"`
	Amount     int64  `json:"amount"`
	CampaignID int64  `json:"campaign_id"`
}
