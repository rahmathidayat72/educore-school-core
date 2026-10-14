# 🏫 Aplikasi Sistem Informasi Sekolah (Fullstack)

Proyek ini terdiri dari dua bagian:
- **`backend/`**: RESTful API Go Native + PostgreSQL dengan autentikasi JWT dan CORS middleware.
- **`frontend/`**: Single Page Application (SPA) berbasis Vue 3 + Vite + Bootstrap 5 + Pinia + Vue Router.

---

## 📁 Struktur Direktori

```text
app_sekolah/
├── backend/          # Server Go Native (Port 8080)
│   ├── config/       # Database & Env Loader
│   ├── features/     # Modul Auth, Users, Guru, Kelas, Siswa, Mapel
│   ├── helper/       # Middleware JWT, CORS, Logging, Response
│   ├── router/       # Router Endpoint
│   ├── bruno/        # Koleksi API Test untuk Bruno
│   ├── .env          # Konfigurasi Environment DB & JWT
│   └── main.go       # Entry Point
│
└── frontend/         # Web Portal & Dashboard Vue 3 (Port 5173)
    ├── src/
    │   ├── api/      # Axios Client & Service Modul
    │   ├── stores/   # Pinia Auth Store
    │   ├── router/   # Vue Router & Navigation Guards
    │   ├── views/    # Halaman Publik, Login, & Dashboard Admin
    │   └── components/
    └── package.json
```

---

## 🚀 Cara Menjalankan

### 1. Jalankan Backend (Terminal 1)
```bash
cd /home/rahmath/Documents/rahmath/app_sekolah/backend
go run main.go
```
*Backend akan berjalan di `http://localhost:8080`*

### 2. Jalankan Frontend (Terminal 2)
```bash
cd /home/rahmath/Documents/rahmath/app_sekolah/frontend
npm run dev
```
*Frontend akan berjalan di `http://localhost:5173`*

---

## 🔐 Akun Default untuk Login

- **Email**: `admin@sekolah.id`
- **Password**: `adminpassword`
- **Role**: `admin`
