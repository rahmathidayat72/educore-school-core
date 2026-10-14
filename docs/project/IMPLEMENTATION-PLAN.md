# Master Implementation Plan — EduCore School

## 1. Roadmap Fase Implementasi

```text
Phase 01 ── Fondasi & Tata Kelola Proyek (Standar AGENTS.md, Git, Kebersihan Repo)
   │
Phase 02 ── Multi-Unit Yayasan & Evolusi Skema Data (school_units & school_unit_id)
   │
Phase 03 ── Identity, Role Super Admin, & Server-Side Tenant Isolation
   │
Phase 04 ── Penguatan Modul Inti (Siswa, Wali Murid, Guru, Rombel Multi-Jenjang)
   │
Phase 05 ── Modul Penjadwalan Bebas Bentrok (Conflict-Free Timetable)
   │
Phase 06 ── Audit Trail Asinkron, CORS Hardening, & Dashboard Multi-Unit
   │
Phase 07 ── Standarisasi Kontrak REST API untuk Satelit & Final Acceptance Verification
```

---

## 2. Rincian Task Tiap Fase

### Phase 01: Fondasi & Tata Kelola Proyek
- **Task 1.1:** Setup struktur dokumen proyek di `docs/project/` dan pelacakan AI di `.agent/`.
- **Task 1.2:** Inisialisasi file `AGENTS.md` di root `School-Core-Management/` dan sinkronisasi workspace.
- **Task 1.3:** Pembersihan artefak sampah lokal (`server.log`, binary kompilasi `go_rest_native_sekolah`) dan penyelarasan `.gitignore`.
- **Task 1.4:** Rekomendasi inisialisasi Git repository terpadu di root workspace.

### Phase 02: Multi-Unit Yayasan & Skema Database
- **Task 2.1:** Pembuatan skrip migrasi terversi SQL untuk tabel `school_units`.
- **Task 2.2:** Penambahan foreign key `school_unit_id` pada seluruh entitas operasional (`users`, `guru`, `kelas`, `siswa`, `mata_pelajaran`, `jadwal_pelajaran`, `transaction_logs`).
- **Task 2.3:** Penyiapan data default unit percontohan (*pilot unit*) dan backfill data yang sudah ada.

### Phase 03: Identity, Role, & Tenant Isolation
- **Task 3.1:** Pembaruan entitas `users` dan JWT middleware untuk menyertakan klaim `school_unit_id` dan role `super_admin`.
- **Task 3.2:** Implementasi isolasi data unit sekolah pada service & model (query WHERE `school_unit_id`).
- **Task 3.3:** Endpoint CRUD Unit Sekolah khusus Super Admin.

### Phase 04: Penguatan Modul Inti (Siswa, Wali, Guru, Rombel)
- **Task 4.1:** Penyesuaian tabel kelas untuk jenjang SD (Fase A-C), SMP (Fase D), SMA (Fase E-F), dan SMK (Program Keahlian).
- **Task 4.2:** Pembuatan tabel `academic_years` dan `class_enrollments` untuk menjaga rekam jejak kenaikan kelas siswa.
- **Task 4.3:** Pembuatan tabel `guardians` dan `student_guardians` untuk relasi siswa dan orang tua/wali murid.
- **Task 4.4:** Penyesuaian UI frontend untuk mendukung penautan wali murid dan input kelas multi-jenjang.

### Phase 05: Penjadwalan Bebas Bentrok (Conflict-Free Schedule)
- **Task 5.1:** Implementasi fungsi deteksi bentrok jadwal guru dan ruangan di `features/jadwal/service`.
- **Task 5.2:** Unit testing komprehensif logika anti-bentrok jadwal.
- **Task 5.3:** Penyesuaian UI frontend kalender/tabel jadwal dengan peringatan bentrok waktu.

### Phase 06: Audit Trail, CORS Hardening, & Dashboard
- **Task 6.1:** Refactor logging middleware menjadi goroutine asinkron dengan masking kredensial otomatis.
- **Task 6.2:** Pengetatan konfigurasi CORS dari `*` menjadi domain whitelist terdaftar.
- **Task 6.3:** Dashboard Yayasan untuk Super Admin (rekapitulasi lintas unit) dan Dashboard Operasional untuk Admin Unit.

### Phase 07: Kontrak REST API & Final Acceptance
- **Task 7.1:** Penyediaan endpoint REST alias terstandarisasi untuk konsumsi aplikasi satelit.
- **Task 7.2:** Pembaruan koleksi Bruno (`backend/bruno/`) untuk seluruh skenario positif dan negatif.
- **Task 7.3:** Eksekusi verifikasi akhir (seluruh unit test pass, frontend build clean, verifikasi acceptance criteria).
