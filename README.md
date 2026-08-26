# SB-Golang-DonationSystem

REST API untuk platform donasi digital. Aplikasi ini mengelola kategori, campaign, update campaign, donatur, transaksi donasi, notifikasi email, dan bukti donasi PDF.

## Daftar Isi

- [Menjalankan Aplikasi](#menjalankan-aplikasi)
- [Autentikasi](#autentikasi)
- [Format Response](#format-response)
- [Dokumentasi Endpoint](#dokumentasi-endpoint)
- [Notifikasi Email](#notifikasi-email)
- [Struktur Direktori](#struktur-direktori)

## Menjalankan Aplikasi

### Prasyarat

- Go sesuai versi pada `go.mod`.
- PostgreSQL yang dapat diakses aplikasi.
- File `.env` pada root project.

### Konfigurasi Environment

```dotenv
DBHOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=your-password
DB_NAME=donation_system
DB_SSL_MODE=disable
JWT_SECRET=your-secret-minimum-32-characters
JWT_ISSUER=donation-system
JWT_EXPIRY_HOURS=24
```

`JWT_SECRET` wajib memiliki panjang minimal 32 karakter. Jangan commit password database, JWT secret, atau kredensial SMTP.

### Database dan Server

Jalankan migration user:

```bash
psql "$DATABASE_URL" -f migrations/003_users.sql
```

Buat user administrator menggunakan password hash bcrypt, lalu jalankan server:

```bash
go run ./cmd/server
```

Server berjalan pada `http://localhost:8082`. Semua endpoint menggunakan prefix `/api`.

```text
Base URL: http://localhost:8082/api
Content-Type: application/json
```

## Autentikasi

Endpoint yang memerlukan JWT ditandai **Protected**. Setelah login, kirim token pada setiap request protected:

```http
Authorization: Bearer <jwt-token>
```

Token dibuat dengan algoritma HS256 dan memuat issuer, subject user, waktu terbit, waktu mulai berlaku, dan waktu kedaluwarsa.

### POST `/auth/register` - Public

Mendaftarkan akun user baru.

Request body:

```json
{
  "email": "admin@example.com",
  "password": "password123"
}
```

Aturan:

- Email wajib valid.
- Password minimal 8 karakter.
- Email yang sudah terdaftar menghasilkan `409 Conflict`.
- Password disimpan sebagai bcrypt hash.

Response `201 Created`:

```json
{
  "status": "success",
  "message": "Registration successful",
  "data": {
    "id": 1,
    "email": "admin@example.com",
    "is_active": true,
    "created_at": "2026-08-26T10:00:00Z"
  }
}
```

### POST `/auth/login` - Public

Memvalidasi user dan menghasilkan JWT.

Request body:

```json
{
  "email": "admin@example.com",
  "password": "password123"
}
```

Response `200 OK`:

```json
{
  "status": "success",
  "data": {
    "token": "<jwt-token>",
    "user": {
      "id": 1,
      "email": "admin@example.com",
      "is_active": true,
      "created_at": "2026-08-26T10:00:00Z"
    }
  }
}
```

Email atau password salah menghasilkan `401 Unauthorized`.

## Format Response

Response JSON sukses umumnya berbentuk:

```json
{
  "status": "success",
  "message": "Operation completed successfully",
  "data": {}
}
```

Response error:

```json
{
  "status": "error",
  "message": "Deskripsi error"
}
```

Kode yang umum digunakan: `200` sukses, `201` berhasil membuat data, `202` request asynchronous diterima, `400` input atau business rule tidak valid, `401` JWT tidak valid/tidak ada, `404` data tidak ditemukan, `409` konflik data, dan `500` error server/database.

## Dokumentasi Endpoint

### Category

#### GET `/categories` - Public

Mengambil seluruh kategori.

Response `200 OK`:

```json
{
  "status": "success",
  "message": "Categories retrieved successfully",
  "data": [
    { "id": 1, "category": "Pendidikan", "is_active": true }
  ]
}
```

#### POST `/categories` - Protected

Membuat kategori baru.

```json
{
  "category": "Kesehatan",
  "is_active": true
}
```

`category` wajib tidak kosong. Response `201 Created` berisi object kategori.

#### PUT `/categories/:id` - Protected

Mengubah kategori berdasarkan `id`.

```json
{
  "category": "Kesehatan dan Sosial",
  "is_active": true
}
```

Response `200 OK` berisi kategori yang telah diubah. ID tidak valid menghasilkan `400`; data tidak ditemukan menghasilkan `404`.

#### DELETE `/categories/:id` - Protected

Menghapus kategori berdasarkan `id`.

Response `200 OK`:

```json
{
  "status": "success",
  "message": "Category deleted successfully"
}
```

### Campaign

#### GET `/campaigns` - Public

Mengambil seluruh campaign.

Response `200 OK` mengembalikan array campaign dengan field `id`, `campaign`, `category_id`, `description`, `is_active`, `min_amount`, `target_amount`, `thumbnail`, `created_at`, dan `updated_at`.

#### GET `/campaigns/:id` - Public

Mengambil detail satu campaign. `id` harus berupa bilangan bulat positif. Response `404` dikembalikan jika campaign tidak ditemukan.

#### POST `/campaigns` - Protected

Membuat campaign baru.

```json
{
  "campaign": "Bantuan Pendidikan Anak",
  "category_id": 1,
  "description": "Pengadaan perlengkapan sekolah untuk anak-anak.",
  "is_active": true,
  "min_amount": 10000,
  "target_amount": 50000000,
  "thumbnail": "https://example.com/campaign.jpg"
}
```

Aturan:

- `campaign` wajib diisi.
- `category_id` harus positif, tersedia, dan aktif.
- `min_amount` tidak boleh negatif.
- `target_amount` harus lebih besar dari nol.
- `min_amount` tidak boleh lebih besar dari `target_amount`.

Response `201 Created` berisi campaign yang dibuat.

#### PUT `/campaigns/:id` - Protected

Mengubah seluruh data campaign berdasarkan `id`. Body dan aturan validasinya sama dengan endpoint create. Response `200 OK` berisi campaign yang diubah.

#### DELETE `/campaigns/:id` - Protected

Menghapus campaign berdasarkan `id`. Response `200 OK` berisi pesan sukses; campaign yang tidak ditemukan menghasilkan `404`.

### Campaign Update

#### GET `/campaigns/:id/updates` - Public

Mengambil seluruh update untuk campaign tertentu. Response `200 OK`:

```json
{
  "status": "success",
  "message": "Campaign updates retrieved successfully",
  "data": [
    {
      "id": 3,
      "campaign_id": 1,
      "title": "Penyaluran Dana Tahap 1",
      "content": "Dana telah disalurkan kepada penerima bantuan.",
      "created_at": "2026-08-26T10:00:00Z",
      "updated_at": "2026-08-26T10:00:00Z"
    }
  ]
}
```

#### POST `/campaigns/:id/updates` - Protected

Membuat update pada campaign aktif.

```json
{
  "title": "Penyaluran Dana Tahap 1",
  "content": "Dana telah disalurkan kepada penerima bantuan."
}
```

`title` dan `content` wajib diisi. Campaign harus ada dan aktif. Response `201 Created`.

#### PUT `/campaigns/:id/updates/:update_id` - Protected

Mengubah update tertentu yang dimiliki campaign tersebut. Body sama dengan endpoint create. Response `200 OK`; update yang tidak ditemukan atau tidak sesuai campaign menghasilkan `404`.

#### DELETE `/campaigns/:id/updates/:update_id` - Protected

Menghapus update berdasarkan kombinasi `id` campaign dan `update_id`. Response `200 OK`.

### Donor

Body donor digunakan oleh create dan update:

```json
{
  "donor": "Budi Santoso",
  "gender": "L",
  "email": "budi@example.com",
  "phone": "081234567890"
}
```

Aturan: seluruh field wajib diisi, email harus valid, gender harus `P` atau `L`, dan nomor telepon harus 8 sampai 15 karakter.

#### POST `/donors` - Protected

Membuat donor baru. Response `201 Created` berisi object donor. Nomor telepon yang sudah digunakan menghasilkan `400`.

#### GET `/donors` - Protected

Mengambil seluruh donor. Response `200 OK` mengembalikan array donor.

#### GET `/donors/:id` - Protected

Mengambil detail donor berdasarkan ID. ID harus positif; donor yang tidak ditemukan menghasilkan `404`.

#### PUT `/donors/:id` - Protected

Mengubah data donor berdasarkan ID. Body mengikuti format donor di atas. Response `200 OK`.

#### DELETE `/donors/:id` - Protected

Menghapus donor berdasarkan ID. Response `200 OK`; donor yang tidak ditemukan menghasilkan `404`.

### Transaction

#### POST `/transactions` - Public

Membuat transaksi donasi baru.

```json
{
  "donor": "Budi Santoso",
  "gender": "L",
  "email": "budi@example.com",
  "phone": "081234567890",
  "is_anonymous": false,
  "campaign_id": 1,
  "amount": 100000,
  "payment_method": "bank_transfer",
  "prayer": "Semoga campaign ini berjalan lancar"
}
```

Aturan:

- `campaign_id` harus menunjuk campaign yang aktif.
- `amount` harus lebih besar dari nol dan minimal sebesar `min_amount` campaign.
- Donor dicari berdasarkan email.
- Jika donor belum ada, `donor`, `gender`, `email`, dan `phone` wajib diisi.
- Gender donor baru harus `P` atau `L`.
- Nomor telepon donor baru harus 8 sampai 15 karakter.
- Email donor baru harus valid.
- Jika `is_anonymous` bernilai `true`, nama ditampilkan sebagai `Hamba Allah` pada response transaksi dan history.

Response `201 Created`:

```json
{
  "status": "success",
  "message": "Transaction created successfully",
  "data": {
    "id": 1,
    "donor_id": 1,
    "donor": "Budi Santoso",
    "email": "budi@example.com",
    "phone": "081234567890",
    "gender": "L",
    "is_anonymous": false,
    "campaign_id": 1,
    "amount": 100000,
    "payment_method": "bank_transfer",
    "prayer": "Semoga campaign ini berjalan lancar",
    "created_at": "2026-08-26T10:00:00Z",
    "updated_at": "2026-08-26T10:00:00Z"
  }
}
```

Contoh request untuk donor yang sudah terdaftar:

```json
{
  "email": "budi@example.com",
  "is_anonymous": true,
  "campaign_id": 1,
  "amount": 100000,
  "payment_method": "bank_transfer",
  "prayer": "Semoga bermanfaat"
}
```

#### GET `/transactions` - Protected

Mengambil seluruh transaksi beserta data donor. Response `200 OK` mengembalikan array transaksi.

#### PUT `/transactions/:id` - Protected

Mengubah campaign dan detail pembayaran transaksi.

```json
{
  "is_anonymous": false,
  "campaign_id": 2,
  "amount": 150000,
  "payment_method": "e_wallet",
  "prayer": "Semoga bermanfaat"
}
```

Campaign harus aktif, nominal harus valid dan memenuhi minimum campaign. Response `200 OK`; transaksi tidak ditemukan menghasilkan `404`.

#### DELETE `/transactions/:id` - Protected

Menghapus transaksi berdasarkan ID. Response `200 OK`; transaksi tidak ditemukan menghasilkan `404`.

#### GET `/transactions/history/:id_campaign` - Public

Mengambil riwayat donasi untuk campaign. Response `200 OK`:

```json
{
  "status": "success",
  "message": "Transaction history retrieved successfully",
  "data": [
    {
      "donor": "Hamba Allah",
      "prayer": "Semoga bermanfaat",
      "amount": 100000,
      "campaign_id": 1
    }
  ]
}
```

#### GET `/print-transaction/:transaction_id` - Public

Menghasilkan bukti donasi dalam format PDF. Response sukses memiliki `Content-Type: application/pdf` dan `Content-Disposition` inline dengan nama file `bukti-donasi-<id>.pdf`. Transaksi yang tidak ditemukan menghasilkan `404`.

### Notification

#### POST `/transactions/notifications/:id_transaksi` - Protected

Memulai pengiriman receipt melalui email secara asynchronous. Tidak membutuhkan request body.

Response `202 Accepted`:

```json
{
  "status": "success",
  "message": "Notification email is being sent",
  "email": "budi@example.com"
}
```

Endpoint hanya mengonfirmasi bahwa proses telah dimulai. Kegagalan SMTP dicatat pada log server. ID transaksi invalid menghasilkan `400`; transaksi tidak ditemukan menghasilkan `404`.

## Notifikasi Email

Konfigurasi email dapat disimpan pada database melalui migration notifikasi. Sistem membaca template dengan `type = 'email'` dan konfigurasi SMTP aktif. Placeholder template yang tersedia:

```text
@donatur @program @id_transaksi @tgldonasi @nominal
```

Jika konfigurasi SMTP database tidak tersedia, service menggunakan environment variable berikut:

```dotenv
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@example.com
SMTP_PASSWORD=your-gmail-app-password
SMTP_SENDER_NAME=Donation System
SMTP_SENDER_EMAIL=your-email@example.com
```

Gunakan Gmail App Password, bukan password akun Gmail biasa. Email dikirim dengan lampiran PDF receipt.

## Struktur Direktori

```text
SB-Golang-DonationSystem/
├── asset/img/                 # Logo dan tanda tangan PDF
├── cmd/server/main.go         # Entry point, dependency wiring, dan route
├── internal/
│   ├── campaign/              # Campaign handler, service, repository, model
│   ├── campaign_update/       # Update informasi campaign
│   ├── category/              # Manajemen kategori
│   ├── config/                # Environment dan koneksi PostgreSQL
│   ├── donor/                 # Manajemen donatur
│   ├── middleware/            # JWT middleware
│   ├── notification/          # Template, SMTP, dan email asynchronous
│   ├── pdf/                   # Generator receipt PDF
│   ├── transaction/           # Transaksi dan history donasi
│   └── user/                  # Register dan login
├── migrations/                # SQL migration
├── .env                       # Konfigurasi lokal, jangan dipublikasi
├── go.mod                     # Module dan dependency Go
└── README.md                  # Dokumentasi API
```

Setiap domain menggunakan pola:

```text
handler.go       # HTTP request dan response
service.go       # Business rule
repository.go    # Query database
model.go         # Model dan request struct
```
