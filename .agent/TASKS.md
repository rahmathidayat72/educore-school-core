# Task Tracking & Backlog — EduCore School

## 1. Active Tasks (Phase 01)

| Task ID | Deskripsi Task | Assignee | Status | DoR | DoD | Requirement Ref |
|---|---|:---:|:---:|:---:|:---:|---|
| `task-01.1` | Setup struktur tata kelola dokumen di `docs/project/` dan `.agent/` | AI Agent | `✅ Selesai dan terverifikasi` | Ready | Done | OBJ-001, OBJ-002 |
| `task-01.2` | Penyelarasan file `AGENTS.md` di root `School-Core-Management/` | AI Agent | `✅ Selesai dan terverifikasi` | Ready | Done | SOP-001 |
| `task-01.3` | Pembersihan artefak lokal (`server.log`, binary) & update `.gitignore` | AI Agent | `✅ Selesai dan terverifikasi` | Ready | Done | SOP-001 |
| `task-01.4` | Rekomendasi/setup kontrol versi Git terpadu | AI Agent | `✅ Selesai dan terverifikasi` | Ready | Done | SOP-001 |

---

## 2. Upcoming Roadmap Tasks

### Phase 02: Multi-Unit Yayasan & Skema Database
- [x] `task-02.1`: Buat skrip migrasi SQL tabel `school_units`. (Status: `✅ Selesai dan terverifikasi`)
- [x] `task-02.2`: Tambahkan kolom `school_unit_id` pada seluruh entitas operasional. (Status: `✅ Selesai dan terverifikasi`)
- [x] `task-02.3`: Penyiapan data default unit percontohan & data backfill. (Status: `✅ Selesai dan terverifikasi`)

### Phase 03: Identity, Role Super Admin, & Tenant Isolation
- [x] `task-03.1`: Perbarui entitas `users` dan JWT claims (`school_unit_id`, `role: super_admin`). (Status: `✅ Selesai dan terverifikasi`)
- [ ] `task-03.2`: Implementasi isolasi data unit sekolah pada service & model (WHERE `school_unit_id`). (Status: `📋 Siap`)
- [ ] `task-03.3`: Implementasi endpoint CRUD Unit Sekolah khusus Super Admin. (Status: `⬜ Belum siap`)

### Phase 04: Penguatan Modul Inti (Siswa, Wali, Guru, Rombel)
- [ ] `task-04.1`: Penyesuaian tabel rombel/kelas untuk jenjang SD (Fase A-C), SMP (Fase D), SMA (Fase E-F), dan SMK (Jurusan).
- [ ] `task-04.2`: Pembuatan tabel `academic_years` dan `class_enrollments`.
- [ ] `task-04.3`: Pembuatan tabel relasi wali murid `guardians` dan `student_guardians`.
- [ ] `task-04.4`: Penyesuaian UI frontend untuk penautan wali murid dan rombel multi-jenjang.

### Phase 05: Penjadwalan Bebas Bentrok (Conflict-Free Schedule)
- [ ] `task-05.1`: Implementasi deteksi bentrok jadwal guru & ruangan di `features/jadwal/service`.
- [ ] `task-05.2`: Unit test komprehensif logika anti-bentrok.
- [ ] `task-05.3`: UI frontend kalender/tabel jadwal dengan peringatan bentrok waktu.

### Phase 06: Audit Trail, CORS Hardening, & Dashboard
- [ ] `task-06.1`: Refactor audit logging middleware menjadi asinkron dengan payload sanitization.
- [ ] `task-06.2`: Pengetatan CORS configuration.
- [ ] `task-06.3`: Dashboard agregat Super Admin & dashboard operasional Admin Unit.

### Phase 07: Kontrak REST API & Final Acceptance
- [ ] `task-07.1`: Standarisasi REST alias endpoint untuk modul satelit.
- [ ] `task-07.2`: Update koleksi Bruno testing (`backend/bruno/`).
- [ ] `task-07.3`: Pengujian regresi penuh & final acceptance report.
