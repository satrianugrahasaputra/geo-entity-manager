# Geo Entity Manager

Aplikasi web untuk menampilkan dan mengelola *entity* berlokasi geografis (kendaraan, perangkat IoT, fasilitas) di atas peta interaktif. Fitur meliputi: tampil di peta, detail, tambah, ubah, hapus, dengan validasi di frontend dan backend. Spesifikasi lengkap ada di [docs/PRD.md](./docs/PRD.md).

## Stack

| Area | Teknologi |
| --- | --- |
| Frontend | React 18 + TypeScript + Vite, react-leaflet (OpenStreetMap), TanStack Query, React Hook Form + Zod, Tailwind CSS |
| Backend | Go + Gin, go-playground/validator, pgx, golang-migrate, `log/slog` |
| Database | PostgreSQL 16 |
| Infra lokal | Docker Compose |

## Cara Menjalankan

Syarat: Docker & Docker Compose sudah terpasang.

1. Jalankan semua service (DB, API, Web):
   ```bash
   docker-compose up --build
   ```
2. Aplikasi Web dapat diakses di: `http://localhost:80` (atau port yang dipetakan oleh Docker).
3. Backend API dapat diakses di: `http://localhost:8080/api/v1`

**Untuk Pengembangan Lokal (tanpa Docker):**
- **Backend**: `cd backend && go run ./cmd/server` (pastikan env `DATABASE_URL` sudah mengarah ke Postgres aktif)
- **Frontend**: `cd frontend && npm install && npm run dev`

## Alasan Pemilihan Library

- **Go + Gin + pgx**: Performa tinggi, efisien dalam resource, Gin mudah dikonfigurasi, dan `pgx` adalah driver native Postgres tercepat.
- **React + Vite**: Proses build sangat cepat (HMR) dan standar industri.
- **Zod + React Hook Form**: Integrasi mulus untuk validasi skema yang ketat dan efisien (menghindari render berulang).
- **TanStack Query**: Manajemen state server-side yang andal (caching, loading, error handling out of the box).
- **Leaflet (react-leaflet)**: Alternatif open-source ringan untuk Mapbox/Google Maps.

## Workflow Penggunaan Agentic AI

Aplikasi ini dikembangkan menggunakan Agentic AI (Antigravity). AI memandu proses dari inisiasi proyek (S0) hingga dokumentasi (S18) melalui pendekatan iteratif per *vertical slice*:
- **Perencanaan (Manual + AI)**: Rencana di-breakdown jadi 18 slice.
- **Backend (AI)**: AI menyusun skema, migration, routing, repository, serta testing per fungsi dengan panduan konvensi Clean Architecture sederhana.
- **Frontend (AI)**: Scaffold via Vite CLI otomatis oleh agen, penulisan React Component, custom hooks (TanStack), mapping Leaflet.
- **Verifikasi**: AI otomatis menjalankan `go test` dan `tsc` sebelum commit untuk memastikan *type-safety* dan lolos tes.

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

- Tes otomatis (Unit/Integration test) untuk **Frontend** dengan Vitest belum ditulis.
- Pagination di sisi Frontend belum sepenuhnya menggunakan tombol page/infinite scroll (hanya menampilkan data berdasarkan limit dari backend yang dikonfigurasi ke 5000 max).
- Di luar scope sesuai PRD: autentikasi, pelacakan real-time, riwayat lokasi, filter bounding box, dark mode.
