package notification

import (
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Mutiaratf/SB-Golang-DonationSystem/internal/pdf"
	"gopkg.in/gomail.v2"
)

type Service struct{ repository *Repository }

func NewService(repository *Repository) *Service { return &Service{repository: repository} }

func (s *Service) GetRecipientEmail(id int64) (string, error) {
	transaction, err := s.repository.GetTransaction(id)
	if err != nil {
		return "", err
	}
	return transaction.Email, nil
}

func (s *Service) SendAsync(id int64) {
	go func() {
		if err := s.Send(id); err != nil {
			log.Printf("notification email failed for transaction %d: %v", id, err)
			return
		}
		log.Printf("notification email sent for transaction %d", id)
	}()
}

func (s *Service) Send(id int64) error {
	transaction, err := s.repository.GetTransaction(id)
	if err != nil {
		return err
	}
	template, err := s.repository.GetEmailTemplate()
	if err != nil {
		return err
	}
	smtpConfig, err := s.smtpConfig()
	if err != nil {
		return err
	}
	receiptBytes, err := pdf.GenerateReceipt(&pdf.Receipt{
		ID: transaction.ID, Donor: transaction.Donor, Phone: transaction.Phone,
		Campaign: transaction.Campaign, Amount: transaction.Amount,
		PaymentMethod: transaction.PaymentMethod, CreatedAt: transaction.CreatedAt,
		CompanyName: transaction.CompanyName, CompanyLogo: transaction.CompanyLogo,
		CompanyAddress: transaction.CompanyAddress, CompanyPhone: transaction.CompanyPhone,
		CompanyEmail: transaction.CompanyEmail, CompanyWebsite: transaction.CompanyWebsite,
		Director: transaction.Director, Sign: transaction.Sign,
	}, time.Now())
	if err != nil {
		return fmt.Errorf("generate receipt PDF: %w", err)
	}

	message := gomail.NewMessage()
	message.SetHeader("From", message.FormatAddress(smtpConfig.SenderEmail, smtpConfig.SenderName))
	message.SetHeader("To", transaction.Email)
	message.SetHeader("Subject", interpolate(template.Subject, transaction))
	message.SetBody("text/html", interpolateHTML(template.Body, transaction))
	message.Attach(fmt.Sprintf("Bukti_Donasi_%d.pdf", id), gomail.SetCopyFunc(func(writer io.Writer) error {
		_, err := writer.Write(receiptBytes)
		return err
	}))

	dialer := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, smtpConfig.Password)
	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil
}

func (s *Service) smtpConfig() (*SMTPConfig, error) {
	config, err := s.repository.GetSMTPConfig()
	if err == nil {
		return config, nil
	}
	if !errors.Is(err, ErrSMTPConfigNotFound) {
		return nil, err
	}
	port, parseErr := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if parseErr != nil {
		return nil, ErrSMTPConfigNotFound
	}
	config = &SMTPConfig{Host: os.Getenv("SMTP_HOST"), Port: port, Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"), SenderName: os.Getenv("SMTP_SENDER_NAME"), SenderEmail: os.Getenv("SMTP_SENDER_EMAIL")}
	if config.Host == "" || config.Username == "" || config.Password == "" || config.SenderEmail == "" {
		return nil, ErrSMTPConfigNotFound
	}
	return config, nil
}

func interpolate(value string, transaction *TransactionData) string {
	return strings.NewReplacer(
		"@donatur", transaction.Donor,
		"@program", transaction.Campaign,
		"@id_transaksi", pdf.TransactionNumber(transaction.ID, transaction.CreatedAt),
		"@tgldonasi", transaction.CreatedAt.Format("02-01-2006 15:04"),
		"@nominal", formatRupiah(transaction.Amount),
	).Replace(value)
}

func interpolateHTML(value string, transaction *TransactionData) string {
	return strings.NewReplacer(
		"@donatur", html.EscapeString(transaction.Donor),
		"@program", html.EscapeString(transaction.Campaign),
		"@id_transaksi", html.EscapeString(pdf.TransactionNumber(transaction.ID, transaction.CreatedAt)),
		"@tgldonasi", html.EscapeString(transaction.CreatedAt.Format("02-01-2006 15:04")),
		"@nominal", html.EscapeString(formatRupiah(transaction.Amount)),
	).Replace(value)
}

func formatRupiah(amount int64) string {
	value := strconv.FormatInt(amount, 10)
	for index := len(value) - 3; index > 0; index -= 3 {
		value = value[:index] + "." + value[index:]
	}
	return "Rp " + value
}
