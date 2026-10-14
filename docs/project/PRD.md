# Product Requirements Document (PRD)
## EduCore School — School Core Management

## 1. Metadata dan Revision History
- **Proyek:** EduCore School
- **Modul:** School Core Management
- **Versi:** 1.0.0
- **Tanggal Persetujuan:** 14 Oktober 2026
- **Status:** APPROVED (Gerbang 1B)
- **Referensi BRD:** `docs/project/BRD.md`

---

## 2. Ringkasan Produk
EduCore School Core Management adalah sistem manajemen data operasional sekolah berbasis web portal dan RESTful API terpusat. Produk ini memungkinkan yayasan pendidikan mengelola multi-unit sekolah (SD/SMP/SMA/SMK) dengan hak akses berjenjang, data siswa terintegrasi dengan wali murid, manajemen pendidik, penempatan rombongan belajar berbasis tahun ajaran, dan penjadwalan pelajaran bebas bentrok.

---

## 3. Persona dan Role Matrix

### 3.1 Persona Pengguna
1. **Super Admin Yayasan:** Pimpinan yayasan yang memantau agregat metrik seluruh unit sekolah dan mengelola master unit sekolah.
2. **Admin Unit Sekolah:** Staf tata usaha yang mengelola operasional harian data siswa, guru, kelas, mapel, dan jadwal di unit sekolahnya.
3. **Guru / Tenaga Pendidik:** Mengakses jadwal mengajar dan melihat daftar siswa di rombel yang diajar.
4. **Orang Tua / Wali:** Memantau jadwal belajar, data guru wali kelas, dan informasi anak yang terhubung sah.
5. **Siswa:** Melihat jadwal pelajaran rombelnya dan informasi akademik pribadi.
6. **Alumni:** Mengakses riwayat kelulusan dan pembaruan data kontak/karir.

### 3.2 Role and Access Matrix

| Resource | Super Admin | Admin Unit | Guru | Wali Murid | Siswa | Alumni |
|---|:---:|:---:|:---:|:---:|:---:|:---:|
| `school_units` | CRUD | Read (Unit) | Read (Unit) | Read (Unit) | Read (Unit) | No |
| `users` | CRUD | CRUD (Unit) | Self Read/Update | Self Read/Update | Self | Self |
| `guru` | Read Semua | CRUD (Unit) | Read (Unit) | Read (Pengampu) | Read | No |
| `kelas` | Read Semua | CRUD (Unit) | Read / Manage Binaan | Read (Rombel Anak) | Read | No |
| `siswa` | Read Semua | CRUD (Unit) | Read (Siswa Kelas) | Read (Anak Sendiri) | Self | No |
| `mata_pelajaran` | Read Semua | CRUD (Unit) | Read (Unit) | Read (Mapel Anak) | Read | No |
| `jadwal_pelajaran` | Read Semua | CRUD (Unit) | Read (Jadwal Sendiri) | Read (Jadwal Anak) | Read | No |
| `alumni` | Read Semua | CRUD (Unit) | Read (Unit) | No | No | Self Profile |
| `transaction_logs`| Read Semua | Read (Unit) | No | No | No | No |
| `dashboard` | Global Stats | Unit Stats | Teacher Stats | Parent Stats | Student Stats | Profile |

---

## 4. Scope Produk

### 4.1 MVP (Fase 1)
- CRUD Unit Sekolah Yayasan (`school_units`).
- Autentikasi JWT dengan klaim `role` dan `school_unit_id`.
- Master Data Siswa, Guru, Kelas, Mapel, dan Jadwal dengan filter unit sekolah otomatis.
- Deteksi bentrok jadwal pelajaran (guru dan ruangan).
- Manajemen data Orang Tua / Wali Murid dan pemetaan relasi ke Siswa.
- Riwayat kelas siswa berbasis tahun ajaran (`class_enrollments`).
- Audit trail transaksi otomatis.
- Dashboard operasional terpisah (Super Admin Yayasan vs Admin Unit).

### 4.2 Release Berikutnya (Fase 2+)
- Smart Attendance (Presensi QR Code / RFID).
- Parent Portal PWA mandiri.
- E-Rapor & Penilaian Kurikulum Merdeka.
- Modul Keuangan Sekolah (SPP & Tagihan).
- Modul PPDB / SPMB Terpadu.
- Modul CBT dan LMS.

---

## 5. Functional Requirements (FR)

### Modul A: Multi-Unit Yayasan & Autentikasi
- **FR-001:** Super Admin dapat mengelola unit sekolah (`id`, `nama_sekolah`, `jenjang`, `npsn`, `status`).
- **FR-002:** Sistem menyediakan endpoint login `/login` yang memvalidasi email/password dan menghasilkan token JWT dengan claims `id`, `role`, dan `school_unit_id`.
- **FR-003:** Backend memverifikasi `school_unit_id` pada setiap query database untuk role selain Super Admin (*server-side tenant isolation*).

### Modul B: Akademik & Rombel
- **FR-004:** Admin unit dapat mengelola tahun ajaran (`academic_years`) dengan validasi tepat satu tahun ajaran yang aktif per unit.
- **FR-005:** Pembuatan kelas mendukung atribut jenjang:
  - SD: Fase A, B, C / Tingkat 1–6.
  - SMP: Fase D / Tingkat 7–9.
  - SMA: Fase E, F / Tingkat 10–12 dengan peminatan.
  - SMK: Memiliki atribut khusus Program/Konsentrasi Keahlian (Jurusan).
- **FR-006:** Sistem mencatat penempatan siswa ke kelas melalui `class_enrollments` yang terikat tahun ajaran sehingga histori kelas tidak tertimpa.

### Modul C: Siswa, Guru, dan Wali Murid
- **FR-007:** Admin unit dapat mengelola data siswa lengkap (NIS, NISN, status, biodata) dengan validasi keunikan NISN.
- **FR-008:** Admin unit dapat menautkan akun wali murid dengan satu atau lebih siswa di bawah yayasan (*many-to-many relationship*).
- **FR-009:** Admin unit dapat mengelola data guru beserta penugasan mata pelajaran.
- **FR-010:** Mengubah status siswa menjadi "Lulus" otomatis menampilkan data pada inventaris data alumni unit sekolah.

### Modul D: Penjadwalan & Operasional
- **FR-011:** Sistem memvalidasi ketersediaan guru dan ruangan saat pembuatan jadwal; menolak jika terjadi bentrok jam mengajar.
- **FR-012:** Sistem mencatat seluruh mutasi data ke tabel `transaction_logs` di background tanpa memperlambat response HTTP.
- **FR-013:** Dashboard menyajikan grafik distribusi siswa per kelas, rasio gender, ringkasan guru, dan pemenuhan jadwal mengajar.

---

## 6. Non-Functional Requirements (NFR)
- **NFR-SEC-001:** Password wajib di-hash menggunakan `bcrypt` dengan salt cost minimal 10.
- **NFR-SEC-002:** Token JWT kedaluwarsa setelah durasi tertentu dan diverifikasi menggunakan secret key aman.
- **NFR-SEC-003:** Otorisasi wajib ditegakkan di sisi server (*default deny* dan *role-based boundary*).
- **NFR-PERF-001:** Waktu pemrosesan API rata-rata < 250 ms untuk endpoint query standar.
- **NFR-REL-001:** Mutasi yang melibatkan relasi ganda wajib menggunakan transaksi database ACID (`BEGIN ... COMMIT`).
- **NFR-UX-001:** Antarmuka responsif pada resolusi desktop, tablet, dan smartphone.

---

## 7. Acceptance Criteria Produk
- [x] Super Admin dapat beralih konteks untuk melihat data unit SD, SMP, SMA, atau SMK.
- [x] Admin unit SD terbukti tidak dapat mengakses atau melihat data unit SMP/SMA/SMK.
- [x] Riwayat kelas siswa tetap tercatat saat siswa naik ke tingkat berikutnya.
- [x] Satu akun wali murid dapat terhubung ke dua siswa di unit yang berbeda.
- [x] Pembuatan jadwal dengan guru dan waktu yang sama pada dua kelas berbeda berhasil ditolak dengan pesan bentrok.
- [x] Seluruh endpoint API lolos pengujian Bruno dan unit testing Go.
