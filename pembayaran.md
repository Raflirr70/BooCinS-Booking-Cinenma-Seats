# Alur Pembuatan Ticket + Pembayaran Midtrans

## Overview

```
Frontend (Postman / Web)
  │
  ▼
POST /api/v1/booking
  │
  ├─ 1. Validasi input (schedule_id, seat_ids[], first_name, last_name)
  ├─ 2. Lock seat (SELECT FOR UPDATE → status = 'locked')
  ├─ 3. Hitung total harga (jumlah seat × harga film - diskon membership)
  ├─ 4. Buat Transaction (status = 'pending')
  ├─ 5. Buat Ticket per seat (status = 'active')
  ├─ 6. Panggil Midtrans Snap → dapat snap_token + redirect_url
  ├─ 7. Return snap_token ke frontend (popup Midtrans)
  │
  ▼
POST /api/v1/booking/notification  ← Midtrans webhook
  │
  ├─ settlement  → Transaction.status = 'paid', ScheduleSeat.status = 'sold'
  ├─ expire      → Transaction.status = 'cancelled', rollback seat → 'available'
  ├─ cancel/deny → Transaction.status = 'cancelled', rollback seat → 'available'
  │                 + jika sudah terbayar → Midtrans Refund API
  ▼
Done
```

---

## 1. Data yang Dikirim dari Frontend

### Request: `POST /api/v1/booking`

Seat yang dipesan diambil dari response `GET /api/v1/schedules/:id/seats` — user memilih (ceklis) seat yang `is_available: true`.

```json
{
  "schedule_id": 1,
  "seat_ids": [3, 5, 7],
  "first_name": "John",
  "last_name": "Doe"
}
```

### Aturan field `first_name` dan `last_name`

| Kondisi | Perilaku |
|---|---|
| **Sudah login** (token JWT ada) | `first_name` dan `last_name` **diabaikan dari body**, langsung diambil dari `users` table berdasarkan `user_id` di JWT claim. `Transaction.user_id` diisi. |
| **Belum login** (guest) | `first_name` dan `last_name` **wajib diisi** di body. Buat `GuestOrder` dengan email (opsional, untuk kirim e-ticket). `Transaction.guest_order_id` diisi. |

### Perhitungan Total Harga

```
price_per_seat  = Film.Price
jumlah_seat     = len(seat_ids)
subtotal        = price_per_seat × jumlah_seat

// Diskon membership (hanya user login yang punya membership)
discount        = User.Membership.Discount  // contoh: 10.00 = 10%
total_price     = subtotal × (1 - discount/100)
```

Contoh: Film harga Rp 50.000, user pesan 3 seat, membership Gold (diskon 10%)
→ `50000 × 3 × 0.90 = Rp 135.000`

---

## 2. Alur Detail di Backend (Satu DB Transaction)

Semua langkah di bawah ini berjalan di dalam **satu `gorm.DB.Transaction()`**. Jika ada error di langkah mana pun, seluruh query di-rollback otomatis.

### Step 1 — Validasi Schedule

```go
schedule, err := scheduleRepo.FindByID(ctx, req.ScheduleID)
// Pastikan schedule ada dan status = true (aktif)
```

### Step 2 — Lock Seat (Mencegah Bentrok / Race Condition)

```sql
SELECT * FROM schedule_seats
WHERE schedule_id = ? AND seat_id IN (?, ?, ?)
AND status = 'available'
FOR UPDATE
```

Menggunakan `FOR UPDATE` agar row-level lock aktif di PostgreSQL. Request lain yang mencoba lock seat yang sama akan **menunggu** sampai transaksi ini selesai (commit/rollback).

```go
// Di repository, buat method baru:
func (r *scheduleSeatRepository) LockSeats(ctx context.Context, scheduleID uint, seatIDs []uint) ([]entity.ScheduleSeat, error) {
    var seats []entity.ScheduleSeat
    err := r.db.WithContext(ctx).
        Raw(`SELECT * FROM schedule_seats
             WHERE schedule_id = ? AND seat_id IN ?
             AND status = 'available'
             FOR UPDATE`, scheduleID, seatIDs).
        Scan(&seats).Error
    return seats, err
}
```

**Validasi setelah lock:**
- Jumlah seat yang berhasil di-lock harus == `len(seat_ids)` dari request.
- Jika kurang → salah satu seat sudah diambil orang lain → return error `"satu atau lebih seat sudah tidak tersedia"` → rollback.

### Step 3 — Update Status Seat → `locked`

```go
now := time.Now()
for _, ss := range lockedSeats {
    ss.Status = "locked"
    ss.LockedAt = &now
    scheduleSeatRepo.Update(ctx, &ss)
}
```

### Step 4 — Hitung Harga

```go
film := schedule.Film
pricePerSeat := film.Price
subtotal := pricePerSeat * float64(len(seatIDs))

var discount float64
if user != nil && user.Membership != nil {
    discount = user.Membership.Discount
}
totalPrice := subtotal * (1 - discount/100)
```

### Step 5 — Buat Transaction Record

```go
transaction := &entity.Transaction{
    UserID:     userID,        // nil jika guest
    GuestOrderID: guestOrderID, // nil jika login
    Status:     "pending",
    TotalPrice: totalPrice,
    Source:     "online",
}
transactionRepo.Create(ctx, transaction)
```

### Step 6 — Buat Ticket Per Seat

```go
for _, ss := range lockedSeats {
    ticket := &entity.Ticket{
        UserID:         userID,
        ScheduleSeatID: ss.ID,
        TransactionID:  transaction.ID,
        QRToken:        uuid.New().String(),
        Status:         "active",
    }
    ticketRepo.Create(ctx, ticket)
}
```

### Step 7 — Panggil Midtrans Snap

```go
import "github.com/midtrans/midtrans-go/snap"

snapClient := snap.Client{}
snapClient.New(cfg.Midtrans.ServerKey, midtrans.Sandbox) // atau Production

snapReq := &snap.Request{
    TransactionDetails: midtrans.TransactionDetails{
        OrderID:  fmt.Sprintf("BOOCINS-%d-%d", transaction.ID, time.Now().Unix()),
        GrossAmt: int64(totalPrice),
    },
    CustomerDetail: &midtrans.CustomerDetail{
        FName: firstName,
        LName: lastName,
    },
}

snapResp, err := snapClient.CreateTransaction(snapReq)
if err != nil {
    return err // → rollback semua (seat kembali available, transaction & ticket dihapus)
}
```

### Step 8 — Simpan Snap Token, Commit Transaction

```go
// Simpan order_id dan snap_token di transaction (tambah field jika perlu)
transaction.MidtransOrderID = snapReq.TransactionDetails.OrderID
transaction.SnapToken = snapResp.Token
transactionRepo.Update(ctx, transaction)

// return ke frontend
return snapResp.Token, snapResp.RedirectURL
```

**Jika langkah 7 (Midtrans) gagal:** `gorm.Transaction()` akan rollback otomatis karena function return error. Semua seat kembali `available`, transaction dan ticket tidak tersimpan.

---

## 3. Response ke Frontend

```json
{
  "status": "success",
  "message": "booking created",
  "data": {
    "transaction_id": 42,
    "order_id": "BOOCINS-42-1696500000",
    "snap_token": "d1f2e3a4-...",
    "redirect_url": "https://app.sandbox.midtrans.com/snap/v2/vtweb/d1f2e3a4-...",
    "total_price": 135000,
    "seats": [
      {"seat_id": 3, "label": "A", "number": 3},
      {"seat_id": 5, "label": "A", "number": 5},
      {"seat_id": 7, "label": "B", "number": 2}
    ]
  }
}
```

Frontend menampilkan popup Midtrans menggunakan `snap_token`:
```js
snap.pay(data.snap_token, {
    onSuccess: function(result) { /* redirect ke halaman sukses */ },
    onPending: function(result) { /* tampilkan "menunggu pembayaran" */ },
    onError:   function(result) { /* tampilkan error */ }
});
```

---

## 4. Webhook Midtrans (Notification Handler)

### Endpoint: `POST /api/v1/booking/notification`

Midtrans mengirim notifikasi ke URL ini setiap kali status pembayaran berubah. **Endpoint ini publik** (tanpa auth), tapi diverifikasi menggunakan signature key dari Midtrans.

```go
// Verifikasi signature
orderId := notification["order_id"].(string)
statusCode := notification["status_code"].(string)
grossAmount := notification["gross_amount"].(string)
signatureKey := notification["signature_key"].(string)

hash := sha512(orderId + statusCode + grossAmount + cfg.Midtrans.ServerKey)
if hash != signatureKey {
    return 403 // signature tidak valid, abaikan
}
```

### Mapping Status Midtrans → Aksi

| Midtrans `transaction_status` | Aksi di BooCinS |
|---|---|
| `capture` / `settlement` | `Transaction.status = 'paid'`, semua `ScheduleSeat.status = 'sold'`, `ScheduleSeat.locked_at = NULL` |
| `pending` | Tidak ada perubahan (sudah `pending`) |
| `expire` | `Transaction.status = 'cancelled'`, semua `ScheduleSeat.status = 'available'`, semua `Ticket.status = 'cancelled'` |
| `cancel` / `deny` | Sama seperti `expire` |

### Jika Midtrans Sudah Terlanjur Berhasil Tapi Error di Tengah

Skenario: pembayaran sukses di Midtrans, tapi update DB gagal (misal DB timeout).

**Solusi: idempotent webhook + retry dari Midtrans.**

Midtrans akan mengirim ulang notifikasi sampai mendapat response HTTP 200. Jadi:
1. Webhook handler harus **idempotent** — panggilan berulang dengan `order_id` yang sama tidak mengubah state yang sudah benar.
2. Cek `Transaction.status` sebelum update: jika sudah `paid`, langsung return 200.

### Jika Perlu Refund (Pembatalan Setelah Bayar)

```go
import "github.com/midtrans/midtrans-go/coreapi"

coreClient := coreapi.Client{}
coreClient.New(cfg.Midtrans.ServerKey, midtrans.Sandbox)

refundReq := &coreapi.RefundReq{
    Amount: int64(transaction.TotalPrice),
    Reason: "Pembatalan oleh sistem",
}
_, err := coreClient.RefundTransaction(transaction.MidtransOrderID, refundReq)
```

Setelah refund berhasil:
- `Transaction.status = 'refunded'`
- Semua `ScheduleSeat.status = 'available'`
- Semua `Ticket.status = 'cancelled'`

---

## 5. Pencegahan Bentrok Pemesanan (Race Condition)

### Mekanisme: `SELECT ... FOR UPDATE` + DB Transaction

```
User A (seat 3,5)              User B (seat 5,7)
      │                              │
      ▼                              ▼
 BEGIN TX                        BEGIN TX
      │                              │
 SELECT ... FOR UPDATE           SELECT ... FOR UPDATE
 seat 3 ✓ locked                 seat 5 → MENUNGGU (blocked)
 seat 5 ✓ locked                      │
      │                              │
 UPDATE status='locked'               │ (masih nunggu TX A selesai)
      │                              │
 Midtrans call                        │
      │                              │
 COMMIT ─────────────────────►   seat 5 sudah 'locked'
                                 → jumlah seat di-lock ≠ yang diminta
                                 → return error "seat tidak tersedia"
                                 ROLLBACK
```

Tidak ada polling, tidak ada Redis lock, tidak ada optimistic locking — cukup `FOR UPDATE` di PostgreSQL.

### Proteksi Tambahan: Lock Timeout

Supaya user B tidak menunggu selamanya jika user A hang:

```go
tx.Exec("SET LOCAL lock_timeout = '5s'")
```

Jika tidak bisa mendapatkan lock dalam 5 detik → error → rollback → frontend menampilkan "coba lagi".

### Proteksi Tambahan: Expired Lock Cleanup

Seat yang berstatus `locked` tapi `locked_at` sudah > 15 menit (user tidak membayar) harus dikembalikan ke `available`. Ini dilakukan oleh:

```go
// Cron job atau goroutine berkala setiap 1 menit:
db.Exec(`UPDATE schedule_seats
         SET status = 'available', locked_at = NULL
         WHERE status = 'locked'
         AND locked_at < NOW() - INTERVAL '15 minutes'`)
```

Dan update Transaction terkait menjadi `cancelled`, Ticket menjadi `cancelled`.

---

## 6. File dan Lokasi yang Perlu Dibuat / Diubah

### File Baru

| File | Isi |
|---|---|
| `internal/domain/repository/ticket_repository.go` | Interface `TransactionRepository` dan `TicketRepository` |
| `internal/repository/postgres/ticket_repository.go` | Implementasi repo: `Create`, `FindByID`, `Update`, `LockSeats` (di ScheduleSeat repo) |
| `internal/domain/usecase/ticket_usecase.go` | Interface `BookingUsecase` |
| `internal/usecase/ticket_usecase.go` | Implementasi: `CreateBooking()` dengan `db.Transaction()`, Midtrans Snap call, hitung harga |
| `internal/delivery/http/handler/booking_handler.go` | Handler: `CreateBooking` dan `MidtransNotification` |
| `internal/delivery/http/dto/booking_dto.go` | DTO: `CreateBookingRequest`, `BookingResponse`, `MidtransNotificationRequest` |

### File yang Diubah

| File | Perubahan |
|---|---|
| `internal/domain/repository/schedule_repository.go` | Tambah method `LockSeats(ctx, scheduleID, seatIDs) ([]ScheduleSeat, error)` di `ScheduleSeatRepository` |
| `internal/repository/postgres/schedule_repository.go` | Implementasi `LockSeats` dengan `SELECT ... FOR UPDATE` |
| `internal/delivery/http/router/router.go` | Tambah route `POST /booking` (protected) dan `POST /booking/notification` (publik) |
| `cmd/api/main.go` | Wire `TransactionRepo`, `TicketRepo`, `BookingUsecase`, `BookingHandler` |
| `internal/domain/entity/ticket.go` | Tambah field `MidtransOrderID` dan `SnapToken` di `Transaction` |

### Route Baru di Router

```go
// Booking (user login)
protected.POST("/booking", bookingHandler.CreateBooking)

// Booking (guest — tanpa auth)
api.POST("/booking/guest", bookingHandler.CreateBookingGuest)

// Midtrans webhook (publik, diverifikasi lewat signature)
api.POST("/booking/notification", bookingHandler.MidtransNotification)
```

---

## 7. Urutan Implementasi

1. Tambah field `MidtransOrderID` + `SnapToken` di entity `Transaction`
2. Buat `TransactionRepository` + `TicketRepository` (interface + implementasi)
3. Tambah `LockSeats()` di `ScheduleSeatRepository`
4. Buat `BookingUsecase` — inti logika ada di sini
5. Buat DTO (`CreateBookingRequest`, `BookingResponse`)
6. Buat `BookingHandler` (endpoint + webhook)
7. Wire di `main.go` dan daftarkan route di `router.go`
8. Buat cleanup goroutine untuk expired locked seats

---

## 8. Contoh Log Query (Rollback Scenario)

### Skenario: Midtrans gagal di langkah 7

```
[INFO] BEGIN TRANSACTION
[INFO] SELECT * FROM schedules WHERE id = 1
[INFO] SELECT * FROM schedule_seats WHERE schedule_id = 1 AND seat_id IN (3,5,7) AND status = 'available' FOR UPDATE
[INFO] UPDATE schedule_seats SET status = 'locked', locked_at = NOW() WHERE id IN (10, 12, 14)
[INFO] INSERT INTO transactions (user_id, status, total_price, source) VALUES (5, 'pending', 135000, 'online')
[INFO] INSERT INTO tickets (user_id, schedule_seat_id, transaction_id, qr_token, status) VALUES ...
[ERROR] Midtrans Snap API failed: timeout
[INFO] ROLLBACK  ← semua INSERT dan UPDATE dibatalkan, seat kembali 'available'
```

### Skenario: Pembayaran sukses tapi DB update gagal

```
[INFO] Midtrans notification received: order_id=BOOCINS-42-1696500000, status=settlement
[INFO] BEGIN TRANSACTION
[INFO] SELECT * FROM transactions WHERE midtrans_order_id = 'BOOCINS-42-1696500000'
[INFO] UPDATE transactions SET status = 'paid' WHERE id = 42
[ERROR] UPDATE schedule_seats SET status = 'sold' → DB connection lost
[INFO] ROLLBACK  ← transaction.status kembali ke 'pending'
[INFO] Midtrans akan mengirim ulang notifikasi karena response bukan 200
[INFO] (Retry) notification received → kali ini DB available → COMMIT → sukses
```

### Skenario: Perlu refund karena error setelah pembayaran

```
[INFO] Midtrans notification: settlement
[INFO] BEGIN TRANSACTION
[INFO] UPDATE transactions SET status = 'paid'
[INFO] UPDATE schedule_seats SET status = 'sold'
[ERROR] INSERT tickets → unique constraint violation (duplikat qr_token)
[INFO] ROLLBACK
[INFO] Transaction sudah 'paid' di Midtrans tapi DB gagal → trigger refund
[INFO] POST https://api.sandbox.midtrans.com/v2/BOOCINS-42-1696500000/refund
[INFO] Refund berhasil → UPDATE transactions SET status = 'refunded'
[INFO] UPDATE schedule_seats SET status = 'available'
```
