-- Migration Up: 000002_add_school_unit_id_to_entities.up.sql
-- Description: Menambahkan kolom school_unit_id nullable dan index pada entitas operasional

-- 1. Tabel users
ALTER TABLE users ADD COLUMN IF NOT EXISTS school_unit_id VARCHAR(36);
CREATE INDEX IF NOT EXISTS idx_users_school_unit_id ON users (school_unit_id);

-- 2. Tabel guru
ALTER TABLE guru ADD COLUMN IF NOT EXISTS school_unit_id VARCHAR(36);
CREATE INDEX IF NOT EXISTS idx_guru_school_unit_id ON guru (school_unit_id);

-- 3. Tabel kelas
ALTER TABLE kelas ADD COLUMN IF NOT EXISTS school_unit_id VARCHAR(36);
CREATE INDEX IF NOT EXISTS idx_kelas_school_unit_id ON kelas (school_unit_id);

-- 4. Tabel siswa
ALTER TABLE siswa ADD COLUMN IF NOT EXISTS school_unit_id VARCHAR(36);
CREATE INDEX IF NOT EXISTS idx_siswa_school_unit_id ON siswa (school_unit_id);

-- 5. Tabel mata_pelajaran
ALTER TABLE mata_pelajaran ADD COLUMN IF NOT EXISTS school_unit_id VARCHAR(36);
CREATE INDEX IF NOT EXISTS idx_mata_pelajaran_school_unit_id ON mata_pelajaran (school_unit_id);

-- 6. Tabel jadwal_pelajaran
ALTER TABLE jadwal_pelajaran ADD COLUMN IF NOT EXISTS school_unit_id VARCHAR(36);
CREATE INDEX IF NOT EXISTS idx_jadwal_pelajaran_school_unit_id ON jadwal_pelajaran (school_unit_id);

-- 7. Tabel transaction_logs
ALTER TABLE transaction_logs ADD COLUMN IF NOT EXISTS school_unit_id VARCHAR(36);
CREATE INDEX IF NOT EXISTS idx_transaction_logs_school_unit_id ON transaction_logs (school_unit_id);
