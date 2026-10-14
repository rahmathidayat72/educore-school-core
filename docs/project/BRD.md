# Business Requirements Document (BRD)
## EduCore School — School Core Management

## 1. Metadata dan Revision History
- **Proyek:** EduCore School
- **Modul:** School Core Management
- **Versi:** 1.0.0
- **Tanggal Persetujuan:** 14 Oktober 2026
- **Status:** APPROVED (Gerbang 1A)
- **Penyusun:** Antigravity AI Agent
- **Disetujui Oleh:** Product Owner / Tech Lead

---

## 2. Executive Summary
EduCore School dirancang sebagai ekosistem digital terpadu untuk menaungi berbagai jenjang institusi pendidikan (SD/MI, SMP/MTs, SMA/MA, SMK sederajat) di bawah naungan satu yayasan. Modul *School Core Management* bertindak sebagai penyedia data induk tunggal (*Single Source of Truth*) yang mengelola profil sekolah, identitas peserta didik, tenaga pendidik, wali murid, rombel kelas, kurikulum/mata pelajaran, dan jadwal tatap muka. Sistem memastikan data terisolasi secara aman per unit sekolah bagi admin lokal, sekaligus memberikan visibilitas komprehensif bagi pengurus yayasan melalui hak akses Super Admin.

---

## 3. Latar Belakang
Pengelolaan sekolah multi-jenjang dalam satu yayasan kerap menghadapi kendala sinkronisasi data antar-unit dan keterpisahan sistem administrasi. Tanpa sistem inti yang terpusat, pengurus yayasan kesulitan memantau metrik agregat, sementara staf administrasi unit sekolah harus melakukan input data berulang ke berbagai aplikasi penunjang (presensi, e-rapor, keuangan, dan ujian). EduCore School dibangun untuk menghilangkan inefisiensi tersebut melalui arsitektur modular yang dapat dikonsumsi oleh seluruh aplikasi satelit.

---

## 4. Pernyataan dan Dampak Masalah
- **Masalah 1 — Data Terfragmentasi:** Data siswa, guru, dan rombel tersebar di spreadsheet atau database terpisah per sekolah.
  - *Dampak:* Rekonsiliasi data yayasan memakan waktu lama, rawan salah hitung kuota, dan sulit menyusun laporan tahunan.
- **Masalah 2 — Integrasi Rapuh:** Aplikasi pendukung (CBT, SPP, Presensi) memiliki master data masing-masing.
  - *Dampak:* Ketidaksinkronan status siswa aktif/mutasi/lulus pada aplikasi satelit.
- **Masalah 3 — Ketiadaan Isolasi Data Aman:** Sistem single-tenant konvensional tidak mampu memisahkan wewenang admin SD, SMP, SMA/SMK dalam satu server terpusat.
  - *Dampak:* Risiko kebocoran data privasi siswa lintas jenjang atau intervensi data antar-operator unit.

---

## 5. Tujuan Bisnis (Business Objectives)
- **OBJ-001:** Memastikan 100% data induk peserta didik, pendidik, dan rombongan belajar terkelola secara terpusat dengan integritas referensial tinggi.
- **OBJ-002:** Mengakomodasi model yayasan multi-unit dengan pemisahan akses tingkat database (*tenant isolation*) antar-unit sekolah.
- **OBJ-003:** Menyediakan API terbuka berstandar industri bagi modul satelit ekosistem (PPDB, Presensi, E-Rapor, Keuangan, LMS, CBT).
- **OBJ-004:** Mencegah kesalahan operasional penempatan rombel dan bentrok jadwal mengajar guru secara otomatis.

---

## 6. Stakeholder dan RACI

| Stakeholder | Responsible | Accountable | Consulted | Informed |
|---|:---:|:---:|:---:|:---:|
| **Pengurus Yayasan (Super Admin)** | | X | X | |
| **Admin Unit Sekolah (Tata Usaha)** | X | | X | |
| **Kepala Sekolah / Wakil Kurikulum** | | | X | X |
| **Guru / Pendidik** | X | | | X |
| **Orang Tua / Wali Murid** | | | | X |
| **Siswa & Alumni** | | | | X |

---

## 7. Proses As-Is vs To-Be

### 7.1 Proses As-Is
1. Setiap unit sekolah mencatat data siswa dan guru secara mandiri di file lokal atau aplikasi terpisah.
2. Penempatan kelas dilakukan manual di atas lembar kerja spreadsheet; riwayat kelas tahun sebelumnya sering tertimpa.
3. Penyusunan jadwal pelajaran rentan bentrok guru atau ruang kelas karena tidak ada validasi sistemik.
4. Laporan yayasan harus dikompilasi manual melalui korespondensi email atau pesan instan.

### 7.2 Proses To-Be
1. Unit sekolah yayasan terdaftar di Core Management; Super Admin memantau seluruh unit dari satu dashboard.
2. Operator unit menginput siswa dan menempatkan pada rombel terikat tahun ajaran (`class_enrollments`) sehingga rekam jejak historis tetap utuh.
3. Pembuatan jadwal memverifikasi ketersediaan guru dan ruangan secara *real-time*; bentrok langsung ditolak sistem.
4. Aplikasi satelit (Presensi, Keuangan, LMS) mengambil data identitas langsung dari Core API tanpa duplikasi.

---

## 8. Ruang Lingkup Bisnis

### 8.1 In Scope (Fase 1 — Core MVP)
- Manajemen Unit Sekolah Yayasan (`school_units`).
- Manajemen Identitas & Multi-Role User (`super_admin`, `admin`, `guru`, `wali`, `siswa`, `alumni`).
- Master data siswa, orang tua/wali murid, guru, kelas, mata pelajaran, dan jadwal.
- Audit trail transaksi mutasi data.
- Dashboard ringkasan yayasan dan unit sekolah.
- Kontrak REST API terstandarisasi untuk integrasi modular.

### 8.2 Out of Scope (Fase Lanjutan)
- Modul Smart Attendance (Presensi QR/RFID).
- Parent Portal PWA mandiri.
- Portal Alumni Tracer Study mandiri.
- School Finance & Payment Gateway.
- CBT & LMS online.

---

## 9. Business Requirements

| ID | Judul | Deskripsi | Prioritas | Indikator Penerimaan |
|---|---|---|:---:|---|
| **BUS-001** | Multi-Unit Yayasan | Sistem mendukung entitas sekolah jamak di bawah satu yayasan. | High | Super Admin dapat menambah, mengubah, dan memantau unit sekolah. |
| **BUS-002** | Isolasi Data Unit | Akses admin sekolah dibatasi ketat hanya pada data unitnya sendiri. | Critical | Query SQL secara otomatis memfilter berdasarkan `school_unit_id`. |
| **BUS-003** | Fleksibilitas Jenjang | Rombel dan mapel mengakomodasi aturan SD, SMP, SMA, dan SMK. | High | Atribut fase/tingkat dan konsentrasi keahlian/jurusan tersedia. |
| **BUS-004** | Riwayat Akademik | Siswa naik kelas tanpa menghapus rekam jejak rombel tahun sebelumnya. | High | Menggunakan tabel relasi `class_enrollments` terikat tahun ajaran. |
| **BUS-005** | Relasi Wali Murid | Satu akun wali dapat terhubung ke lebih dari satu anak di sekolah yayasan. | Medium | Relasi *many-to-many* antara siswa dan wali murid. |
| **BUS-006** | Validasi Bentrok Jadwal | Sistem mendeteksi bentrok jadwal mengajar guru dan ruang kelas. | High | Penyimpanan jadwal gagal jika ada bentrok waktu yang tumpang tindih. |

---

## 10. Business Rules

- **BRULE-001 (Isolasi Tenant):** Pengguna dengan role `admin` unit dilarang mengakses, mengubah, atau menghapus data milik unit sekolah lain.
- **BRULE-002 (Tahun Ajaran Aktif Tunggal):** Hanya boleh ada tepat satu tahun ajaran aktif per unit sekolah pada satu rentang waktu.
- **BRULE-003 (Siklus Hidup Siswa):** Siswa yang berstatus "Lulus" otomatis masuk ke basis data alumni tanpa menghapus riwayat akademik.
- **BRULE-004 (Keutuhan Audit Log):** Setiap aksi penambahan, perubahan, dan penghapusan data wajib dicatat di audit log dan tidak boleh dimanipulasi oleh operator biasa.

---

## 11. KPI dan Success Metrics
- **KPI-001:** Zero data overlap (0% kebocoran data antar-unit sekolah).
- **KPI-002:** 100% deteksi bentrok jadwal mengajar pada saat penyusunan jadwal.
- **KPI-003:** Kecepatan respon API rata-rata < 250 ms.
- **KPI-004:** Waktu rekapitulasi data yayasan berkurang dari hitungan hari menjadi instan (*real-time*).

---

## 12. Constraint, Dependency, dan Risiko
- **Constraint:** Tetap menggunakan fondasi Go Native + PostgreSQL + Vue 3.
- **Dependency:** Integritas skema data induk harus stabil sebelum modul ekosistem satelit mulai dibangun.
- **Risiko:** Perubahan skema tanpa *database migration tool* dapat merusak backward-compatibility API.

---

## 13. Approval
- **Status:** **DISSETUJUI (APPROVED)**
- **Tanggal:** 14 Oktober 2026
