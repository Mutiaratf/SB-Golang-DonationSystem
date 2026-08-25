package notification

import "time"

type TransactionData struct {
	ID             int64
	Donor          string
	Email          string
	Phone          string
	Campaign       string
	Amount         int64
	PaymentMethod  string
	CreatedAt      time.Time
	CompanyName    string
	CompanyLogo    string
	CompanyAddress string
	CompanyPhone   string
	CompanyEmail   string
	CompanyWebsite string
	Director       string
	Sign           string
}

type Template struct {
	Subject string
	Body    string
}

type SMTPConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	SenderName  string
	SenderEmail string
}
