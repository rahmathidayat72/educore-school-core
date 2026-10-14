# Security Architecture & Risk Assessment — EduCore School

## 1. Metadata Keamanan
- **Proyek:** EduCore School — School Core Management
- **Versi:** 1.0.0
- **Tanggal Tinjauan:** 14 Oktober 2026
- **Status Security Gate:** APPROVED (Gerbang 2)
- **Security Risk Tier:** **`MEDIUM`**
- **Candidate Security Owner:** Tech Lead / Lead Developer (Rahmat Hidayat)

---

## 2. Justifikasi Security Risk Tier (MEDIUM)
1. **Pemrosesan Data PII:** Sistem mengelola data identitas peserta didik (nama lengkap, NISN, NIK orang tua, kontak, alamat rumah) yang tunduk pada kepatuhan pelindungan data pribadi.
2. **Kredensial dan Multi-Role:** Mengatur hak akses multi-role berjenjang (Super Admin Yayasan, Admin Unit, Guru, Siswa, Wali Murid, Alumni) dengan token JWT.
3. **Multi-Unit Yayasan:** Kegagalan isolasi data dapat menyebabkan kebocoran informasi antar-sekolah di bawah yayasan yang sama.
4. **Audit Trail Finansial/Operasional:** Menjadi dasar verifikasi data untuk aplikasi ekosistem lainnya (keuangan SPP, presensi, rapor).

---

## 3. Matriks Cakupan 12 Kategori Ancaman (Threat Coverage Matrix)

| No | Kategori Ancaman | Status | Rincian Mitigasi / Kontrol Pencegahan & Deteksi |
|:---:|---|:---:|---|
| 1 | **Reconnaissance** | `APPLICABLE` | Menyembunyikan banner version pada response header; menonaktifkan directory listing web server; menyembunyikan stack trace error internal di production. |
| 2 | **Scanning & Enumeration** | `APPLICABLE` | Endpoint `/login` mengembalikan pesan generik (*"Email atau password salah"*) untuk mencegah penentuan keberadaan email pengguna (*user enumeration*). |
| 3 | **Vulnerability Assessment** | `APPLICABLE` | Melakukan scan dependensi berkala: `govulncheck ./...` pada Go backend dan `npm audit` pada frontend Vue 3. |
| 4 | **Exploitation (SQLi / XSS)** | `APPLICABLE` | 100% query PostgreSQL menggunakan *parameterized queries* via driver `pgx/v5` (mencegah SQL Injection). Frontend Vue 3 me-render string dengan auto-escaping (mencegah XSS). |
| 5 | **Password Attacks** | `APPLICABLE` | Password pengguna di-hash menggunakan algoritma `bcrypt` dengan cost minimum 10. Memasang rate limiting pada endpoint autentikasi. |
| 6 | **Social Engineering** | `NOT APPLICABLE` | Modul Core tidak memicu pengiriman SMS/OTP keluar langsung pada fase MVP; perubahan kredensial dilakukan melalui verifikasi staf admin unit resmi. |
| 7 | **Web / API Attacks** | `APPLICABLE` | Pengaturan CORS diperketat dengan whitelist origin terdaftar; validasi otorisasi `Bearer <JWT>` pada seluruh endpoint privat; CSRF mitigation melalui stateless header. |
| 8 | **Privilege Escalation** | `APPLICABLE` | Otorisasi role diverifikasi di sisi backend (`RolesAllowedMiddleware`), bukan sekadar pembatasan tampilan navigasi frontend. |
| 9 | **Lateral Movement (Cross-Tenant)** | `APPLICABLE` | Isolasi data unit sekolah: seluruh query data otomatis menyertakan filter `school_unit_id` dari JWT token pengguna untuk mencegah intervensi antar-unit. |
| 10 | **Persistence** | `APPLICABLE` | Token JWT memiliki masa kedaluwarsa terbatas (*expiry*); sistem memeriksa status aktif akun pengguna (`delete_at IS NULL`) pada setiap request. |
| 11 | **Defense Evasion** | `APPLICABLE` | Log transaksi dicatat di tabel terpisah (`transaction_logs`) yang berstatus *append-only* dan tidak dapat diubah oleh operator sekolah biasa. |
| 12 | **Exfiltration & Data Abuse** | `APPLICABLE` | Masking otomatis nilai kredensial sensitif sebelum disimpan ke log; pembatasan paginasi (limit & offset) pada query data massal. |

---

## 4. Standar Penanganan Kredensial & Secrets
- Tidak boleh ada secret, token API, atau password yang di-hardcode di kode program.
- Variabel konfigurasi lingkungan wajib dibaca dari environment variables (`.env` lokal tidak boleh di-commit ke Git).
- Komunikasi produksi wajib menggunakan protokol terenkripsi (HTTPS / TLS 1.3).

---

## 5. Rencana Tanggap Insiden Keamanan (Incident Response)
1. **Deteksi:** Pemantauan anomali log HTTP (banyaknya status 401/403 berturut-turut pada audit log).
2. **Isolasi:** Admin dapat menonaktifkan akun yang dicurigai disusupi secara langsung (`delete_at`).
3. **Pemberitahuan:** Laporkan temuan ke Lead Developer untuk analisis akar masalah (*root cause analysis*).
4. **Perbaikan & Verifikasi:** Rilis patch keamanan terverifikasi melalui prosedur pengujian regresi.
