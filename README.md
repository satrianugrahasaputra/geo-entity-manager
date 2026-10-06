# Geo Entity Manager

Aplikasi web untuk menampilkan dan mengelola *entity* berlokasi geografis (kendaraan, perangkat IoT, fasilitas) di atas peta interaktif. Fitur yang direncanakan: tampil di peta, detail, tambah, ubah, hapus, dengan validasi di frontend dan backend. Spesifikasi lengkap ada di [docs/PRD.md](./docs/PRD.md).

> **Status: dalam pengerjaan.** README ini diperbarui di setiap slice. Bagian yang belum terisi berarti fiturnya memang belum dikerjakan.

## Stack

| Area | Teknologi |
| --- | --- |
| Frontend | React 18 + TypeScript + Vite, react-leaflet (OpenStreetMap), TanStack Query, React Hook Form + Zod, Tailwind CSS |
| Backend | Go + Gin, go-playground/validator, pgx, golang-migrate, `log/slog` |
| Database | PostgreSQL 16 |
| Infra lokal | Docker Compose |

## Cara Menjalankan

_Belum tersedia. Akan diisi setelah backend dan Docker Compose selesai._

## Alasan Pemilihan Library

_Akan diisi (ringkasan PRD bagian 10)._

## Workflow Penggunaan Agentic AI

_Akan diisi secara jujur di akhir pengerjaan: tool yang dipakai, bagian yang dibantu AI, bagian yang ditulis/di-review manual, dan cara verifikasinya._

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
| S15 | Filter & pencarian | ⏳ |
| S16 | Drag marker | ⏳ |
| S17 | Docker frontend, responsif, aksesibilitas | ⏳ |
| S18 | Dokumentasi akhir & verifikasi *fresh clone* | ⏳ |

## Fitur Belum Selesai & Batasan yang Diketahui

- Semua fitur aplikasi masih dalam pengerjaan (lihat tabel progres).
- Di luar scope sesuai PRD: autentikasi, pelacakan real-time, riwayat lokasi, filter bounding box, pagination di UI, dark mode.
