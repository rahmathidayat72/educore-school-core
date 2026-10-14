-- Migration Down: 000003_seed_default_units_and_backfill.down.sql
-- Description: Rollback foreign key constraints dan seed units

-- 1. Hapus Foreign Key Constraints
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_school_unit;
ALTER TABLE guru DROP CONSTRAINT IF EXISTS fk_guru_school_unit;
ALTER TABLE kelas DROP CONSTRAINT IF EXISTS fk_kelas_school_unit;
ALTER TABLE siswa DROP CONSTRAINT IF EXISTS fk_siswa_school_unit;
ALTER TABLE mata_pelajaran DROP CONSTRAINT IF EXISTS fk_mata_pelajaran_school_unit;
ALTER TABLE jadwal_pelajaran DROP CONSTRAINT IF EXISTS fk_jadwal_pelajaran_school_unit;
ALTER TABLE transaction_logs DROP CONSTRAINT IF EXISTS fk_transaction_logs_school_unit;

-- 2. Kosongkan kembali kolom school_unit_id pada unit pilot jika rollback
UPDATE users SET school_unit_id = NULL WHERE school_unit_id = 'unit-sma-001';
UPDATE guru SET school_unit_id = NULL WHERE school_unit_id = 'unit-sma-001';
UPDATE kelas SET school_unit_id = NULL WHERE school_unit_id = 'unit-sma-001';
UPDATE siswa SET school_unit_id = NULL WHERE school_unit_id = 'unit-sma-001';
UPDATE mata_pelajaran SET school_unit_id = NULL WHERE school_unit_id = 'unit-sma-001';
UPDATE jadwal_pelajaran SET school_unit_id = NULL WHERE school_unit_id = 'unit-sma-001';
UPDATE transaction_logs SET school_unit_id = NULL WHERE school_unit_id = 'unit-sma-001';

-- 3. Hapus data unit percontohan
DELETE FROM school_units WHERE id IN ('unit-sma-001', 'unit-smp-001', 'unit-smk-001', 'unit-sd-001');
