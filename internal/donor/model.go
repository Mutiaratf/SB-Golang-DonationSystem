package donor

import "time"

type Donor struct {
	ID        int64     `json:"id"`
	Donor     string    `json:"donor"`
	Gender    string    `json:"gender"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
