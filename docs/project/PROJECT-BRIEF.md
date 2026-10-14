# Project Brief — EduCore School (School Core Management)

## 1. Metadata Proyek
- **Nama Proyek:** EduCore School
- **Modul Inti:** School Core Management
- **Versi:** 1.0.0
- **Status Approval:** Approved (Gerbang 1 & Gerbang 2 Selesai)
- **Tech Stack Utama:** Go Native (`net/http`) + PostgreSQL + Vue 3 (Vite + Pinia + Bootstrap 5)
- **Target Pengguna:** Super Admin Yayasan, Admin Unit Sekolah, Guru, Siswa, Orang Tua/Wali Murid, Alumni

---

## 2. Masalah dan Tujuan
- **Masalah:**
  - Pengelolaan data administrasi dan akademik sekolah yang terfragmentasi antar-unit (SD, SMP, SMA, SMK) dalam satu yayasan.
  - Duplikasi data (*data redundancy*) ketika mengintegrasikan aplikasi pendukung seperti presensi, ujian, e-rapor, dan keuangan.
  - Ketiadaan kontrol hierarki terpusat untuk yayasan dengan hak isolasi data per unit sekolah yang ketat.
- **Tujuan:**
  - Membangun pusat data induk (*Single Source of Truth*) untuk identitas dan referensi akademik yayasan multi-sekolah.
  - Menerapkan arsitektur multi-unit dengan isolasi data terjamin di tingkat database dan backend.
  - Membuka kontrak REST API terstandarisasi untuk integrasi modular dengan aplikasi ekosistem satelit.

---

## 3. Ruang Lingkup (Scope)
- **In Scope (Fase 1 - MVP Core):**
  - Manajemen multi-unit sekolah yayasan (`school_units`).
  - Autentikasi JWT dengan role berjenjang (`super_admin`, `admin`, `guru`, `wali`, `siswa`, `alumni`).
  - Master data siswa, guru, rombel/kelas (lintas jenjang SD/SMP/SMA/SMK), dan mata pelajaran.
  - Penjadwalan pelajaran tatap muka dengan validasi bentrok guru dan ruang.
  - Audit trail dan log mutasi data terpusat.
  - Kontrak REST API terstandarisasi untuk modul ekosistem.
- **Out of Scope (Fase Lanjutan):**
  - Modul Smart Attendance (QR/RFID).
  - Parent Portal mandiri (PWA).
  - School Finance & Payment Gateway.
  - CBT & LMS online.
  - E-Rapor & Asesmen detail.

---

## 4. Indikator Keberhasilan (Success Indicators)
1. Super Admin yayasan dapat memantau seluruh unit sekolah dan beralih konteks unit secara instan.
2. Admin unit sekolah hanya dapat melihat dan mengelola data milik unitnya sendiri (100% isolasi data).
3. Bentrok jadwal guru atau ruangan terdeteksi dan dicegah otomatis saat penyimpanan jadwal.
4. Seluruh endpoint lolos uji integritas API dan pengujian unit dengan cakupan >85%.
5. Waktu respon API rata-rata di bawah 250 ms untuk operasional harian.

---

## 5. Stakeholder & Kontak
- **Sponsor / Yayasan:** Badan Penyelenggara Pendidikan Yayasan.
- **Product Owner & Lead Developer:** Rahmat Hidayat.
- **Target Operator:** Staf Tata Usaha & Tim Kurikulum Unit Sekolah.
