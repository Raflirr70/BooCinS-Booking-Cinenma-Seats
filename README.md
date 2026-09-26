# BooCinS

Backend REST API untuk sistem reservasi tiket bioskop berbasis web (BooCinS — Booking Cinema Seats). Dibangun dengan **Go** menggunakan arsitektur **Clean Architecture**, di mana tiap fitur dipisahkan ke dalam domain, use case, repository, dan delivery layer agar mudah dikembangkan dan diuji.

## Daftar Isi

- [Fitur](#fitur)
- [Tech Stack](#tech-stack)
- [Arsitektur](#arsitektur)
- [Struktur Folder](#struktur-folder)
- [Persyaratan](#persyaratan)
- [Cara Menjalankan](#cara-menjalankan)
- [Konfigurasi Environment](#konfigurasi-environment)
- [Database & Migration](#database--migration)
- [API Endpoint](#api-endpoint)
- [Contoh Request](#contoh-request)
- [Pendekatan Aplikasi](#pendekatan-aplikasi)
- [Testing](#testing)

## Fitur

Fitur sistem berdasarkan:

- **Manager** — monitoring pendapatan (online/offline), monitoring pembelian tiket, kelola akun.
- **Admin** — kelola film, kelola jadwal ruangan & film, kelola promosi, kelola diskon membership, kelola genre.
- **Staff** — pembelian tiket offline.
- **Consumer (member)** — monitoring skor membership, histori tontonan, histori pencarian, kelola bookmark, rekomendasi film.
- **Guest** — pemesanan tiket online.

Fitur yang diimplementasikan di kode:

- Autentikasi (register, login, logout, profile) — **JWT + Redis session/blacklist**
- Manajemen **Film** (CRUD + detail with genres & media)
- Manajemen **Genre** (CRUD)
- Room & Seat (create/update/delete dengan transaksi & soft delete)
- Booking tiket, schedule, promo, memberships, rekomendasi, dll. *(sedang dikembangkan)*
- Pembayaran

## Tech Stack

| Teknologi | Kegunaan |
|---|---|
| **Go** (1.26) | Bahasa pemrograman |
| **Gin** (`github.com/gin-gonic/gin`) | Web framework / HTTP router |
| **GORM** (`gorm.io/gorm`) + `gorm.io/driver/postgres` | ORM & migration schema |
| **PostgreSQL** | Database utama |
| **Redis** (`github.com/go-redis/redis/v8`) | Session/blacklist token, cache |
| **JWT** (`github.com/golang-jwt/jwt/v5`) | Autentikasi token |
| **bcrypt** (`golang.org/x/crypto`) | Hash password |
| **godotenv** (`github.com/joho/godotenv`) | Baca file `.env` |
| **Kafka** (`config.KafkaConfig`) | Konfigurasi messaging tersedia (implementasi menyusul) |
| **miniredis** (`github.com/alicebob/miniredis/v2`) | Redis in-memory untuk unit test |
| **testify** (`github.com/stretchr/testify`) | Assertion untuk unit test |

## Arsitektur

Proyek mengikuti **Clean Architecture**, dengan aliran dependensi satu arah:

```
Delivery layer (handler, middleware, router, DTO)
        ↓
Usecase layer (business logic, transaction)
        ↓
Repository layer (interface) → implementasi (gorm/postgres)
        ↓
Domain layer (entity)
```

- **Domain** berisi entity murni (tanpa dependency framework) serta *contract* interface repository & usecase.
- **Usecase** memuat aturan bisnis, contoh: transaksi create/update room, attach genre ke film.
- **Delivery** berisi handler HTTP, DTO, middleware, dan router.
- **Infrastructure** berisi koneksi DB (dengan AutoMigrate & seed), Redis, dan messaging.

### Transaksi & Rollback

Operasi multi-step dibungkus dalam transaksi GORM agar bisa **rollback** jika terjadi error di tengah proses. Contohnya `Create`/`Update` room: jika ada seat yang gagal dibuat/diupdate, seluruh perubahan (room + seat) dibatalkan.

```go
u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    roomRepo := u.roomRepo.WithTx(tx)
    seatRepo := u.seatRepo.WithTx(tx)
    // ... semua operasi lewat repo yang terikat tx
})
```

### Soft Delete

Beberapa entitas menggunakan **soft delete**: baris tidak dihapus fisik, hanya diberi timestamp `deleted_at`. Semua query otomatis mengecualikan data yang sudah di-soft-delete.

## Struktur Folder

```
BooCinS
├── cmd/
│   └── api/
│       └── main.go                     // entry point, wiring dependency
├── config/
│   └── config.go                       // baca konfigurasi dari environment
├── internal/
│   ├── domain/                         // layer domain (entity + contract)
│   │   ├── entity/
│   │   │   ├── user.go
│   │   │   ├── film.go
│   │   │   ├── schedule.go             
│   │   │   ├── ticket.go
│   │   │   ├── interaction.go          
│   │   │   └── promo.go
│   │   ├── repository/
│   │   └── usecase/
│   ├── delivery/http/
│   │   ├── dto/                        // request & response DTO
│   │   ├── handler/                    // auth, film, genre handler
│   │   ├── middleware/                 // JWT auth & role middleware
│   │   └── router/router.go            // daftar route
│   ├── infrastructure/
│   │   ├── cache/redis.go              // koneksi Redis
│   │   └── database/postgres.go        // koneksi, AutoMigrate, seed
│   ├── repository/postgres/            // implementasi repository (GORM)
│   └── usecase/                        // implementasi business logic + test
├── pkg/
│   ├── jwt/
│   ├── logger/
│   └── response/
├── test/mocks/                         // mock repository untuk unit test
├── .env.example
```

## Persyaratan

- **Go** 1.26 atau lebih baru
- **PostgreSQL**
- **Redis** 6.x

## Cara Menjalankan

1. Clone repository dan masuk ke direktori:

   ```bash
   git clone https://github.com/rafli/boocins.git
   cd BooCinS
   ```

2. Salin `.env.example` menjadi `.env` dan sesuaikan nilainya:

   ```bash
   cp .env.example .env
   ```

3. Pastikan PostgreSQL dan Redis sudah berjalan di mesin Anda.

4. Jalankan server:

   ```bash
   go run ./cmd/api
   ```

5. Server berjalan di `http://localhost:8080` (atau sesuai `SERVER_PORT`).

## Konfigurasi Environment

| Variabel | Default | Keterangan |
|---|---|---|
| `SERVER_PORT` | `8080` | Port server |
| `DB_HOST` | `localhost` | Host PostgreSQL |
| `DB_PORT` | `5432` | Port PostgreSQL |
| `DB_USER` | `user-kau` | User PostgreSQL |
| `DB_PASSWORD` | `pw-kau` | Password PostgreSQL |
| `DB_NAME` | `db-kau` | Nama database |
| `DB_SSLMODE` | `disable` | Mode SSL koneksi DB |
| `REDIS_HOST` | `localhost` | Host Redis |
| `REDIS_PORT` | `6379` | Port Redis |
| `REDIS_PASSWORD` | *(kosong)* | Password Redis |
| `REDIS_DB` | `0` | Nomor database Redis |
| `JWT_SECRET` | `secret-jwt-kau` | Secret JWT (ganti di production!) |
| `JWT_EXPIRATION_HOURS` | `24` | Masa berlaku token (jam) |
| `KAFKA_BROKER` | `localhost:9092` | Broker Kafka |
| `KAFKA_GROUP_ID` | `boocins-group` | Group ID Kafka |

## Database & Migration

- **Jangan** menjalankan `Database.sql` secara manual ke database yang sama dengan AutoMigrate (akan konflik tabel).
- Schema dibuat otomatis oleh **GORM AutoMigrate** saat server pertama kali dijalankan.
- AutoMigrate juga menambahkan kolom baru (mis. `deleted_at` untuk soft delete) saat schema berubah.
- Default data di-seed otomatis:
  - **Role**: `super_admin`, `manager`, `admin`, `staff`, `member`
  - **Membership**: `Bronze`, `Silver`, `Gold`, `Platinum`

### Catatan index unik & soft delete

Jika ingin nama room (mis. "Studio 1") bisa dipakai kembali setelah room di-soft-delete, ubah `uniqueIndex` biasa pada `Room.Name` menjadi index komposit `(name, deleted_at)`. GORM AutoMigrate **tidak** menghapus constraint/index lama, jadi sekali-sekali jalankan SQL manual:

```sql
ALTER TABLE rooms DROP CONSTRAINT rooms_name_key;
```

## Contoh Request

### Register

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Rafli",
    "email": "rafli@example.com",
    "password": "123456789"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "rafli@example.com",
    "password": "123456789"
  }'
```

Response berisi `access_token` yang dipakai pada header `Authorization: Bearer <token>`.

## Pendekatan Aplikasi

- **DTO terpisah dari entity** — request/response memakai DTO di `internal/delivery/http/dto`, sehingga entity domain tetap murni.
- **Soft delete** untuk data historis (user, film, room, schedule).
- **Transaksi database** untuk operasi yang melibatkan banyak baris/tabel.
- **Redis** dipakai untuk mengecek apakah token masih valid (session) dan blacklist token saat logout.
- **Unit test** memakai `miniredis` (tanpa Redis asli) dan mock repository.

## Testing

```bash
go test ./... -v
```