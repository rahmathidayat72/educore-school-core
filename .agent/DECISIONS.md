# Architecture & Business Decision Records (ADR) — EduCore School

## Record 001: Pemilihan Tech Stack Utama
- **Status:** APPROVED
- **Tanggal:** 14 Oktober 2026
- **Konteks:** Sistem membutuhkan performa tinggi, efisiensi memori, dan kemudahan deployment on-premise maupun cloud untuk yayasan multi-sekolah.
- **Keputusan:** Mempertahankan basis kode yang ada: **Go Native (`net/http`)** untuk backend dan **Vue 3 (Vite + Pinia + Bootstrap 5)** untuk frontend, dengan PostgreSQL sebagai DBMS utama via driver `pgx/v5`.
- **Konsekuensi:** Performa maksimal dan latensi rendah tanpa overhead framework pihak ketiga yang berat, namun penanganan routing dan middleware ditulis secara native dan modular.

---

## Record 002: Arsitektur Ekosistem Multi-Repo & Single Source of Truth
- **Status:** APPROVED
- **Tanggal:** 14 Oktober 2026
- **Konteks:** EduCore School mencakup berbagai modul pendukung (Presensi, Keuangan/SPP, E-Rapor, LMS, CBT, PPDB) yang dikembangkan terpisah.
- **Keputusan:** Modul **School Core Management** bertindak sebagai master data induk tunggal (*Single Source of Truth*). Modul lain mengonsumsi data identitas dan akademik melalui RESTful API terstandarisasi.
- **Konsekuensi:** Mencegah redundansi data antar-aplikasi; perubahan skema inti memerlukan sinkronisasi kontrak API.

---

## Record 003: Model Isolasi Tenant Multi-Unit Yayasan
- **Status:** APPROVED
- **Tanggal:** 14 Oktober 2026
- **Konteks:** Yayasan menaungi beberapa jenjang sekolah (SD, SMP, SMA, SMK) di mana admin unit hanya boleh mengakses unitnya sendiri, sedangkan Super Admin yayasan dapat memantau seluruh unit.
- **Keputusan:** Menerapkan *shared database, shared schema* dengan kolom pembeda `school_unit_id` pada setiap entitas operasional, dengan penegakan isolasi mutlak di server-side query database (bukan sekadar filter UI).
- **Konsekuensi:** Biaya operasional database efisien, namun seluruh query model database wajib menyertakan filter `WHERE school_unit_id = $1`.

---

## Record 004: Validasi Bentrok Jadwal Mengajar (Conflict-Free Scheduling)
- **Status:** APPROVED
- **Tanggal:** 14 Oktober 2026
- **Konteks:** Penyusunan jadwal pelajaran di unit sekolah sering mengalami bentrok guru atau ruang kelas.
- **Keputusan:** Memvalidasi ketersediaan guru dan ruangan secara otomatis di backend service layer saat penyimpanan jadwal baru atau perubahan jadwal. Jika ada irisan rentang waktu pada hari yang sama, request ditolak dengan pesan HTTP 400 Bad Request.
- **Konsekuensi:** Mencegah jadwal bertabrakan secara sistemik sebelum data tersimpan ke database.
