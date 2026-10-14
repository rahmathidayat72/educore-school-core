# Test Strategy & Quality Assurance — EduCore School

## 1. Metadata Pengujian
- **Proyek:** EduCore School — School Core Management
- **Versi:** 1.0.0
- **Target Coverage Backend:** > 85% untuk service & model logic
- **Target Frontend:** Zero build errors & Clean linting

---

## 2. Tingkat Piramida Pengujian (Test Pyramid)

```text
       ▲
      / \     E2E / Browser Verification (Manual & Automated Scenario)
     /   \    API Contract Testing (Bruno Collection)
    /     \   Integration Testing (PostgreSQL in Memory / Test DB)
   /_______\  Unit Testing (Go Native `go test` with Testify)
```

---

## 3. Komponen Pengujian

### 3.1 Backend Unit Testing
- **Framework:** Go `testing` package + `github.com/stretchr/testify` (`assert`, `mock`).
- **Cakupan Modul:**
  - `features/users`: Login, Hash password, Role middleware, Token claims.
  - `features/school_units`: CRUD unit sekolah, isolasi query.
  - `features/siswa` & `features/guru`: CRUD data master, validasi NISN/NIP unik, soft delete.
  - `features/kelas` & `features/mapel`: Pemetaan rombel, fase jenjang, konsentrasi keahlian SMK.
  - `features/jadwal`: Deteksi bentrok jadwal guru & ruang (*conflict detection matrix*).
  - `features/dashboard`: Kalkulasi metrik yayasan vs unit sekolah.
  - `features/logs`: Format log payload, masking data sensitif.

### 3.2 API Integration & Contract Testing
- **Tooling:** Koleksi Bruno (`backend/bruno/`).
- **Skenario Verifikasi:**
  - `200 OK / 201 Created`: Request valid dengan kredensial sah.
  - `400 Bad Request`: Format payload JSON tidak valid, field wajib kosong, bentrok jadwal terdeteksi.
  - `401 Unauthorized`: Token hilang, format invalid, atau token kedaluwarsa.
  - `403 Forbidden`: Admin Unit SD mencoba mengakses/mengubah data unit SMP/SMA/SMK.
  - `404 Not Found`: ID resource tidak ditemukan atau sudah di-soft-delete.

### 3.3 Frontend Verification
- **Verifikasi Build:** `npm run build` harus selesai tanpa error kompilasi Vite.
- **Routing & Navigation Guard:**
  - Sesi belum login diarahkan ke `/login`.
  - Siswa/Guru dilarang mengakses rute `/admin/units` atau `/admin/logs`.
  - Status aktif unit sekolah tersimpan di Pinia store.

---

## 4. Perintah Eksekusi Pengujian

```bash
# 1. Menjalankan seluruh unit test backend
cd backend && go test -v ./...

# 2. Menjalankan pengujian dengan kalkulasi code coverage
cd backend && go test -coverprofile=coverage.out ./...
cd backend && go tool cover -func=coverage.out

# 3. Pengujian spesifik modul jadwal (anti-bentrok)
cd backend && go test -v ./features/jadwal/...

# 4. Verifikasi kompilasi frontend Vue 3
cd frontend && npm run build
```

---

## 5. Kriteria Kelulusan Pengujian (Gate Criteria)
1. Seluruh unit test backend berstatus **PASS** (100% lulus, 0 fail).
2. Tidak ada *data race condition* pada logging asinkron (`go test -race ./...`).
3. Frontend build menghasilkan direktori `dist/` tanpa warning fatal atau error.
