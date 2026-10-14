-- Migration Up: 000003_seed_default_units_and_backfill.up.sql
-- Description: Mendaftarkan unit sekolah percontohan, backfill data existing, dan memasang Foreign Key constraint

-- 1. Seed unit sekolah yayasan percontohan (SD, SMP, SMA, SMK)
INSERT INTO school_units (id, nama_sekolah, jenjang, npsn, alamat, telepon, email, status)
VALUES 
    ('unit-sma-001', 'SMA EduCore 1', 'SMA', '20109901', 'Jl. Pendidikan Utama No. 10', '021-5550101', 'sma1@educore.sch.id', 'aktif'),
    ('unit-smp-001', 'SMP EduCore 1', 'SMP', '20109902', 'Jl. Pendidikan Utama No. 12', '021-5550102', 'smp1@educore.sch.id', 'aktif'),
    ('unit-smk-001', 'SMK EduCore Teknologi', 'SMK', '20109903', 'Jl. Industri EduCore No. 5', '021-5550103', 'smk@educore.sch.id', 'aktif'),
    ('unit-sd-001', 'SD EduCore Unggulan', 'SD', '20109904', 'Jl. Melati Harapan No. 3', '021-5550104', 'sd@educore.sch.id', 'aktif')
ON CONFLICT (id) DO NOTHING;

-- 2. Backfill data yang sudah ada ke unit pilot utama ('unit-sma-001')
UPDATE guru SET school_unit_id = 'unit-sma-001' WHERE school_unit_id IS NULL;
UPDATE kelas SET school_unit_id = 'unit-sma-001' WHERE school_unit_id IS NULL;
UPDATE siswa SET school_unit_id = 'unit-sma-001' WHERE school_unit_id IS NULL;
UPDATE mata_pelajaran SET school_unit_id = 'unit-sma-001' WHERE school_unit_id IS NULL;
UPDATE jadwal_pelajaran SET school_unit_id = 'unit-sma-001' WHERE school_unit_id IS NULL;
UPDATE transaction_logs SET school_unit_id = 'unit-sma-001' WHERE school_unit_id IS NULL;
UPDATE users SET school_unit_id = 'unit-sma-001' WHERE school_unit_id IS NULL AND role != 'super_admin';

-- 3. Memasang Foreign Key Constraint ke school_units(id)
ALTER TABLE users 
    ADD CONSTRAINT fk_users_school_unit 
    FOREIGN KEY (school_unit_id) REFERENCES school_units(id) ON DELETE RESTRICT;

ALTER TABLE guru 
    ADD CONSTRAINT fk_guru_school_unit 
    FOREIGN KEY (school_unit_id) REFERENCES school_units(id) ON DELETE RESTRICT;

ALTER TABLE kelas 
    ADD CONSTRAINT fk_kelas_school_unit 
    FOREIGN KEY (school_unit_id) REFERENCES school_units(id) ON DELETE RESTRICT;

ALTER TABLE siswa 
    ADD CONSTRAINT fk_siswa_school_unit 
    FOREIGN KEY (school_unit_id) REFERENCES school_units(id) ON DELETE RESTRICT;

ALTER TABLE mata_pelajaran 
    ADD CONSTRAINT fk_mata_pelajaran_school_unit 
    FOREIGN KEY (school_unit_id) REFERENCES school_units(id) ON DELETE RESTRICT;

ALTER TABLE jadwal_pelajaran 
    ADD CONSTRAINT fk_jadwal_pelajaran_school_unit 
    FOREIGN KEY (school_unit_id) REFERENCES school_units(id) ON DELETE RESTRICT;

ALTER TABLE transaction_logs 
    ADD CONSTRAINT fk_transaction_logs_school_unit 
    FOREIGN KEY (school_unit_id) REFERENCES school_units(id) ON DELETE SET NULL;
