-- Migration Up: 000004_enhance_kelas_multi_jenjang.up.sql
-- Description: Menambahkan indeks performa untuk atribut jenjang (tingkat, fase, jurusan) pada tabel kelas

CREATE INDEX IF NOT EXISTS idx_kelas_tingkat ON kelas (tingkat);
CREATE INDEX IF NOT EXISTS idx_kelas_fase ON kelas (fase);
CREATE INDEX IF NOT EXISTS idx_kelas_jurusan ON kelas (jurusan);
CREATE INDEX IF NOT EXISTS idx_kelas_unit_tingkat ON kelas (school_unit_id, tingkat);
