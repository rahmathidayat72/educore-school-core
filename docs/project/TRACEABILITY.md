# Requirements Traceability Matrix (RTM) — EduCore School

## 1. Pemetaan Menyeluruh: Objective → Business Req → Functional Req → Phase & Task

| Objective ID | Business Req (BRD) | Functional Req (PRD) | Phase Implementasi | Task ID | Komponen / Endpoint Terkait | Kriteria Verifikasi |
|---|---|---|:---:|:---:|---|---|
| **OBJ-001** | BUS-001 (Multi-Unit) | FR-001, FR-002 | Phase 02, 03 | `task-02.1`, `task-03.1` | `school_units`, `features/users` | Super Admin CRUD unit sekolah |
| **OBJ-002** | BUS-002 (Tenant Isolation) | FR-003 | Phase 03 | `task-03.2` | `helper/middleware`, model queries | Query otomatis filter `school_unit_id` |
| **OBJ-001** | BUS-003 (Jenjang SD-SMK) | FR-005 | Phase 04 | `task-04.1` | `features/kelas`, `features/mapel` | Atribut fase & konsentrasi keahlian |
| **OBJ-001** | BUS-004 (Riwayat Rombel) | FR-004, FR-006 | Phase 04 | `task-04.2` | `academic_years`, `class_enrollments` | Riwayat kelas tidak hilang saat naik kelas |
| **OBJ-001** | BUS-005 (Relasi Wali Murid)| FR-008 | Phase 04 | `task-04.3` | `guardians`, `student_guardians` | 1 wali bisa terhubung ke multi-anak |
| **OBJ-004** | BUS-006 (Anti-Bentrok Jadwal)| FR-011 | Phase 05 | `task-05.1` | `features/jadwal/service` | Penolakan otomatis bentrok guru/ruang |
| **OBJ-003** | BUS-002 (Audit & Standard) | FR-012, FR-013 | Phase 06, 07 | `task-06.1`, `task-07.1` | `transaction_logs`, REST endpoints | Payload log tersanitasi & REST alias |

---

## 2. Pemetaan Non-Functional Requirements (NFR)

| NFR ID | Deskripsi Persyaratan | Modul Terkait | Mitigasi Teknis | Verifikasi Pengujian |
|---|---|---|---|---|
| **NFR-SEC-001** | Password bcrypt cost >= 10 | `features/users` | `bcrypt.GenerateFromPassword` | Unit test hash check & login check |
| **NFR-SEC-002** | JWT Token Expiration & Secret | `helper/jwt` | `golang-jwt/jwt/v5` claims `exp` | Unit test expired token response 401 |
| **NFR-SEC-003** | Server-side Tenant Enforcement | Backend Middleware & Model | Filter `school_unit_id` di klausa WHERE | Cross-tenant access test return 403 / empty |
| **NFR-PERF-001** | API Response Time < 250 ms | Seluruh Endpoint REST | Connection pool `pgxpool` & indexing DB | Benchmark testing & latency monitoring |
| **NFR-REL-001** | ACID Transaction Mutasi Relasi| `siswa`, `jadwal`, `enrollment` | `tx, err := pool.Begin(ctx)` | Rollback test saat simpan relasi parsial gagal |
