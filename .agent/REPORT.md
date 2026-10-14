# Execution Report & Evidence Log — EduCore School

## 1. Summary Laporan Eksekusi
- **Proyek:** EduCore School — School Core Management
- **Terakhir Diperbarui:** 14 Oktober 2026
- **Status Sesi Terkini:** Menginisialisasi dokumen tata kelola proyek (Phase 01).

---

## 2. Log Eksekusi Task

### [2026-10-14] Task 01.1 & Task 01.2 — Inisialisasi Struktur Tata Kelola Dokumen & AGENTS.md
- **Status:** `✅ Selesai dan terverifikasi`
- **Tindakan:**
  1. Membuat struktur folder `docs/project/`:
     - `PROJECT-BRIEF.md`: Dokumen ringkasan proyek EduCore School.
     - `BRD.md`: Dokumen kebutuhan bisnis multi-unit yayasan.
     - `PRD.md`: Dokumen spesifikasi produk, persona, role matrix, FR & NFR.
     - `ARCHITECTURE.md`: Arsitektur Clean Feature Go, tenant isolation, dan deteksi bentrok jadwal.
     - `SOP.md`: Standar operasional prosedur pengembangan stack Go + Vue 3.
     - `SECURITY.md`: Penilaian risiko MEDIUM dan matriks 12 ancaman keamanan.
     - `TEST-STRATEGY.md`: Strategi piramida pengujian, target coverage >85%, dan perintah uji.
     - `TRACEABILITY.md`: Matriks keterlacakan Objective → Business Req → Functional Req → Task.
     - `IMPLEMENTATION-PLAN.md`: Roadmap eksekusi Phase 01 hingga Phase 07.
  2. Membuat struktur folder `.agent/`:
     - `PROJECT-STATE.md`: Status lifecycle, metadata gerbang, dan risk tier.
     - `TASKS.md`: Task aktif dan backlog roadmap terperinci.
     - `REPORT.md`: Log bukti verifikasi dan laporan eksekusi berkala.
     - `DECISIONS.md`: Catatan keputusan arsitektur dan bisnis.
     - `OPEN-QUESTIONS.md`: Daftar pertanyaan terbuka.
     - `CHANGE-REQUESTS.md`: Log permintaan perubahan.
  3. Menyalin dan memastikan ketersediaan `AGENTS.md` di root proyek `School-Core-Management/AGENTS.md`.
- **Bukti / Verifikasi:**
  - 9 file dokumen di `School-Core-Management/docs/project/` terverifikasi via `glob`.
  - 6 file tracking di `School-Core-Management/.agent/` terverifikasi via `glob`.
  - File `School-Core-Management/AGENTS.md` terverifikasi via `glob`.

### [2026-10-14] Task 01.3 — Pembersihan Artefak Lokal & Update .gitignore
- **Status:** `✅ Selesai dan terverifikasi`
- **Tindakan:**
  1. Membuat file `.gitignore` di root workspace (`/home/rahmath/Documents/rahmath/app_sekolah/.gitignore`).
  2. Membuat file `.gitignore` di `School-Core-Management/.gitignore`.
  3. Memperbarui `School-Core-Management/backend/.gitignore` untuk mengabaikan binary dan log.
  4. Menghapus file binary Linux lama `School-Core-Management/backend/go_rest_native_sekolah`.
  5. Menghapus log runtime lama `School-Core-Management/backend/server.log`.
- **Bukti / Verifikasi:**
  - `glob` untuk `server.log` dan `go_rest_native_sekolah` mengembalikan `No files found`.
  - `git status` pada backend memverifikasi file binary dan log tidak lagi terlacak.

### [2026-10-14] Task 01.4 — Setup Kontrol Versi Git Monorepo Terpadu
- **Status:** `✅ Selesai dan terverifikasi`
- **Tindakan:**
  1. Mengamankan git lama backend dengan memindahkan `backend/.git` ke `backend/.git.backup`.
  2. Mengamankan git sub-repo `agent/.git` ke `agent/.git.backup`.
  3. Menginisialisasi Git monorepo baru di root `School-Core-Management/` pada branch `main`.
  4. Mengonfigurasi identitas author `Rahmat Hidayat` secara lokal di repository.
  5. Melakukan initial baseline commit: `feat(core): initialize EduCore School Core Management monorepo` (175 files, 33,747 insertions).
- **Bukti / Verifikasi:**
  - `git log -1` memverifikasi commit hash `b2aeff5` berhasil tercatat dengan aman di branch `main`.
  - Seluruh source code Go backend, Vue 3 frontend, docs/project, dan .agent tracking kini berada dalam satu riwayat sinkron.
  - Remote repository GitHub baru dibuat: `https://github.com/rahmathidayat72/educore-school-core`
  - Push branch `main` ke `origin/main` berhasil diverifikasi.

### [2026-10-14] Task 02.1 — Pembuatan Skrip Migrasi SQL Tabel school_units
- **Status:** `✅ Selesai dan terverifikasi`
- **Tindakan:**
  1. Membuat direktori `School-Core-Management/backend/migrations/`.
  2. Menulis skrip migrasi UP `000001_create_school_units_table.up.sql` dengan skema kolom id (UUID/VARCHAR 36), nama_sekolah, jenjang (SD/SMP/SMA/SMK), npsn unik, kontak, alamat, status, dan soft delete.
  3. Menambahkan partial unique index pada `npsn` (WHERE `delete_at IS NULL`) dan indeks performa pada `jenjang`, `status`, dan `delete_at`.
  4. Menulis skrip migrasi DOWN `000001_create_school_units_table.down.sql` untuk rollback aman (`DROP TABLE IF EXISTS school_units CASCADE`).
- **Bukti / Verifikasi:**
  - File migrasi UP dan DOWN terverifikasi ada via `glob`.
  - Sintaks DDL PostgreSQL telah divalidasi.

### [2026-10-14] Task 02.2 — Penambahan Kolom school_unit_id pada Entitas Operasional
- **Status:** `✅ Selesai dan terverifikasi`
- **Tindakan:**
  1. Menulis skrip migrasi UP `000002_add_school_unit_id_to_entities.up.sql` yang menambahkan kolom `school_unit_id VARCHAR(36)` dan B-Tree index pada 7 entitas utama: `users`, `guru`, `kelas`, `siswa`, `mata_pelajaran`, `jadwal_pelajaran`, dan `transaction_logs`.
  2. Menulis skrip migrasi DOWN `000002_add_school_unit_id_to_entities.down.sql` untuk rollback aman drop column & drop index.
- **Bukti / Verifikasi:**
  - File migrasi 000002 UP dan DOWN terverifikasi ada via `glob`.
  - Kolom dibuat bersifat nullable terlebih dahulu agar kompatibel dengan data yang sudah ada sebelum proses backfill.

### [2026-10-14] Task 02.3 — Penyiapan Unit Percontohan & Data Backfill
- **Status:** `✅ Selesai dan terverifikasi`
- **Tindakan:**
  1. Menulis skrip migrasi UP `000003_seed_default_units_and_backfill.up.sql` yang mendaftarkan 4 unit sekolah yayasan (SMA, SMP, SMK, SD).
  2. Melakukan backfill data yang sudah ada agar terhubung ke unit pilot utama (`unit-sma-001`).
  3. Memasang Foreign Key constraint dari tabel operasional (`users`, `guru`, `kelas`, `siswa`, `mata_pelajaran`, `jadwal_pelajaran`, `transaction_logs`) ke `school_units(id)`.
  4. Menulis skrip migrasi DOWN `000003_seed_default_units_and_backfill.down.sql` untuk pemulihan/rollback aman jika dibutuhkan.
- **Bukti / Verifikasi:**
  - File migrasi 000003 UP dan DOWN terverifikasi ada via `glob`.
  - Sintaks FK constraint dan data seed divalidasi.

### [2026-10-14] Task 03.1 — Pembaruan Entitas Users & JWT Claims Multi-Unit
- **Status:** `✅ Selesai dan terverifikasi`
- **Tindakan:**
  1. Menambahkan field `SchoolUnitID *string` pada entitas `UserCore` di `features/auth/entities.go` dan `features/users/entities.go`.
  2. Memperbarui `MetaToken` dan `VerifyTokenHeader` pada `helper/middlewares.go` untuk menyalin dan membaca `school_unit_id` dari JWT claims.
  3. Mengizinkan role `super_admin` pada `AdminOnlyMiddleware` dan `UserSelfOrAdminMiddleware`.
  4. Menyesuaikan query login di `features/auth/model/query.go` untuk men-select `school_unit_id` dan memvalidasi `delete_at IS NULL`.
  5. Menyesuaikan `features/auth/controllers/controller.go` agar klaim `school_unit_id` disisipkan ke token dan response JSON.
  6. Menyesuaikan query CRUD di `features/users/model/querry.go` agar membaca dan menyimpan `school_unit_id`.
  7. Menjalankan seluruh rangkaian unit testing backend Go.
- **Bukti / Verifikasi:**
  - `go test ./...` lulus 100% tanpa error atau regresi.
  - Test suites `features/auth`, `features/users`, dan `helper` seluruhnya berstatus `PASS`.

### [2026-10-14] Task 03.2 — Isolasi Data Unit Sekolah di Level Query & Logic Backend
- **Status:** `✅ Selesai dan terverifikasi`
- **Requirement:** BUS-003, FR-001, FR-004, NFR-SEC-001
- **Tindakan:**
  1. Menambahkan helper `ExtractUserContext(r)` pada `helper/middlewares.go` untuk membaca metadata user, role, unit ID dari JWT claims serta query param `?unit_id=` bagi Super Admin Yayasan.
  2. Mengimplementasikan isolasi data multi-unit pada modul `features/guru` (`entities.go`, `model/model.go`, `model/query.go`, `service/logic.go`, `controllers/controller.go`).
  3. Mengimplementasikan isolasi data multi-unit pada modul `features/siswa` (`entities.go`, `model/model.go`, `model/query.go`, `service/logic.go`, `controllers/controllers.go`).
  4. Mengimplementasikan isolasi data multi-unit pada modul `features/kelas` (`entities.go`, `model/model.go`, `model/query.go`, `service/logic.go`, `controllers/controller.go`).
  5. Mengimplementasikan isolasi data multi-unit pada modul `features/mata_pelajaran` (`entities.go`, `model/model.go`, `model/query.go`, `service/logic.go`, `controllers/controller.go`, `controllers/formatter.go`).
  6. Mengimplementasikan isolasi data multi-unit pada modul `features/jadwal` (`entities.go`, `model/query.go`, `service/logic.go`, `controllers/controller.go`) termasuk pembatasan deteksi bentrok jadwal (`CheckConflict`) hanya di dalam unit sekolah yang sama.
  7. Memperbarui mock unit test yang terdampak (`guru`, `kelas`, `mata_pelajaran`, `siswa`) untuk mendukung signature variadic `schoolUnitID ...string`.
- **Bukti / Verifikasi:**
  - `go test -v ./...` berhasil dijalankan pada root backend: seluruh paket dan mock test berstatus `PASS` (100% lulus).
  - `go build ./...` berhasil tanpa peringatan maupun error kompilasi.

### [2026-10-14] Task 03.3 — Implementasi Endpoint CRUD Unit Sekolah Khusus Super Admin
- **Status:** `✅ Selesai dan terverifikasi`
- **Requirement:** BUS-003, FR-001, FR-004, NFR-SEC-001
- **Tindakan:**
  1. Menambahkan middleware otorisasi `SuperAdminOnlyMiddleware` pada `helper/middlewares.go` untuk melindungi operasi perubahan data yayasan.
  2. Membuat entity domain `SchoolUnitCore` dan interface `DataSchoolUnitInterface` serta `ServiceSchoolUnitInterface` di `features/school_units/entities.go`.
  3. Mengimplementasikan model database PostgreSQL `SchoolUnitModel` di `features/school_units/model/` dengan query parameterized (`SelectAll`, `SelectByID`, `Insert`, `Update`, `Delete`, `CheckNPSNExists`).
  4. Mengimplementasikan service layer `SchoolUnitService` di `features/school_units/service/logic.go` dengan validasi jenjang (SD, SMP, SMA, SMK), validasi status, dan pengecekan keunikan NPSN.
  5. Menulis unit testing komprehensif `features/school_units/service/logic_test.go` dengan mock repository.
  6. Membuat DTO dan response serializer di `features/school_units/controllers/formatter.go`.
  7. Membuat HTTP presentation layer `SchoolUnitController` di `features/school_units/controllers/controller.go`.
  8. Mendaftarkan route HTTP `/school-units`, `/school-units/detail`, `/school-units/tambah`, `/school-units/update`, dan `/school-units/deleted` pada `router/router.go`.
- **Bukti / Verifikasi:**
  - `go test -v ./features/school_units/...` berstatus `PASS` (100% lulus seluruh skenario pengujian: get all, get by id, validation error, duplicate NPSN, update, dan delete).
  - `go test -v ./...` berstatus `PASS` (100% lulus seluruh paket backend tanpa regresi).
  - `go build ./...` sukses tanpa error maupun peringatan kompilasi.

### [2026-10-14] Task 04.1 — Penyesuaian Tabel Rombel/Kelas untuk Multi-Jenjang (SD-SMK)
- **Status:** `✅ Selesai dan terverifikasi`
- **Requirement:** OBJ-001, BUS-003, FR-005, NFR-SEC-001
- **Tindakan:**
  1. Menulis skrip migrasi UP `000004_enhance_kelas_multi_jenjang.up.sql` untuk menambahkan indeks performa pada kolom `tingkat`, `fase`, `jurusan`, dan `(school_unit_id, tingkat)`.
  2. Menulis skrip migrasi DOWN `000004_enhance_kelas_multi_jenjang.down.sql` untuk rollback aman.
  3. Mengimplementasikan fungsi inferensi dan validasi Kurikulum Merdeka `ResolveFaseDanJurusan` di `features/kelas/service/logic.go`:
     - Normalisasi tingkat (Romawi/Arab ke angka standar 1–13).
     - Auto-assignment fase jika kosong: Tingkat 1–2 $\rightarrow$ A, 3–4 $\rightarrow$ B, 5–6 $\rightarrow$ C, 7–9 $\rightarrow$ D, 10 $\rightarrow$ E, 11–13 $\rightarrow$ F.
     - Validasi konsistensi jika fase diinput manual (menolak jika bertentangan dengan tingkat).
     - Validasi ketat SMK: wajib menyertakan atribut jurusan / program keahlian.
     - Sanitasi SD & SMP: otomatis mengosongkan atribut jurusan.
  4. Menambahkan dukungan filter query parameter `?tingkat=`, `?fase=`, dan `?jurusan=` pada controller handler `Kelas` di `features/kelas/controllers/controller.go`.
  5. Menulis unit testing `TestResolveFaseDanJurusan` di `features/kelas/service/logic_test.go` yang menguji skenario auto-assign SD/SMP/SMA, penolakan konflik fase, validasi jurusan SMK, dan sanitasi jenjang dasar.
- **Bukti / Verifikasi:**
  - `go test -v ./features/kelas/...` berstatus `PASS` (seluruh 7 skenario pengujian baru dan test eksisting lulus 100%).
  - `go test -v ./...` berstatus `PASS` (zero regression pada seluruh modul backend).
  - `go build ./...` sukses tanpa peringatan atau error kompilasi.
