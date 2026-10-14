-- Migration Down: 000004_enhance_kelas_multi_jenjang.down.sql
-- Description: Menghapus indeks performa untuk atribut jenjang pada tabel kelas

DROP INDEX IF EXISTS idx_kelas_unit_tingkat;
DROP INDEX IF EXISTS idx_kelas_jurusan;
DROP INDEX IF EXISTS idx_kelas_fase;
DROP INDEX IF EXISTS idx_kelas_tingkat;
