# Geo Entity Manager

Aplikasi web untuk menampilkan dan mengelola *entity* berlokasi geografis (kendaraan, perangkat IoT, fasilitas) di atas peta interaktif. Fitur meliputi: tampil di peta, detail, tambah, ubah, hapus, dengan validasi di frontend dan backend. Spesifikasi lengkap ada di [docs/PRD.md](./docs/PRD.md).

## Stack

| Area | Teknologi |
| --- | --- |
| Frontend | React 19 + TypeScript + Vite, react-leaflet (OpenStreetMap), TanStack Query, React Hook Form + Zod, Tailwind CSS, Vitest |
| Backend | Go + Gin, go-playground/validator, pgx, golang-migrate, `log/slog` |
| Database | PostgreSQL 16 |
| Infra lokal | Docker Compose |

## Cara Menjalankan

Syarat: Docker & Docker Compose sudah terpasang.

1. Jalankan semua service (DB, API, Web):
   ```bash
   docker compose up --build
   ```
   Saat start, API otomatis menjalankan migrasi (termasuk data contoh) lalu terhubung ke Postgres.
2. Aplikasi Web: `http://localhost`
3. Backend API langsung: `http://localhost:8080/api/v1` (health check: `/healthz`)

### Arsitektur koneksi FE → API

Frontend memanggil API lewat path relatif **`/api/v1`** (same-origin), sehingga tidak bergantung pada CORS:

- **Docker**: Nginx di container `web` ([frontend/nginx.conf](./frontend/nginx.conf)) meneruskan `/api/` ke `api:8080`.
- **Dev lokal**: Vite dev server mem-proxy `/api` ke `http://localhost:8080` ([frontend/vite.config.ts](./frontend/vite.config.ts)).
- Bila API berada di origin lain, set `VITE_API_URL` saat build dan tambahkan origin FE ke `CORS_ORIGINS` di backend.

### Pengembangan Lokal (tanpa Docker untuk app)

```bash
docker compose up -d db                       # hanya Postgres
cd backend
$env:DATABASE_URL="postgres://geo:change-me@localhost:5432/geo_entities?sslmode=disable"  # PowerShell
go run ./cmd/server

cd frontend && npm install && npm run dev     # http://localhost:5173
```

Variabel env backend: `DATABASE_URL` (wajib), `HTTP_PORT` (default 8080), `CORS_ORIGINS`, `LOG_LEVEL`, `SHUTDOWN_TIMEOUT`. Contoh ada di [.env.example](./.env.example).

### Test, Lint, Typecheck

```bash
cd backend && go test ./... && golangci-lint run
cd frontend && npm test && npm run lint && npm run typecheck
```

## Alasan Pemilihan Library

- **Go + Gin + pgx**: Performa tinggi, efisien dalam resource, Gin mudah dikonfigurasi, dan `pgx` adalah driver native Postgres tercepat.
- **React + Vite**: Proses build sangat cepat (HMR) dan standar industri.
- **Zod + React Hook Form**: Integrasi mulus untuk validasi skema yang ketat dan efisien (menghindari render berulang).
- **TanStack Query**: Manajemen state server-side yang andal (caching, loading, error handling out of the box).
- **Leaflet (react-leaflet)**: Alternatif open-source ringan untuk Mapbox/Google Maps.
- **Vitest**: Satu konfigurasi dengan Vite, API kompatibel Jest, cepat untuk test skema/komponen.
- **Nginx (container web)**: Menyajikan build SPA sekaligus reverse proxy `/api` agar FE dan API satu origin.

## Workflow Penggunaan Agentic AI

Aplikasi ini dikembangkan menggunakan Agentic AI (Antigravity). AI memandu proses dari inisiasi proyek (S0) hingga dokumentasi (S18) melalui pendekatan iteratif per *vertical slice*:
- **Perencanaan (Manual + AI)**: Rencana di-breakdown jadi 18 slice.
- **Backend (AI)**: AI menyusun skema, migration, routing, repository, serta testing per fungsi dengan panduan konvensi Clean Architecture sederhana.
- **Frontend (AI)**: Scaffold via Vite CLI otomatis oleh agen, penulisan React Component, custom hooks (TanStack), mapping Leaflet.
- **Verifikasi**: AI menjalankan `go test`, `tsc`, dan `vitest` setelah perubahan.
- **Review & debugging manual (User)**: Pengguna menjalankan `docker compose up --build` dan menguji aplikasi di browser. Dari pengujian manual ini ditemukan bahwa data tidak tampil; AI lalu menelusuri dan memperbaiki penyebabnya:
  - `cmd/server/main.go` belum menyambungkan DB, migrasi, dan route `/api/v1/entities` (sebelumnya hanya health check aktif).
  - Request FE dari `http://localhost` diblokir CORS → diganti pola proxy same-origin (Nginx/Vite).
  - Env `PORT` di compose diganti `HTTP_PORT` sesuai config backend.
  - Skema Zod `name` belum di-trim seperti backend → diperbaiki dan dikunci dengan test.

> Catatan jujur: tabel progres di bawah sempat ditandai selesai sebelum aplikasi diuji end-to-end. Masalah di atas baru terdeteksi saat pengujian manual, sehingga verifikasi via Docker oleh manusia tetap diperlukan.

## Progres per Slice

| Slice | Isi | Status |
| --- | --- | --- |
| S0 | Housekeeping repo | ✅ |
| S1 | Kerangka backend + health check | ✅ |
| S2 | Database & migrasi | ✅ |
| S3 | Model, DTO, validasi, format error | ✅ |
| S4 | Create + Get by ID | ✅ |
| S5 | List + filter + pencarian + pagination | ✅ |
| S6 | PUT + PATCH | ✅ |
| S7 | DELETE | ✅ |
| S8 | Seed data, OpenAPI, Dockerfile backend | ✅ |
| S9 | Kerangka frontend + API client | ✅ |
| S10 | Skema Zod | ✅ |
| S11 | Peta + marker + state loading/empty/error | ✅ |
| S12 | Sidebar daftar + panel detail | ✅ |
| S13 | Form tambah + pilih lokasi di peta | ✅ |
| S14 | Edit + hapus | ✅ |
| S15 | Filter & pencarian | ✅ |
| S16 | Drag marker | ✅ |
| S17 | Docker frontend, responsif, aksesibilitas | ✅ |
| S18 | Dokumentasi akhir & verifikasi *fresh clone* | ✅ |

## Fitur Belum Selesai & Batasan yang Diketahui

- Test **Frontend** baru mencakup skema Zod (`src/schemas/entity.test.ts`); test komponen (Map, Sidebar, EntityForm) dengan Testing Library belum ditulis.
- Test integrasi repository Postgres hanya berjalan bila `TEST_DATABASE_URL` di-set (di-skip jika tidak).
- Data contoh dimasukkan lewat migrasi `000002_seed_data`, sehingga selalu ikut ter-apply; env `SEED_DATA` saat ini belum berpengaruh.
- Pagination di sisi Frontend belum memakai tombol page/infinite scroll (mengambil semua data, hard cap 5000 dari backend).
- `npm run lint` masih memberi 2 warning React (`setState` di dalam effect) yang belum di-refactor.
- Di luar scope sesuai PRD: autentikasi, pelacakan real-time, riwayat lokasi, filter bounding box, dark mode.
