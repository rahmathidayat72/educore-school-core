# Standard Operating Procedure (SOP)
## EduCore School Engineering

## 1. Prosedur Awal Setiap Sesi
Sebelum mulai bekerja di setiap sesi, AI dan developer WAJIB:
1. Membaca `AGENTS.md` di root proyek.
2. Membaca `.agent/PROJECT-STATE.md` untuk mengetahui status dan task aktif.
3. Membaca `.agent/TASKS.md` dan `.agent/REPORT.md` terbaru.
4. Menjalankan pemeriksaan read-only terhadap Git status dan source code.
5. Meminta persetujuan Gerbang 3 sebelum menulis atau mengubah kode.

---

## 2. Aturan Coding dan Arsitektur
- **Go Backend:**
  - Format kode wajib menggunakan `gofmt` standar.
  - Jangan menambahkan framework HTTP berat (Gin/Fiber/Echo) jika `net/http` sudah mencukupi.
  - Seluruh query database wajib menggunakan parameterized queries (`$1, $2, ...`) via pgxpool. Dilarang konkatenasi string SQL.
  - Tangani error secara eksplisit; jangan abaikan error (`_ = err`).
  - Struktur domain harus konsisten: `entities.go`, `model/`, `service/`, `controllers/`.
- **Vue 3 Frontend:**
  - Gunakan Composition API (`<script setup>`).
  - Kelola autentikasi dan konteks pengguna di Pinia store (`src/stores/auth.js`).
  - Endpoint API dipusatkan di folder `src/api/` menggunakan instance Axios terpusat.
  - Lindungi rute privat dengan navigation guard berbasis role di `src/router/index.js`.

---

## 3. Aturan Database dan Migrasi
1. Dilarang melakukan perubahan DDL langsung di database production tanpa skrip migrasi yang teruji.
2. Setiap perubahan skema tabel wajib menyertakan:
   - Skrip migrasi maju (*UP*).
   - Skrip rollback (*DOWN*).
   - Penanganan data default / backfill jika kolom baru bersifat `NOT NULL`.
3. Terapkan soft delete (`delete_at TIMESTAMP`) untuk entitas data operasional utama (siswa, guru, kelas, mapel, jadwal).

---

## 4. Security Baseline
1. **Zero Secret in Code:** Dilarang commit file `.env`, kredensial database, atau JWT secret ke Git.
2. **Password:** Wajib di-hash menggunakan `bcrypt` dengan cost >= 10.
3. **CORS:** Batasi origin yang diizinkan sesuai lingkungan (bukan wildcard `*` di mode produksi).
4. **Isolasi Unit:** Seluruh query mutasi data wajib memvalidasi kecocokan `school_unit_id` pengguna.

---

## 5. Testing & Verification Commands
- **Backend Test:**
  ```bash
  cd backend && go test ./...
  cd backend && go test -cover ./...
  ```
- **Frontend Build & Lint:**
  ```bash
  cd frontend && npm run build
  ```

---

## 6. Definition of Ready (DoR) dan Definition of Done (DoD)

### Definition of Ready (DoR):
- [ ] Task memiliki ID dan requirement terkait (BUS/FR).
- [ ] Scope dan out-of-scope dijelaskan eksplisit.
- [ ] File dan simbol yang terdampak telah dipetakan.
- [ ] Rencana pengujian dan perintah eksekusi terklasifikasi jelas.
- [ ] Task dan code plan disetujui pengguna (Gerbang 3).

### Definition of Done (DoD):
- [ ] Implementasi sesuai code plan yang disetujui.
- [ ] Unit test dan pengujian regresi lulus (*PASS*).
- [ ] Frontend build sukses tanpa error.
- [ ] Dokumentasi dan tracking (`.agent/TASKS.md`, `.agent/REPORT.md`, `.agent/PROJECT-STATE.md`) telah diperbarui.
