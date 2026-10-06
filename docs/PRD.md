# Product Requirements Document — Geo Entity Manager

|  |  |
| --- | --- |
| **Versi** | 1.0 (Draft untuk implementasi) |
| **Tanggal** | 6 Oktober 2026 |
| **Penulis** | \[Nama Kandidat\] |
| **Konteks** | Take-Home Test — Software Developer |
| **Stack** | React + TypeScript · Go · PostgreSQL |

---

## 1. Ringkasan Eksekutif

**Geo Entity Manager** adalah aplikasi web untuk menampilkan dan mengelola *entity* yang memiliki lokasi geografis (kendaraan, perangkat IoT, fasilitas, dll.) di atas peta interaktif. Pengguna dapat melihat sebaran entitas, membuka detailnya, serta menambah, mengubah, dan menghapusnya. Validasi dilakukan di frontend dan backend.

**Tujuan dokumen:** menjadi acuan tunggal scope, perilaku sistem, desain teknis, dan kriteria selesai agar implementasi dalam waktu 2 hari tetap terarah.

## 2. Latar Belakang & Masalah

Organisasi yang mengelola aset bergerak maupun tetap membutuhkan satu tampilan spasial untuk memantau posisi dan status aset. Data tersebar dalam daftar/tabel tanpa konteks lokasi sehingga sulit dipantau dan diperbarui.

## 3. Tujuan & Non-Tujuan

**Tujuan**

- G1. Visualisasi seluruh entitas pada peta dengan status yang mudah dibedakan.
- G2. CRUD entitas lengkap (create, read, update, delete) dari UI.
- G3. Validasi konsisten di FE dan BE; BE adalah *source of truth*.
- G4. Kode bersih, teruji, mudah dijalankan (satu perintah), terdokumentasi.

**Non-Tujuan (out of scope)**

- Autentikasi/otorisasi multi-user dan RBAC.
- Pelacakan real-time (WebSocket/MQTT) dan riwayat pergerakan.
- Geofencing, routing, clustering skala jutaan titik.
- Aplikasi mobile native.

## 4. Persona & Use Case

| Persona | Kebutuhan |
| --- | --- |
| **Operator** | Melihat semua entitas di peta, memfilter berdasarkan tipe/status, membuka detail. |
| **Administrator Data** | Menambah entitas baru, memperbarui atribut/lokasi, menghapus entitas usang. |
| **Reviewer/Penguji** | Menjalankan aplikasi cepat dan menilai kualitas kode, validasi, dan dokumentasi. |

**Alur utama:** buka aplikasi → peta menampilkan marker → klik marker untuk detail → *Edit* / *Delete*; atau klik *Add Entity* → klik peta/isi koordinat → simpan.

## 5. Scope & Prioritas (MoSCoW)

| Prioritas | Fitur |
| --- | --- |
| **Must** | Tampil di peta, detail, tambah, ubah, hapus, validasi FE+BE, persistensi DB, dokumentasi |
| **Should** | Filter tipe/status, pencarian nama, pilih lokasi via klik peta, drag marker, toast & konfirmasi hapus, Docker Compose, unit test |
| **Could** | Pagination, pencarian bounding box, Swagger/OpenAPI, seed data, dark mode |
| **Won't** | Auth, real-time, riwayat lokasi |

## 6. Functional Requirements

| ID | Requirement | Prio | Acceptance Criteria |
| --- | --- | --- | --- |
| FR-01 | Menampilkan seluruh entitas sebagai marker pada peta | Must | Given data tersimpan, when halaman dibuka, then setiap entitas muncul di koordinatnya; marker berwarna sesuai status |
| FR-02 | Melihat detail entitas | Must | When marker/list item diklik, then panel/popup menampilkan nama, tipe, status, deskripsi, koordinat, waktu dibuat/diubah |
| FR-03 | Menambah entitas | Must | When form valid disubmit, then entitas tersimpan di DB dan marker langsung muncul tanpa reload |
| FR-04 | Mengubah entitas | Must | When form edit disubmit, then perubahan tersimpan dan marker/detail ter-update |
| FR-05 | Menghapus entitas | Must | When user mengonfirmasi, then data terhapus dari DB dan marker hilang; jika ID tidak ada, tampil pesan error |
| FR-06 | Validasi input FE | Must | Field tidak valid ditandai inline, tombol simpan tidak mengirim request |
| FR-07 | Validasi data BE | Must | Request tidak valid ditolak 4xx dengan detail error per-field; data tidak tersimpan |
| FR-08 | Pilih lokasi lewat klik peta | Should | When mode tambah aktif dan peta diklik, then lat/lng terisi otomatis dan marker sementara tampil |
| FR-09 | Drag marker untuk update lokasi | Should | When marker dilepas di titik baru, then lokasi diperbarui via API; rollback jika gagal |
| FR-10 | Filter & pencarian | Should | Filter tipe/status dan pencarian nama memperbarui marker dan daftar |
| FR-11 | Daftar entitas (sidebar) | Should | Klik item memfokuskan peta ke entitas (fly-to) |
| FR-12 | Penanganan state loading/empty/error | Must | Skeleton saat memuat, pesan saat kosong, pesan error ramah + opsi coba lagi |

## 7. Model Data

**Entity**

| Field | Tipe | Aturan |
| --- | --- | --- |
| `id` | UUID | PK, di-generate server |
| `name` | string | wajib, 3–100 karakter, di-trim |
| `type` | enum | `vehicle`, `iot_device`, `facility`, `other` |
| `status` | enum | `active`, `inactive`, `maintenance`, `offline` |
| `description` | string | opsional, maks 500 karakter |
| `latitude` | float64 | wajib, −90 s.d. 90 |
| `longitude` | float64 | wajib, −180 s.d. 180 |
| `created_at` / `updated_at` | timestamptz | otomatis (UTC) |

Constraint DB: `NOT NULL`, `CHECK` untuk rentang lat/lng, `CHECK` enum, indeks pada `(type, status)` dan `(latitude, longitude)`. Opsi pengembangan: PostGIS bila dibutuhkan query spasial lanjutan.

## 8. Aturan Validasi (FE & BE)

| Field | Aturan | FE (Zod) | BE (validator) |
| --- | --- | --- | --- |
| name | wajib, 3–100, bukan hanya spasi | ✔ | ✔ |
| type / status | harus anggota enum | ✔ | ✔ |
| description | ≤ 500 karakter | ✔ | ✔ |
| latitude / longitude | angka valid dalam rentang; bukan NaN/Inf | ✔ | ✔ |
| id (path) | format UUID valid | – | ✔ |
| body | JSON valid, ukuran maks 1 MB, field asing ditolak | – | ✔ |

> Prinsip: FE memberi umpan balik cepat; BE tetap memvalidasi ulang seluruh input dan tidak pernah mempercayai klien.

## 9. Desain API (REST, JSON, prefix `/api/v1`)

| Method | Endpoint | Deskripsi | Sukses |
| --- | --- | --- | --- |
| GET | `/entities` | Daftar; query: `type`, `status`, `q`, `bbox`, `page`, `limit` | 200 |
| GET | `/entities/{id}` | Detail | 200 / 404 |
| POST | `/entities` | Buat entitas | 201 + `Location` |
| PUT | `/entities/{id}` | Ganti seluruh atribut | 200 / 404 |
| PATCH | `/entities/{id}` | Ubah sebagian (mis. lokasi/status) | 200 / 404 |
| DELETE | `/entities/{id}` | Hapus | 204 / 404 |
| GET | `/healthz` | Health check | 200 |

**Contoh request `POST /entities`**

```json
{
  "name": "Truk Logistik 01",
  "type": "vehicle",
  "status": "active",
  "description": "Rute Semarang–Ungaran",
  "latitude": -7.1396,
  "longitude": 110.4203
}
```

**Format error seragam**

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Input tidak valid",
    "details": [{ "field": "latitude", "message": "harus antara -90 dan 90" }]
  }
}
```

Kode status: `400` JSON rusak · `404` tidak ditemukan · `422` validasi gagal · `500` error internal (tanpa membocorkan detail).

## 10. Arsitektur & Pilihan Teknologi

```
Browser (React + TS + Leaflet)  ──HTTP/JSON──▶  Go REST API  ──SQL──▶  PostgreSQL
```

Backend berlapis: **handler → service → repository**, dengan DTO terpisah dari model domain.

| Area | Pilihan | Alasan |
| --- | --- | --- |
| Frontend | React 18 + TypeScript + Vite | Wajib oleh soal; Vite cepat dan ringan |
| Peta | Leaflet + react-leaflet + OpenStreetMap | Gratis, tanpa API key sehingga reviewer langsung bisa menjalankan; ekosistem matang |
| Server state | TanStack Query | Cache, invalidasi, loading/error state otomatis |
| UI state | React state / Zustand (bila perlu) | Ringan, tanpa boilerplate |
| Form & validasi | React Hook Form + Zod | Skema tipe-aman, performa baik |
| Styling | Tailwind CSS | Cepat dan konsisten |
| Backend | Go + Gin | Populer, dokumentasi luas, binding & validasi bawaan |
| Validasi BE | go-playground/validator | Standar de facto, tag deklaratif |
| Database | PostgreSQL + pgx | Handal, mendukung CHECK constraint, siap PostGIS |
| Migrasi | golang-migrate | Skema ter-versi dan dapat diulang |
| Testing | Go `testing` + testify; Vitest + Testing Library | Cakupan unit dan komponen |
| Deploy lokal | Docker Compose | Satu perintah menjalankan DB + API + web |

## 11. Spesifikasi UI/UX

- **Layout:** sidebar kiri (pencarian, filter, daftar) dan peta di area utama; header dengan tombol *Add Entity*.
- **Marker:** warna per status — hijau (active), abu-abu (inactive), kuning (maintenance), merah (offline); ikon sesuai tipe.
- **Detail:** panel kanan/popup berisi atribut lengkap, tombol *Edit* dan *Delete*.
- **Form:** modal/drawer untuk tambah & edit, pesan error inline, tombol nonaktif saat *submitting*.
- **Hapus:** dialog konfirmasi.
- **Umpan balik:** toast sukses/gagal, skeleton saat memuat.
- **Aksesibilitas & responsif:** label form, fokus keyboard, tata letak sidebar menjadi *drawer* pada layar kecil.

## 12. Non-Functional Requirements

| Kategori | Target |
| --- | --- |
| Performa | Daftar 1.000 entitas dimuat \< 1 detik (lokal); respons API p95 \< 200 ms |
| Keamanan | Query berparameter (anti SQL injection), CORS terbatas ke origin FE, batas ukuran body, pesan error aman, konfigurasi via env, tanpa secret di repo |
| Keandalan | Graceful shutdown, health check, migrasi otomatis saat start |
| Maintainability | Struktur modular, lint (golangci-lint, ESLint), format (gofmt, Prettier), commit terstruktur |
| Observability | Logging terstruktur (`slog`) dengan request ID |
| Portabilitas | Berjalan di Linux/macOS/Windows via Docker |

## 13. Struktur Repository

```
.
├── backend/        # cmd/server, internal/{handler,service,repository,model,config}, migrations/
├── frontend/       # src/{api,components,features/entities,hooks,schemas,types}
├── docker-compose.yml
├── .env.example
├── README.md       # cara menjalankan, pilihan library, workflow AI, fitur belum selesai
├── CLAUDE.md       # konteks & aturan kerja untuk agentic AI
└── docs/PRD.md
```

## 14. Strategi Pengujian

- **Backend:** unit test validasi & service; test handler dengan `httptest`; test repository terhadap DB uji.
- **Frontend:** test skema Zod, komponen form (validasi, submit), dan hook data dengan mock API.
- **Manual E2E:** skenario CRUD lengkap + kasus tepi (koordinat di luar rentang, nama kosong, hapus ID tidak ada, API mati).

## 15. Rencana Kerja (2 Hari)

| Hari | Fokus | Output |
| --- | --- | --- |
| **Hari 1** | Setup repo, Docker Compose, migrasi DB; API CRUD + validasi + error handling; unit test BE; `CLAUDE.md` | API berfungsi dan teruji (Postman/cURL) |
| **Hari 2** | Setup FE, peta + marker, detail, form tambah/edit/hapus, validasi Zod, filter; integrasi; test FE; polishing; README | Aplikasi end-to-end siap kirim |

Cadangan waktu ±2 jam di akhir untuk perbaikan bug dan verifikasi *fresh clone*. Bila waktu terbatas, kerjakan seluruh **Must** terlebih dahulu, lalu **Should**.

## 16. Risiko & Mitigasi

| Risiko | Dampak | Mitigasi |
| --- | --- | --- |
| Waktu hanya 2 hari | Fitur tidak selesai | Prioritas MoSCoW; kerjakan vertical slice; catat yang belum selesai di README |
| Kompleksitas integrasi Leaflet dengan React | Marker/ikon bermasalah | Gunakan react-leaflet; atur ikon marker sejak awal |
| Perbedaan aturan validasi FE vs BE | Perilaku tidak konsisten | Satu tabel aturan (bagian 8) sebagai acuan |
| Reviewer gagal menjalankan aplikasi | Penilaian buruk | Docker Compose, `.env.example`, uji *fresh clone* |
| Kode hasil AI tidak dipahami | Sulit dipertanggungjawabkan saat wawancara | Review manual setiap bagian, tulis test sendiri, pahami setiap keputusan |

## 17. Rencana Dokumentasi & Pengiriman

**README.md wajib memuat:**

1. Deskripsi singkat, fitur, dan tangkapan layar.
2. **Cara menjalankan** (prasyarat, `cp .env.example .env`, `docker compose up --build`, URL akses, cara menjalankan test, opsi tanpa Docker).
3. **Alasan pemilihan library** (ringkas dari bagian 10).
4. **Workflow penggunaan Agentic AI** — jelaskan jujur: tool yang dipakai, tahapan yang dibantu AI (perencanaan, scaffolding, implementasi, test, dokumentasi), bagian yang ditulis/di-review manual, dan cara verifikasi hasilnya.
5. **Fitur belum selesai & batasan** (wajib bila ada).
6. Keputusan desain dan pengembangan lanjutan.

**Checklist pengiriman**

- [ ] Repo Git rapi dengan riwayat commit bermakna
- [ ] `CLAUDE.md`/`AGENT.md` disertakan
- [ ] README lengkap dan sudah diuji dari *fresh clone*
- [ ] Tidak ada secret/credential di repo
- [ ] Email dikirim ke kedua alamat penerima sebelum tenggat, berisi tautan repo

## 18. Definition of Done

- Seluruh requirement **Must** lolos acceptance criteria.
- Validasi bekerja di FE dan BE dengan pesan jelas.
- Aplikasi berjalan dengan satu perintah dari *fresh clone*.
- Test utama lulus; lint dan format bersih.
- Dokumentasi lengkap, termasuk fitur yang belum selesai.

## 19. Asumsi & Pertanyaan Terbuka

- **Asumsi:** satu pengguna tanpa login; semua entitas ditampilkan bersamaan (data \< beberapa ribu); koordinat berformat WGS84 desimal; waktu disimpan dalam UTC.
- **Terbuka:** apakah penguji mengharapkan atribut tambahan (mis. *metadata* kustom per tipe)? Bila ya, tambahkan kolom `attributes` bertipe JSONB sebagai pengembangan lanjutan.