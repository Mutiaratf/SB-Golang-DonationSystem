package notification

import (
	"database/sql"
	"errors"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrTemplateNotFound    = errors.New("email template not found")
	ErrSMTPConfigNotFound  = errors.New("active SMTP configuration not found")
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) GetTransaction(id int64) (*TransactionData, error) {
	item := new(TransactionData)
	row := r.db.QueryRow(`SELECT t.id, d.donor, d.email, d.phone, c.campaign, t.amount,
		t.payment_method, t.created_at, cp.company_name, cp.logo, cp.address,
		cp.phone, cp.email, cp.website, cp.director, cp.sign
		FROM transactions t
		JOIN donors d ON d.id = t.donor_id
		JOIN campaigns c ON c.id = t.campaign_id
		LEFT JOIN company_profile cp ON true
		WHERE t.id = $1`, id)
	var companyName, companyLogo, companyAddress, companyPhone, companyEmail sql.NullString
	var companyWebsite, director, sign sql.NullString
	if err := row.Scan(&item.ID, &item.Donor, &item.Email, &item.Phone, &item.Campaign,
		&item.Amount, &item.PaymentMethod, &item.CreatedAt, &companyName, &companyLogo,
		&companyAddress, &companyPhone, &companyEmail, &companyWebsite, &director, &sign); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}
	item.CompanyName, item.CompanyLogo, item.CompanyAddress = companyName.String, companyLogo.String, companyAddress.String
	item.CompanyPhone, item.CompanyEmail, item.CompanyWebsite = companyPhone.String, companyEmail.String, companyWebsite.String
	item.Director, item.Sign = director.String, sign.String
	return item, nil
}

func (r *Repository) GetEmailTemplate() (*Template, error) {
	item := new(Template)
	err := r.db.QueryRow(`SELECT subject, body FROM templates
		WHERE type = 'email' ORDER BY id LIMIT 1`).Scan(&item.Subject, &item.Body)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTemplateNotFound
	}
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *Repository) GetSMTPConfig() (*SMTPConfig, error) {
	item := new(SMTPConfig)
	err := r.db.QueryRow(`SELECT host, port, username, password, sender_name, sender_email
		FROM smtp_configs WHERE is_active = true ORDER BY id LIMIT 1`).Scan(
		&item.Host, &item.Port, &item.Username, &item.Password, &item.SenderName, &item.SenderEmail)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrSMTPConfigNotFound
	}
	if err != nil {
		return nil, err
	}
	return item, nil
}
