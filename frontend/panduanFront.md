# Panduan Frontend BooCinS (React + Vite + TS)

> File ini WAJIB dibaca oleh AI dan programmer sebelum sentuh kode di folder `frontend/`.
> Bahasa dibuat sederhana. Ikuti kata per kata. Jangan improvisasi.

## 0. Tujuan

1. Halaman utama (`/`) ambil data asli dari `GET /api/v1/home`.
2. Login (`POST /api/v1/auth/login`) lalu lempar user sesuai `role`:
   - `admin` -> `/admin` (dashboard + sidebar: Kelola Film, Room, Schedule, Promo)
   - `manager` -> `/manager` (dashboard + statistik pendapatan, booking, schedule)
   - `super_admin` -> `/superadmin` (dashboard + log aktivitas + statistik pengunjung)
   - `staff` -> `/staff` (kasir pembelian tiket cash)
   - `member` -> `/` (halaman utama lagi)
3. Dashboard `super_admin, admin, staff, manager` = TANPA fetch API dulu. Pakai data hardcode (tulis manual di file).

---

## 1. ATURAN ANTI-HALUSINASI (WAJIB, TIDAK BOLEH DILANGGAR)

AI dan programmer DILARANG:

1. DILARANG bikin folder baru di luar daftar di Bagian 2. Kalau butuh folder baru, tanya manusia dulu.
2. DILARANG pakai `fetch()`. Hanya boleh pakai `axios` dari file `src/lib/api.ts`.
3. DILARANG pakai `localStorage` langsung. Hanya boleh pakai `src/store/auth.ts` (zustand + persist).
4. DILARANG bikin state auth baru pakai `Context`, `Redux`, atau `useState` global. Auth hanya di `src/store/auth.ts`.
5. DILARANG ganti library. Stack dikunci di Bagian 3. Jangan tambah `redux, swr, react-query, next, shadcn` tanpa izin.
6. DILARANG ubah bentuk response backend. Backend selalu balas bentuk ini:
   ```json
   { "status": "success", "message": "...", "data": { ... } }
   ```
7. DILARANG ubah nama role. Role yang sah hanya 5 string ini (huruf kecil, persis):
   `super_admin | manager | admin | staff | member`
8. DILARANG ubah nama file yang sudah ada: `src/lib/api.ts`, `src/store/auth.ts`, `src/components/RequireAuth.tsx`, `src/App.tsx`.
9. Penamaan file halaman: huruf kecil + underscore, contoh: `home_page.tsx`, `admin_films_page.tsx`. Satu file = satu halaman.
10. Setiap halaman WAJIB punya 3 state: `loading`, `error`, `data`. Tidak boleh fetch tanpa `try/catch`.

Kalau AI tidak tahu, AI harus jawab: "Tidak ada di panduan, saya berhenti" — jangan mengarang API baru.

---

## 2. Struktur Folder Baku (Jangan Diubah)

```
frontend/src/
  App.tsx                  -> daftar Route saja, tidak ada logic fetch di sini
  main.tsx                 -> jangan diubah
  lib/api.ts               -> satu-satunya axios instance, baseURL: /api/v1
  store/auth.ts            -> simpan token + user + fungsi homePathFor()
  components/
    RequireAuth.tsx        -> penjaga route per role
    Navbar.tsx             -> navbar publik
    Sidebar.tsx            -> sidebar khusus /admin
    FilmCard.tsx           -> kartu film di home
    PromoBanner.tsx        -> banner promo di home
  layouts/
    PublicLayout.tsx       -> bungkus halaman publik (navbar + outlet)
    AdminLayout.tsx        -> bungkus halaman admin (sidebar + outlet)
    DashboardLayout.tsx    -> bungkus manager/superadmin/staff (tanpa sidebar kompleks)
  pages/
    auth/login_page.tsx    -> SUDAH ADA, tinggal perbaiki redirect role
    auth/register_page.tsx -> SUDAH ADA, jangan diubah
    home/home_page.tsx     -> BARU, fetch /home
    admin/admin_dashboard_page.tsx      -> BARU, hardcode
    admin/admin_films_page.tsx          -> BARU, hardcode
    admin/admin_rooms_page.tsx          -> BARU, hardcode
    admin/admin_schedules_page.tsx      -> BARU, hardcode
    admin/admin_promos_page.tsx         -> BARU, hardcode
    manager/manager_dashboard_page.tsx  -> BARU, hardcode
    superadmin/superadmin_dashboard_page.tsx -> BARU, hardcode
    staff/staff_cashier_page.tsx        -> BARU, hardcode
  data/
    mock_admin.ts          -> semua hardcode admin taruh sini
    mock_manager.ts        -> semua hardcode manager taruh sini
    mock_superadmin.ts     -> semua hardcode superadmin taruh sini
    mock_staff.ts          -> semua hardcode staff taruh sini
```

Aturan: komponen kecil di `components/`, halaman di `pages/`, data palsu di `data/`. Jangan campur.

---

## 3. Stack Terkunci

- `react` + `react-dom` v19
- `react-router-dom` v7 (pakai `BrowserRouter, Routes, Route, Navigate, Outlet, Link, useNavigate`)
- `axios` (hanya lewat `src/lib/api.ts`)
- `zustand` (hanya untuk auth)
- `tailwindcss` v4 (untuk gaya, jangan pakai CSS file baru kalau bisa pakai class Tailwind)
- `vite` proxy: `/api` -> `http://localhost:8080` (sudah diset di `vite.config.ts`, jangan diubah)

File `src/lib/api.ts` yang benar (jangan diubah polanya):

```ts
import axios from 'axios'
import { useAuth } from '../store/auth'
const api = axios.create({ baseURL: '/api/v1' })
api.interceptors.request.use((config) => {
  const token = useAuth.getState().token
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})
export default api
```

File `src/store/auth.ts` yang benar (wajib seperti ini):

```ts
export interface AuthUser {
  id: number
  first_name: string
  last_name: string
  email: string
  role: { id: number; name: string }
}
// role.name pasti salah satu dari: super_admin, manager, admin, staff, member

export function homePathFor(role: string): string {
  if (role === 'super_admin') return '/superadmin'
  if (role === 'admin') return '/admin'
  if (role === 'manager') return '/manager'
  if (role === 'staff') return '/staff'
  return '/' // member dan guest ke halaman utama
}
```

File `src/components/RequireAuth.tsx` yang benar (perbaiki file lama yang rusak):

```tsx
import type { ReactNode } from 'react'
import { Navigate } from 'react-router-dom'
import { useAuth, homePathFor } from '../store/auth'

export default function RequireAuth({ roles, children }: { roles?: string[]; children: ReactNode }) {
  const token = useAuth((s) => s.token)
  const user = useAuth((s) => s.user)
  if (!token || !user) return <Navigate to="/login" replace />
  if (roles && !roles.includes(user.role.name)) {
    return <Navigate to={homePathFor(user.role.name)} replace />
  }
  return <>{children}</>
}
```

---

## 4. Kontrak API (Jangan Dikarang)

### 4.1 Login

- Request: `POST /api/v1/auth/login`
- Body: `{ "email": "a@b.com", "password": "123456" }`
- Balasan sukses (`res.data.data`):
  ```json
  { "user": { "id": 1, "first_name": "...", "last_name": "...", "email": "...", "role": { "id": 1, "name": "admin" } }, "token": "jwt..." }
  ```
- Cara pakai di `login_page.tsx`:
  ```tsx
  const res = await api.post('/auth/login', { email, password })
  const { user, token } = res.data.data as { user: AuthUser; token: string }
  setAuth(token, user)
  navigate(homePathFor(user.role.name), { replace: true })
  ```
- Jangan `navigate('/')` untuk semua role. Wajib pakai `homePathFor()`.

### 4.2 Home

- Request: `GET /api/v1/home` (publik, boleh tanpa token)
- Balasan sukses (`res.data.data`):
  ```json
  {
    "user_id": 0, "name": "",
    "promo": { "id": 1, "title": "...", "description": "...", "img": "..." },
    "films": [
      { "film_id": 1, "name": "...", "cover": "...", "duration": 120, "price": 50000, "genres": ["Action"], "schedule": [
        { "schedule_id": 1, "room_id": 1, "room_name": "Studio 1", "date": "2026-10-07", "time": "19:00", "nseat": 100, "booked_seats": 10, "available_seats": 90, "seats": 10 }
      ]}
    ]
  }
  ```
- `promo` bisa `null`. Kode harus cek `if (promo)` sebelum tampilkan.
- Type TS yang wajib dipakai di `home_page.tsx`:
  ```ts
  interface ScheduleHome { schedule_id: number; room_name: string; date: string; time: string; available_seats: number }
  interface FilmHome { film_id: number; name: string; cover: string; duration: number; price: number; genres: string[]; schedule: ScheduleHome[] }
  interface HomeData { user_id: number; name: string; promo: { id: number; title: string; description: string; img: string } | null; films: FilmHome[] }
  ```
- Pola fetch baku:
  ```tsx
  const [data, setData] = useState<HomeData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  useEffect(() => {
    (async () => {
      try {
        const res = await api.get('/home')
        setData(res.data.data as HomeData)
      } catch (e) { setError(apiMessage(e)) }
      finally { setLoading(false) }
    })()
  }, [])
  if (loading) return <p>Loading...</p>
  if (error) return <p>{error}</p>
  if (!data || data.films.length === 0) return <p>Belum ada film</p>
  ```

---

## 5. RUNDOWN PENGERJAAN (Kerjakan Urut, 1 Tahap = 1 Prompt AI)

### Tahap 0 — Siap-siap (5 menit)
1. Jalankan backend: `go run ./cmd/api` (pastikan jalan di `http://localhost:8080`).
2. Jalankan frontend: `cd frontend; npm install; npm run dev`.
3. Buka `http://localhost:5173/` dan pastikan halaman BooCinS muncul.

### Tahap 1 — Perbaiki Login + Routing Role (Wajib Pertama)
File yang disentuh HANYA:
- `src/store/auth.ts` -> perbaiki `homePathFor()` seperti Bagian 3.
- `src/components/RequireAuth.tsx` -> timpa total seperti Bagian 3.
- `src/pages/auth/login_page.tsx` -> ganti `navigate('/')` jadi `navigate(homePathFor(user.role.name), { replace: true })`.
- `src/App.tsx` -> tulis ulang jadi daftar route di bawah (belum bikin halamannya tidak apa, bikin file kosong dulu).

`App.tsx` target akhir:
```tsx
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import LoginPage from './pages/auth/login_page'
import RegisterPage from './pages/auth/register_page'
import HomePage from './pages/home/home_page'
import RequireAuth from './components/RequireAuth'
import AdminLayout from './layouts/AdminLayout'
import DashboardLayout from './layouts/DashboardLayout'
// ... import halaman lain
export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/admin" element={<RequireAuth roles={['admin']}><AdminLayout /></RequireAuth>}>
          <Route index element={<AdminDashboardPage />} />
          <Route path="films" element={<AdminFilmsPage />} />
          <Route path="rooms" element={<AdminRoomsPage />} />
          <Route path="schedules" element={<AdminSchedulesPage />} />
          <Route path="promos" element={<AdminPromosPage />} />
        </Route>
        <Route path="/manager" element={<RequireAuth roles={['manager']}><DashboardLayout /></RequireAuth>}>
          <Route index element={<ManagerDashboardPage />} />
        </Route>
        <Route path="/superadmin" element={<RequireAuth roles={['super_admin']}><DashboardLayout /></RequireAuth>}>
          <Route index element={<SuperadminDashboardPage />} />
        </Route>
        <Route path="/staff" element={<RequireAuth roles={['staff']}><DashboardLayout /></RequireAuth>}>
          <Route index element={<StaffCashierPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}
```
Test: login tiap role harus mental ke path yang benar. Member harus mental ke `/`.

### Tahap 2 — Halaman Utama (`/` pakai `/api/v1/home`)
Buat urut:
1. `src/pages/home/home_page.tsx` (fetch seperti Bagian 4.2).
2. `src/components/Navbar.tsx` (logo + tombol Login/Logout + nama user dari `useAuth`).
3. `src/components/PromoBanner.tsx` (props: `promo`, kalau null tampilkan kotak "Belum ada promo").
4. `src/components/FilmCard.tsx` (props: `film: FilmHome`, tampilkan cover, nama, genre, harga, jadwal ringkas).
Test: matikan backend -> harus muncul pesan error, bukan layar putih. Promo null -> tidak crash.

### Tahap 3 — Admin (`/admin` + Sidebar, SEMUA hardcode)
1. `src/data/mock_admin.ts`:
   ```ts
   export const mockFilms = [{ id: 1, title: 'Contoh Film', genre: 'Action', duration: 120, price: 50000, status: 'Tayang' }]
   export const mockRooms = [{ id: 1, name: 'Studio 1', capacity: 100 }]
   export const mockSchedules = [{ id: 1, film: 'Contoh Film', room: 'Studio 1', date: '2026-10-07', time: '19:00' }]
   export const mockPromos = [{ id: 1, title: 'Promo 1', desc: 'Beli 1 gratis 1' }]
   ```
2. `src/layouts/AdminLayout.tsx` (sidebar kiri + `<Outlet/>`; menu: Dashboard, Kelola Film, Room, Schedule, Promo pakai `<Link to="/admin/films">` dst).
3. `src/components/Sidebar.tsx` (dipakai oleh AdminLayout saja).
4. 5 halaman admin, tiap halaman: judul + tabel dari mock + tombol Tambah/Edit/Hapus yang untuk sekarang hanya `alert('Mock: belum sambung API')`. DILARANG fetch API di tahap ini.
Test: login sebagai admin -> sidebar muncul, klik tiap menu ganti halaman tanpa reload.

### Tahap 4 — Manager (`/manager`, hardcode)
1. `src/data/mock_manager.ts`:
   ```ts
   export const mockRevenue = { online: 1500000, offline: 800000, total: 2300000 }
   export const mockBookings = [{ id: 1, film: 'Contoh', total: 5, type: 'online' }]
   export const mockSchedulesToday = [{ id: 1, film: 'Contoh', room: 'Studio 1', time: '19:00', booked: 10 }]
   ```
2. `src/pages/manager/manager_dashboard_page.tsx` (3 kartu angka + 2 tabel sederhana).
3. Pakai `DashboardLayout` (header sederhana + `<Outlet/>`, tanpa sidebar admin).
Test: login manager -> ke `/manager`, angka muncul.

### Tahap 5 — Superadmin (`/superadmin`, hardcode)
1. `src/data/mock_superadmin.ts`:
   ```ts
   export const mockVisitors = { today: 120, week: 800 }
   export const mockLogs = [{ id: 1, user: 'admin@x.com', action: 'Hapus film', time: '10:00' }]
   ```
2. `src/pages/superadmin/superadmin_dashboard_page.tsx` (2 kartu pengunjung + tabel log aktivitas).
Test: login super_admin -> ke `/superadmin`.

### Tahap 6 — Staff (`/staff`, hardcode)
1. `src/data/mock_staff.ts`:
   ```ts
   export const mockTickets = [{ film: 'Contoh', price: 50000, qty: 1 }]
   ```
2. `src/pages/staff/staff_cashier_page.tsx` (kiri: pilih film + jumlah tiket, kanan: total bayar + kembalian + tombol "Bayar Cash" -> `alert('Mock bayar berhasil')`).
3. Tidak ada sidebar. Tampilan besar dan sederhana untuk kasir.
Test: login staff -> ke `/staff`, hitung total benar.

### Tahap 7 — Member + Finishing
1. Member tidak punya dashboard. Setelah login langsung `navigate('/')`. Halaman utama harus tunjukkan nama member di navbar (ambil dari `useAuth`).
2. Tambah tombol Logout di semua layout: `clearAuth()` lalu `navigate('/login')`.
3. Jalankan `npm run build` di folder frontend. Wajib sukses tanpa error TS. Kalau error, perbaiki type, jangan pakai `any`.

---

## 6. Contoh Prompt Aman Untuk AI Rendah (Copy-Paste)

> "Kamu adalah helper React pemula. Baca file `frontend/panduanFront.md` Bagian 1 sampai 4 dulu. Kerjakan HANYA Tahap X. Jangan buat folder baru. Jangan pakai fetch. Hanya pakai axios dari `src/lib/api.ts` dan auth dari `src/store/auth.ts`. Tampilkan kode lengkap per file. Kalau butuh API yang tidak ada di Bagian 4, berhenti dan tanya."

Jangan beri 2 tahap sekaligus ke AI. Satu chat = satu tahap. Ini mencegah arsitektur berantakan.

---

## 7. Checklist Sebelum Dianggap Selesai

- [ ] `npm run dev` jalan, `npm run build` sukses.
- [ ] `/` ambil data asli `/api/v1/home`, ada loading + error + kosong.
- [ ] Login 5 role mental ke path benar (admin/manager/super_admin/staff/member).
- [ ] Buka URL orang lain (misal staff buka `/admin`) -> mental balik ke dashboardnya, bukan error.
- [ ] Dashboard admin/manager/superadmin/staff TIDAK ada `api.get/post` di dalamnya (cek pakai cari teks `api.` -> harus 0 hasil di folder `pages/admin, pages/manager, pages/superadmin, pages/staff`).
- [ ] Tidak ada `fetch(`, tidak ada `localStorage.getItem` di `src/` (kecuali di `store/auth.ts` bawaan zustand).
- [ ] File baru ikut penamaan `huruf_kecil_underscore.tsx`.

Cara cek cepat (di folder `frontend/`):
```powershell
Select-String -Path "src/pages/admin/*.tsx", "src/pages/manager/*.tsx", "src/pages/superadmin/*.tsx", "src/pages/staff/*.tsx" -Pattern "api\."
# harus tidak ada hasil
Select-String -Path "src/**/*.tsx", "src/**/*.ts" -Pattern "fetch\(|localStorage"
# harus tidak ada hasil (kecuali persist bawaan)
```
