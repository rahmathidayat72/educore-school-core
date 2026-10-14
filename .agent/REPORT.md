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

### [2026-10-14] Task 04.2 — Pembuatan Tabel academic_years dan class_enrollments
- **Status:** `✅ Selesai dan terverifikasi`
- **Requirement:** OBJ-001, BUS-004, FR-004, FR-006, NFR-SEC-001
- **Tindakan:**
  1. Menulis skrip migrasi UP `000005_create_academic_years_and_enrollments.up.sql`:
     - Tabel `academic_years`: `id`, `school_unit_id` (FK), `tahun_ajaran`, `semester`, `is_active`, `tanggal_mulai`, `tanggal_selesai`.
     - Partial unique index `idx_academic_years_single_active` untuk menegakkan aturan tepat satu tahun ajaran aktif per unit sekolah (`WHERE is_active = true AND delete_at IS NULL`).
     - Tabel `class_enrollments`: `id`, `school_unit_id` (FK), `academic_year_id` (FK), `kelas_id` (FK), `siswa_id` (FK), `status`, `keterangan`.
     - Unique index `idx_class_enrollments_unique_student_year` untuk mencegah duplikasi pendaftaran siswa di tahun ajaran yang sama.
     - Seed awal tahun ajaran pilot `unit-sma-001` ('2026/2027', 'Ganjil', active).
     - Backfill otomatis riwayat siswa existing ke `class_enrollments`.
  2. Menulis skrip migrasi DOWN `000005_create_academic_years_and_enrollments.down.sql` untuk rollback aman.
  3. Mengimplementasikan domain entities, DTOs, dan interfaces di `backend/features/academic_years/entities.go`.
  4. Mengimplementasikan data layer model database PostgreSQL di `backend/features/academic_years/model/query.go` dengan query parameterized dan isolasi multi-unit.
  5. Mengimplementasikan business logic service layer di `backend/features/academic_years/service/logic.go` dengan validasi semester ('Ganjil'/'Genap'), validasi status pendaftaran ('Aktif', 'Lulus', 'Pindah', 'Tinggal Kelas'), proteksi penghapusan tahun ajaran yang sedang aktif, dan aktivasi atomik.
  6. Menulis comprehensive unit tests di `backend/features/academic_years/service/logic_test.go` dan `backend/features/academic_years/controllers/controller_test.go`.
  7. Mengimplementasikan presentation HTTP layer di `backend/features/academic_years/controllers/controller.go` dan formatter di `formatter.go`.
  8. Mendaftarkan route HTTP `/academic-years`, `/academic-years/detail`, `/academic-years/active`, `/academic-years/tambah`, `/academic-years/update`, `/academic-years/set-active`, `/academic-years/deleted`, `/class-enrollments`, `/class-enrollments/detail`, `/class-enrollments/tambah`, `/class-enrollments/update`, `/class-enrollments/deleted`, dan `/class-enrollments/student-history` pada `backend/router/router.go`.
- **Bukti / Verifikasi:**
  - `go test -v ./features/academic_years/...` berstatus `PASS` (seluruh skenario pengujian unit service dan controller lulus 100%).
  - `go test ./...` berstatus `PASS` (zero regression di seluruh 11 modul backend).
  - `go build ./...` sukses tanpa error maupun peringatan kompilasi.

### [2026-10-14] Task 04.3 — Pembuatan Tabel Relasi Wali Murid (guardians dan student_guardians)
- **Status:** `✅ Selesai dan terverifikasi`
- **Requirement:** OBJ-001, BUS-005, FR-008, NFR-SEC-001
- **Tindakan:**
  1. Menulis skrip migrasi UP `000006_create_guardians_and_student_guardians.up.sql`:
     - Tabel `guardians`: `id`, `user_id` (FK ke `users(id)`), `nama`, `nik`, `hubungan`, `pekerjaan`, `penghasilan`, `telepon`, `email`, `alamat`.
     - Tabel `student_guardians`: `id`, `guardian_id` (FK ke `guardians(id)`), `siswa_id` (FK ke `siswa(id)`), `is_primary`.
     - Unique index `idx_student_guardians_unique` pada `(guardian_id, siswa_id)` WHERE `delete_at IS NULL` untuk mencegah duplikasi pendaftaran relasi.
     - Backfill otomatis dari siswa existing yang memiliki data `nama_wali` dan `no_telepon_wali` ke tabel `guardians` dan `student_guardians`.
  2. Menulis skrip migrasi DOWN `000006_create_guardians_and_student_guardians.down.sql` untuk rollback aman.
  3. Mengimplementasikan domain entities, DTOs, dan interfaces di `backend/features/guardians/entities.go` (mendukung `SiswaList` multi-anak per wali murid).
  4. Mengimplementasikan data layer model database PostgreSQL di `backend/features/guardians/model/query.go` dengan query parameterized dan join info siswa, kelas, dan unit sekolah.
  5. Mengimplementasikan business logic service layer di `backend/features/guardians/service/logic.go` dengan validasi data, normalisasi hubungan (Ayah, Ibu, Wali, Kakak), pemutusan relasi (`UnlinkStudent`), dan penautan siswa (`LinkStudent`).
  6. Menulis unit testing komprehensif di `backend/features/guardians/service/logic_test.go` dan `backend/features/guardians/controllers/controller_test.go`.
  7. Mengimplementasikan presentation HTTP layer di `backend/features/guardians/controllers/controller.go` dan formatter di `formatter.go`.
  8. Mendaftarkan route HTTP `/guardians`, `/guardians/detail`, `/guardians/my-profile`, `/guardians/tambah`, `/guardians/update`, `/guardians/deleted`, `/guardians/link-student`, `/guardians/unlink-student`, `/guardians/students`, dan `/guardians/by-student` pada `backend/router/router.go`.
- **Bukti / Verifikasi:**
  - `go test -v ./features/guardians/...` berstatus `PASS` (seluruh skenario pengujian unit service dan controller lulus 100%).
  - `go test ./...` berstatus `PASS` (zero regression di seluruh modul backend).
  - `go build ./...` sukses tanpa error maupun peringatan kompilasi.

### [2026-10-14] Task 04.4 — Penyesuaian UI Frontend untuk Penautan Wali Murid dan Rombel Multi-Jenjang
- **Status:** `✅ Selesai dan terverifikasi`
- **Requirement:** OBJ-001, BUS-004, BUS-005, FR-006, FR-008, NFR-UX-001
- **Tindakan:**
  1. Membuat modul client axios baru di `frontend/src/api/guardiansApi.js` untuk rute CRUD wali murid dan relasi link/unlink student.
  2. Membuat modul client axios baru di `frontend/src/api/academicYearsApi.js` untuk master tahun ajaran dan riwayat penempatan rombel.
  3. Memperbarui `frontend/src/views/admin/Kelas.vue`:
     - Menambahkan filter dan input form tingkat sekolah komprehensif dari Kelas 1 s.d. 13 (SD, SMP, SMA, SMK 4 Tahun).
     - Menyesuaikan pilihan fase Kurikulum Merdeka secara berurutan: Fase A (SD Kelas 1-2), B (SD Kelas 3-4), C (SD Kelas 5-6), D (SMP Kelas 7-9), E (SMA Kelas 10), dan F (SMA/SMK Kelas 11-13).
  4. Memperbarui `frontend/src/views/admin/Siswa.vue`:
     - Menghubungkan dialog detail biodata siswa dengan pemanggilan relasi wali murid `getGuardiansBySiswaId`.
     - Menyediakan modal interaktif "Kelola / Tautkan Wali" (`showLinkModal`) yang memungkinkan admin memilih wali murid terdaftar atau mendaftarkan wali baru secara langsung (*on-the-fly*).
     - Menyediakan fitur pelepasan tautan (*unlink*) wali murid langsung dari antarmuka web.
- **Bukti / Verifikasi:**
  - `npm run build` di direktori `frontend/` berstatus `PASS` (Vite production build sukses tanpa error sintaks maupun bundling).
  - `go test ./...` dan `go build ./...` di direktori `backend/` tetap berstatus `PASS` (Zero regression).

### [2026-10-14] Task 05.1 & 05.2 — Verifikasi Deteksi Bentrok Jadwal Guru & Ruangan Serta Unit Testing Komprehensif
- **Status:** `✅ Selesai dan terverifikasi`
- **Requirement:** OBJ-001, FR-009, NFR-REL-001
- **Tindakan:**
  1. Melakukan audit komprehensif pada query database `CheckConflict` (`backend/features/jadwal/model/query.go`):
     - Memverifikasi pengecekan irisan jam mengajar guru pada hari dan semester/tahun ajaran aktif yang sama (`j2.jam_mulai < $7::time AND j2.jam_selesai > $6::time`).
     - Memverifikasi pengecekan tabrakan jadwal rombel/kelas.
     - Memverifikasi pengecekan okupansi ruangan sekolah agar tidak ada 2 kelas di ruangan yang sama pada jam bertabrakan.
  2. Menulis unit test komprehensif di `backend/features/jadwal/service/logic_test.go`:
     - Test validasi waktu (jam selesai harus lebih besar dari jam mulai).
     - Test deteksi bentrok jadwal guru mengajar (`ConflictInfo.HasConflict = true`).
     - Test deteksi bentrok okupansi ruangan sekolah.
     - Test jadwal valid tanpa bentrok (`ConflictInfo.HasConflict = false`).
- **Bukti / Verifikasi:**
  - `go test -v ./features/jadwal/service/...` berstatus `PASS` (100% skenario uji anti-bentrok lulus).
  - `go test ./...` berstatus `PASS` (zero regression di seluruh modul backend).
  - `go build ./...` sukses tanpa error maupun peringatan kompilasi.

### [2026-10-14] Task 05.3 — UI Frontend Kalender / Matriks Jadwal Mingguan Bebas Bentrok
- **Status:** `✅ Selesai dan terverifikasi`
- **Requirement:** OBJ-001, FR-009, NFR-UX-001
- **Tindakan:**
  1. Menambahkan switcher mode tampilan di `frontend/src/views/admin/JadwalPelajaran.vue`:
     - Mode `Tabel Sesi`: Daftar linear lengkap dengan filter multi-kriteria dan aksi cepat.
     - Mode `Matriks Mingguan`: Grid kalender 6 hari KBM (Senin - Sabtu) yang mengelompokkan sesi pembelajaran per hari secara kronologis.
  2. Menyediakan fungsi filter kronologis `getJadwalByDay(dayInt)` untuk mengurutkan sesi belajar berdasarkan jam mulai.
  3. Memastikan peringatan error konflik dari backend tertangkap dan ditampilkan dengan jelas di modal tambah/ubah jadwal.
- **Bukti / Verifikasi:**
  - `npm run build` di direktori `frontend/` berstatus `PASS` (Vite build sukses dalam 860ms tanpa error).
  - `go test ./...` dan `go build ./...` di direktori `backend/` tetap berstatus `PASS` (Zero regression).

### [2026-10-14] Task 06.1 — Refactor Audit Logging Middleware Asinkron & Payload Sanitization
- **Status:** `✅ Selesai dan terverifikasi`
- **Requirement:** OBJ-001, NFR-SEC-001, NFR-PERF-001
- **Tindakan:**
  1. Memperketat fungsi sanitasi `maskSensitiveData` di `backend/helper/transaksi_log.go` sesuai `AGENTS.md` Bagian 18.5:
     - Mendeteksi dan mereduksi field sensitif: `password`, `kata_sandi`, `old_password`, `new_password`, `access_token`, `refresh_token`, `token`, `secret`, `private_key`, `authorization`, `api_key` menjadi label aman `[REDACTED]`.
  2. Memperluas identifikasi modul aktivitas di `DetermineActivityDetail` untuk mencakup modul satelit baru:
     - Modul `Unit Sekolah` (`/school-units`)
     - Modul `Tahun Ajaran` (`/academic-years`)
     - Modul `Riwayat Rombel Siswa` (`/class-enrollments`)
     - Modul `Wali Murid` (`/guardians`)
  3. Menulis unit testing `TestMaskSensitiveData` dan `TestDetermineActivityDetail_NewModules` di `backend/helper/transaksi_log_test.go` untuk memvalidasi zero credential leak pada audit trail.
- **Bukti / Verifikasi:**
  - `go test -v ./helper/...` berstatus `PASS` (100% skenario pengujian sanitasi & audit log lulus).
  - `go test ./...` berstatus `PASS` (zero regression di seluruh 12 paket backend).
  - `go build ./...` sukses tanpa error maupun peringatan kompilasi.

### [2026-10-14] Task 06.2 — Pengetatan Konfigurasi CORS (CORS Hardening)
- **Status:** `✅ Selesai dan terverifikasi`
- **Requirement:** NFR-SEC-001, OBJ-001
- **Tindakan:**
  1. Memperbarui `CORSMiddleware` di `backend/helper/middlewares.go` untuk membaca `CORS_ALLOWED_ORIGINS` dari environment variable.
  2. Menggantikan wildcard liar dengan validasi origin dinamis yang aman, mendukung `Access-Control-Allow-Credentials: true`, header `X-School-Unit-ID`, dan max age preflight cache 24 jam.
  3. Menulis unit test `TestCORSMiddleware` di `backend/helper/transaksi_log_test.go` untuk memverifikasi handling preflight request OPTIONS dan origin matching.
- **Bukti / Verifikasi:**
  - `go test -v ./helper/...` berstatus `PASS` (Termasuk `TestCORSMiddleware`).
  - `go test ./...` dan `go build ./...` berstatus `PASS` (Zero regression).

### [2026-10-14] Task 06.3 — Dashboard Agregat Yayasan & Operasional Admin Unit
- **Status:** `✅ Selesai dan terverifikasi`
- **Requirement:** BUS-001, FR-001, FR-004
- **Tindakan:**
  1. Menambahkan atribut akumulasi master yayasan pada entitas backend `DashboardTotals`:
     - `total_unit`: Total unit sekolah aktif yayasan.
     - `total_wali`: Total wali murid terdaftar dalam sistem.
     - `total_tahun_aktif`: Total periode tahun ajaran aktif berjalan.
  2. Memperbarui query akumulasi di `backend/features/dashboard/model/query.go` untuk menghitung ketiga metrik di atas secara efisien.
  3. Memperbarui antarmuka `frontend/src/views/admin/Dashboard.vue` dengan menambahkan kartu ringkasan visual untuk Unit Sekolah Yayasan, Wali Murid Terhubung, dan Tahun Ajaran Aktif.
- **Bukti / Verifikasi:**
  - `npm run build` di `frontend/` berstatus `PASS` (Vite build sukses dalam 905ms).
  - `go test ./...` dan `go build ./...` berstatus `PASS` (Zero regression).
