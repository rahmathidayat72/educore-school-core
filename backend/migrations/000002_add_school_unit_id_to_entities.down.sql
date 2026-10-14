-- Migration Down: 000002_add_school_unit_id_to_entities.down.sql
-- Description: Rollback penambahan kolom school_unit_id dan index terkait

-- 1. Tabel users
DROP INDEX IF EXISTS idx_users_school_unit_id;
ALTER TABLE users DROP COLUMN IF EXISTS school_unit_id;

-- 2. Tabel guru
DROP INDEX IF EXISTS idx_guru_school_unit_id;
ALTER TABLE guru DROP COLUMN IF EXISTS school_unit_id;

-- 3. Tabel kelas
DROP INDEX IF EXISTS idx_kelas_school_unit_id;
ALTER TABLE kelas DROP COLUMN IF EXISTS school_unit_id;

-- 4. Tabel siswa
DROP INDEX IF EXISTS idx_siswa_school_unit_id;
ALTER TABLE siswa DROP COLUMN IF EXISTS school_unit_id;

-- 5. Tabel mata_pelajaran
DROP INDEX IF EXISTS idx_mata_pelajaran_school_unit_id;
ALTER TABLE mata_pelajaran DROP COLUMN IF EXISTS school_unit_id;

-- 6. Tabel jadwal_pelajaran
DROP INDEX IF EXISTS idx_jadwal_pelajaran_school_unit_id;
ALTER TABLE jadwal_pelajaran DROP COLUMN IF EXISTS school_unit_id;

-- 7. Tabel transaction_logs
DROP INDEX IF EXISTS idx_transaction_logs_school_unit_id;
ALTER TABLE transaction_logs DROP COLUMN IF EXISTS school_unit_id;
