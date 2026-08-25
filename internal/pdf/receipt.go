package pdf

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type Receipt struct {
	ID             int64
	Donor          string
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

const (
	directorTitle = "Direktur Yayasan Peduli Negeri"
	logoPath      = "asset/img/logo.png"
	signPath      = "asset/img/sign.png"
)

func GenerateReceipt(receipt *Receipt, printedAt time.Time) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(16, 14, 16)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()

	green := []int{31, 111, 88}
	pdf.SetDrawColor(green[0], green[1], green[2])
	pdf.SetLineWidth(0.8)
	pdf.Line(16, 39, 194, 39)

	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(10, 18, 24)
	pdf.Text(43, 22, fallback(receipt.CompanyName, "Yayasan Peduli Negeri"))
	pdf.SetFont("Arial", "", 9)
	pdf.SetFont("Arial", "", 7)
	pdf.Text(43, 26, "Alamat: "+receipt.CompanyAddress)
	pdf.Text(43, 29, "Phone: "+receipt.CompanyPhone)
	pdf.Text(43, 32, "Email: "+receipt.CompanyEmail)
	pdf.Text(43, 35, "Website: "+receipt.CompanyWebsite)
	addImage(pdf, logoPath, 17, 15, 26, 20)

	pdf.SetFont("Arial", "B", 24)
	pdf.SetXY(16, 51)
	pdf.CellFormat(178, 10, "BUKTI DONASI", "", 0, "C", false, 0, "")

	number := TransactionNumber(receipt.ID, receipt.CreatedAt)
	pdf.SetFillColor(236, 239, 241)
	pdf.RoundedRect(57, 66, 96, 11, 5, "1234", "F")
	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(57, 69)
	pdf.CellFormat(96, 5, "No. Transaksi: "+number, "", 0, "C", false, 0, "")

	sectionTitle(pdf, "Informasi Donatur", 84, green)
	row(pdf, "Nama", receipt.Donor, 96)
	row(pdf, "No. HP", receipt.Phone, 106)

	sectionTitle(pdf, "Informasi Donasi", 135, green)
	row(pdf, "Campaign", receipt.Campaign, 147)
	pdf.SetFont("Arial", "B", 9)
	row(pdf, "Nominal Donasi", formatRupiah(receipt.Amount), 157)
	row(pdf, "Metode Pembayaran", receipt.PaymentMethod, 167)
	row(pdf, "Tanggal Donasi", formatDateTime(receipt.CreatedAt), 177)

	pdf.SetFillColor(236, 239, 241)
	pdf.RoundedRect(16, 190, 178, 29, 4, "1234", "F")
	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(16, 199)
	pdf.CellFormat(178, 5, "Terima kasih atas kepedulian dan kepercayaannya.", "", 0, "C", false, 0, "")
	pdf.SetFont("Arial", "", 8.5)
	pdf.SetXY(30, 205)
	pdf.MultiCell(150, 4, "Donasi yang Anda berikan akan kami salurkan untuk tujuan sosial dan kemanusiaan\ndan bukan merupakan transaksi pencucian uang ataupun tindakan ilegal lainnya.", "", "C", false)

	pdf.SetFont("Arial", "", 9)
	pdf.Text(148, 231, "Bandung, "+formatDate(printedAt))
	pdf.Text(143, 238, directorTitle)
	if !addImage(pdf, signPath, 143, 239, 40, 22) {
		pdf.SetFont("Arial", "I", 16)
		pdf.Text(157, 253, receipt.Director)
	}
	pdf.SetDrawColor(30, 30, 30)
	pdf.Line(143, 261, 194, 261)
	pdf.SetFont("Arial", "B", 8)
	pdf.SetXY(143, 266)
	pdf.CellFormat(51, 4, fallback(receipt.Director, "Direktur"), "", 0, "C", false, 0, "")

	drawFooterWave(pdf, printedAt)

	var output bytes.Buffer
	if err := pdf.Output(&output); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func drawFooterWave(pdf *gofpdf.Fpdf, printedAt time.Time) {
	pdf.SetFillColor(132, 159, 137)
	points := []gofpdf.PointType{
		{0, 279.9}, {11.7, 279.9}, {23.3, 279.9}, {35, 280.2},
		{46.7, 281.1}, {58.3, 282.2}, {70, 283.4}, {81.7, 284.8},
		{93.3, 286.0}, {105, 287.0}, {116.7, 287.9}, {128.3, 288.4},
		{140, 288.4}, {151.7, 288.4}, {163.3, 288.4}, {175, 288.4},
		{186.7, 288.4}, {198.3, 288.4}, {210, 288.4},
		{210, 297}, {0, 297},
	}
	pdf.Polygon(points, "F")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "", 8)
	pdf.Text(17, 290, "Tanggal Cetak: "+formatDate(printedAt))
}

func sectionTitle(pdf *gofpdf.Fpdf, title string, y float64, green []int) {
	pdf.SetTextColor(10, 18, 24)
	pdf.SetFont("Arial", "B", 14)
	pdf.Text(17, y, title)
	pdf.SetDrawColor(green[0], green[1], green[2])
	pdf.SetLineWidth(1.1)
	pdf.Line(17, y+4, 34, y+4)
}

func row(pdf *gofpdf.Fpdf, label, value string, y float64) {
	pdf.SetTextColor(10, 18, 24)
	pdf.SetFont("Arial", "B", 9)
	pdf.Text(18, y, label)
	pdf.SetFont("Arial", "", 9)
	pdf.Text(66, y, ":")
	pdf.SetFont("Arial", "", 9)
	pdf.Text(71, y, value)
	pdf.SetDrawColor(110, 120, 125)
	pdf.SetLineWidth(0.25)
	pdf.Line(16, y+3, 194, y+3)
}

func TransactionNumber(id int64, date time.Time) string {
	return fmt.Sprintf("TRX-%s-%05d", date.Format("20060102"), id)
}

func formatRupiah(amount int64) string {
	value := fmt.Sprintf("%d", amount)
	for index := len(value) - 3; index > 0; index -= 3 {
		value = value[:index] + "." + value[index:]
	}
	return "Rp " + value
}

func formatDate(value time.Time) string     { return value.Format("02 January 2006") }
func formatDateTime(value time.Time) string { return value.Format("02 January 2006, 15:04 WIB") }
func fallback(value, defaultValue string) string {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}

func addImage(pdf *gofpdf.Fpdf, source string, x, y, width, height float64) bool {
	data, imageType, err := readImage(source)
	if err != nil {
		return false
	}
	name := fmt.Sprintf("image-%d-%d", int(x), int(y))
	pdf.RegisterImageOptionsReader(name, gofpdf.ImageOptions{ImageType: imageType}, bytes.NewReader(data))
	if pdf.Err() {
		return false
	}
	pdf.ImageOptions(name, x, y, width, height, false, gofpdf.ImageOptions{ImageType: imageType}, 0, "")
	return !pdf.Err()
}

func readImage(source string) ([]byte, string, error) {
	var data []byte
	var err error
	parsed, parseErr := url.Parse(source)
	if parseErr == nil && (parsed.Scheme == "http" || parsed.Scheme == "https") {
		client := &http.Client{Timeout: 5 * time.Second}
		response, requestErr := client.Get(source)
		if requestErr != nil {
			return nil, "", requestErr
		}
		defer response.Body.Close()
		if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
			return nil, "", fmt.Errorf("image request returned %s", response.Status)
		}
		data, err = io.ReadAll(io.LimitReader(response.Body, 5<<20))
	} else {
		data, err = os.ReadFile(source)
	}
	if err != nil {
		return nil, "", err
	}
	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || config.Width == 0 || config.Height == 0 {
		return nil, "", fmt.Errorf("invalid image")
	}
	imageType := strings.ToLower(http.DetectContentType(data)[len("image/"):])
	if imageType == "jpeg" {
		return data, "JPG", nil
	}
	if imageType == "png" {
		return data, "PNG", nil
	}
	return nil, "", fmt.Errorf("unsupported image type")
}
