package transaction

import (
	"database/sql"
	"errors"

	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/pdf"
)

var ErrTransactionNotFound = errors.New("transaction not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Create(transaction *Transaction) error {
	return r.db.QueryRow(`INSERT INTO transactions
		(donor_id, email, phone, gender, campaign_id, amount, payment_method, prayer, is_anonymous)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`, transaction.DonorID, transaction.Email, transaction.Phone, transaction.Gender, transaction.CampaignID,
		transaction.Amount, transaction.PaymentMethod, transaction.Prayer,
		transaction.IsAnonymous).Scan(&transaction.ID, &transaction.CreatedAt, &transaction.UpdatedAt)
}

func (r *Repository) GetAll() ([]*Transaction, error) {
	rows, err := r.db.Query(`SELECT t.id, t.donor_id, d.donor, d.email, d.phone, d.gender,
		t.is_anonymous, t.campaign_id, t.amount, t.payment_method, t.prayer,
		t.created_at, t.updated_at FROM transactions t JOIN donors d ON d.id = t.donor_id ORDER BY t.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTransactions(rows)
}

func (r *Repository) GetReceipt(id int64) (*pdf.Receipt, error) {
	item := new(pdf.Receipt)
	var companyName, companyLogo, companyAddress, companyPhone, companyEmail, companyWebsite, director, sign sql.NullString
	row := r.db.QueryRow(`SELECT t.id, d.donor, d.phone,
		c.campaign, t.amount, t.payment_method, t.created_at,
		cp.company_name, cp.logo, cp.address, cp.phone, cp.email, cp.website,
		cp.director, cp.sign
		FROM transactions t
		JOIN donors d ON d.id = t.donor_id
		JOIN campaigns c ON c.id = t.campaign_id
		LEFT JOIN company_profile cp ON true
		WHERE t.id = $1`, id)
	if err := row.Scan(&item.ID, &item.Donor, &item.Phone,
		&item.Campaign, &item.Amount, &item.PaymentMethod, &item.CreatedAt,
		&companyName, &companyLogo, &companyAddress, &companyPhone, &companyEmail,
		&companyWebsite, &director, &sign); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}
	item.CompanyName = companyName.String
	item.CompanyLogo = companyLogo.String
	item.CompanyAddress = companyAddress.String
	item.CompanyPhone = companyPhone.String
	item.CompanyEmail = companyEmail.String
	item.CompanyWebsite = companyWebsite.String
	item.Director = director.String
	item.Sign = sign.String
	return item, nil
}

func (r *Repository) GetHistory(campaignID int64) ([]*TransactionHistory, error) {
	rows, err := r.db.Query(`SELECT
		CASE WHEN t.is_anonymous THEN 'Hamba Allah' ELSE d.donor END AS donor,
		t.prayer, t.amount, t.campaign_id
		FROM transactions t JOIN donors d ON d.id = t.donor_id
		WHERE t.campaign_id = $1 ORDER BY t.id`, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*TransactionHistory
	for rows.Next() {
		item := new(TransactionHistory)
		if err := rows.Scan(&item.Donor, &item.Prayer, &item.Amount, &item.CampaignID); err != nil {
			return nil, err
		}
		transactions = append(transactions, item)
	}
	return transactions, rows.Err()
}

func scanTransactions(rows *sql.Rows) ([]*Transaction, error) {
	var transactions []*Transaction
	for rows.Next() {
		item := new(Transaction)
		if err := rows.Scan(&item.Donor, &item.CampaignID, &item.Amount, &item.Prayer); err != nil {
			return nil, err
		}
		transactions = append(transactions, item)
	}
	return transactions, rows.Err()
}

func (r *Repository) Update(id int64, transaction *Transaction) error {
	result, err := r.db.Exec(`UPDATE transactions SET campaign_id = $1, amount = $2,
		payment_method = $3, prayer = $4, is_anonymous = $5,
		updated_at = CURRENT_TIMESTAMP WHERE id = $6`, transaction.CampaignID,
		transaction.Amount, transaction.PaymentMethod, transaction.Prayer,
		transaction.IsAnonymous, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrTransactionNotFound
	}
	return nil
}

func (r *Repository) Delete(id int64) error {
	result, err := r.db.Exec(`DELETE FROM transactions WHERE id = $1`, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrTransactionNotFound
	}
	return nil
}
