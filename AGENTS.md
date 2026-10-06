# AGENTS.md — Geo Entity Manager

Panduan kerja untuk agentic AI di repo ini. Baca file ini dan `docs/PRD.md` sebelum mengubah apa pun.

## Ringkasan Proyek
Aplikasi web untuk menampilkan dan mengelola entitas berlokasi geografis (kendaraan, perangkat IoT, fasilitas) di peta. Fitur: tampil di peta, detail, tambah, ubah, hapus, dengan validasi di frontend dan backend.

## Tech Stack
- **Frontend:** React 18 + TypeScript (strict) + Vite, react-leaflet (OpenStreetMap), TanStack Query, React Hook Form + Zod, Tailwind CSS, Vitest + Testing Library
- **Backend:** Go 1.22+, Gin, go-playground/validator, pgx, golang-migrate, `log/slog`
- **Database:** PostgreSQL 16
- **Infra lokal:** Docker Compose

## Struktur Repo
```
backend/   cmd/server, internal/{handler,service,repository,model,config}, migrations/
frontend/  src/{api,components,features/entities,hooks,schemas,types}
docs/PRD.md
```

## Perintah Umum
```bash
docker compose up --build        # jalankan semua (DB + API + web)
cd backend && go test ./...      # test backend
cd backend && golangci-lint run  # lint backend
cd frontend && npm run dev       # dev server
cd frontend && npm test          # test frontend
cd frontend && npm run lint && npm run typecheck
```

## Aturan Domain (sumber kebenaran: PRD bagian 7–9)
- `type`: `vehicle | iot_device | facility | other`
- `status`: `active | inactive | maintenance | offline`
- `name` 3–100 karakter (di-trim), `description` ≤ 500, `latitude` −90..90, `longitude` −180..180
- API di `/api/v1/entities`; error selalu berformat `{"error":{"code","message","details":[{"field","message"}]}}`
- Status code: 400 JSON rusak atau UUID path tidak valid, 404 tidak ada, 413 body > 1 MB, 422 validasi gagal (termasuk field asing, salah tipe, dan query param tidak valid), 500 internal
- `GET /entities` mengembalikan `{"data":[...],"meta":{"total","page","limit"}}`; `page`/`limit` opsional (tanpa keduanya → semua data, hard cap 5000)
- PATCH = partial update (hanya field yang dikirim); `description` tidak pernah null (default `""`)
- Aturan validasi FE (Zod) dan BE harus sama; jika salah satu berubah, ubah keduanya beserta test-nya.

## Konvensi Kode
**Backend (Go)**
- Layering: handler → service → repository. Handler tidak mengakses DB langsung.
- Gunakan DTO terpisah dari model domain; semua query berparameter (tidak ada string concatenation SQL).
- Error dibungkus dengan konteks (`fmt.Errorf("...: %w", err)`); jangan bocorkan error internal ke klien.
- Pass `context.Context` sebagai argumen pertama; konfigurasi lewat env.
- Format dengan `gofmt`; test bergaya table-driven.

**Frontend (TypeScript)**
- Tanpa `any`; tipe API ada di `src/types`, skema Zod di `src/schemas`.
- Akses API hanya lewat `src/api`; data server dikelola TanStack Query (jangan duplikasi ke state lokal).
- Komponen fungsional kecil dan fokus; styling dengan Tailwind.
- Setiap tampilan data wajib menangani state loading, empty, dan error.

## Cara Kerja yang Diharapkan
1. Baca PRD dan kode terkait, lalu sampaikan rencana singkat sebelum mengubah banyak file.
2. Kerjakan per *vertical slice* kecil (mis. satu endpoint beserta test-nya), bukan semua sekaligus.
3. Tulis atau perbarui test bersamaan dengan perubahan; jalankan test, lint, dan typecheck sebelum menyatakan selesai.
4. Jangan menambah dependency tanpa alasan jelas; sebutkan alasannya.
5. Jangan mengubah scope di luar PRD; jika ada kebutuhan di luar itu, tanyakan lebih dulu.
6. Commit kecil dengan pesan jelas (Conventional Commits, mis. `feat(api): add create entity endpoint`).

## Hal yang Tidak Boleh Dilakukan
- Commit secret/credential atau file `.env` (hanya `.env.example`).
- Menonaktifkan validasi, lint, atau test agar build lolos.
- Menghapus atau menimpa migrasi yang sudah ada; buat migrasi baru.
- Mengarang fitur yang belum dikerjakan di README; fitur yang belum selesai harus dicatat jujur.

## Dokumentasi yang Harus Selalu Terbarui
`README.md` harus memuat: cara menjalankan, alasan pemilihan library, penjelasan workflow penggunaan Agentic AI (bagian mana dibantu AI dan bagian mana ditulis/di-review manual), serta daftar fitur yang belum selesai dan batasan yang diketahui.
