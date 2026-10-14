-- Migration Up: 000005_create_academic_years_and_enrollments.up.sql
-- Description: Membuat tabel academic_years dan class_enrollments untuk riwayat penempatan rombel

-- 1. Tabel Tahun Ajaran
CREATE TABLE IF NOT EXISTS academic_years (
    id VARCHAR(36) PRIMARY KEY,
    school_unit_id VARCHAR(36) NOT NULL,
    tahun_ajaran VARCHAR(20) NOT NULL, -- Contoh: '2025/2026', '2026/2027'
    semester VARCHAR(10) NOT NULL,     -- 'Ganjil', 'Genap'
    is_active BOOLEAN DEFAULT false,
    tanggal_mulai DATE NULL,
    tanggal_selesai DATE NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP WITH TIME ZONE NULL,
    CONSTRAINT fk_academic_years_school_unit 
        FOREIGN KEY (school_unit_id) REFERENCES school_units(id) ON DELETE RESTRICT
);

-- Hanya boleh ada 1 tahun ajaran & semester yang aktif per unit sekolah
CREATE UNIQUE INDEX IF NOT EXISTS idx_academic_years_single_active
    ON academic_years (school_unit_id)
    WHERE is_active = true AND delete_at IS NULL;

-- Index performa query
CREATE INDEX IF NOT EXISTS idx_academic_years_unit_tahun 
    ON academic_years (school_unit_id, tahun_ajaran);
CREATE INDEX IF NOT EXISTS idx_academic_years_is_active 
    ON academic_years (is_active);
CREATE INDEX IF NOT EXISTS idx_academic_years_delete_at 
    ON academic_years (delete_at);

-- 2. Tabel Penempatan Kelas Siswa (Class Enrollments)
CREATE TABLE IF NOT EXISTS class_enrollments (
    id VARCHAR(36) PRIMARY KEY,
    school_unit_id VARCHAR(36) NOT NULL,
    academic_year_id VARCHAR(36) NOT NULL,
    kelas_id VARCHAR(36) NOT NULL,
    siswa_id VARCHAR(36) NOT NULL,
    status VARCHAR(20) DEFAULT 'Aktif', -- 'Aktif', 'Lulus', 'Pindah', 'Tinggal Kelas'
    keterangan TEXT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP WITH TIME ZONE NULL,
    CONSTRAINT fk_class_enrollments_school_unit 
        FOREIGN KEY (school_unit_id) REFERENCES school_units(id) ON DELETE RESTRICT,
    CONSTRAINT fk_class_enrollments_academic_year 
        FOREIGN KEY (academic_year_id) REFERENCES academic_years(id) ON DELETE RESTRICT,
    CONSTRAINT fk_class_enrollments_kelas 
        FOREIGN KEY (kelas_id) REFERENCES kelas(id) ON DELETE RESTRICT,
    CONSTRAINT fk_class_enrollments_siswa 
        FOREIGN KEY (siswa_id) REFERENCES siswa(id) ON DELETE RESTRICT
);

-- Mencegah pendaftaran ganda siswa di rombel berbeda pada tahun ajaran yang sama
CREATE UNIQUE INDEX IF NOT EXISTS idx_class_enrollments_unique_student_year
    ON class_enrollments (school_unit_id, academic_year_id, siswa_id)
    WHERE delete_at IS NULL;

-- Index performa query pencarian rombel & siswa
CREATE INDEX IF NOT EXISTS idx_class_enrollments_unit_year_kelas
    ON class_enrollments (school_unit_id, academic_year_id, kelas_id);
CREATE INDEX IF NOT EXISTS idx_class_enrollments_siswa
    ON class_enrollments (siswa_id);
CREATE INDEX IF NOT EXISTS idx_class_enrollments_status
    ON class_enrollments (status);
CREATE INDEX IF NOT EXISTS idx_class_enrollments_delete_at
    ON class_enrollments (delete_at);

-- 3. Seed data tahun ajaran percontohan untuk unit SMA pilot ('unit-sma-001')
INSERT INTO academic_years (id, school_unit_id, tahun_ajaran, semester, is_active, tanggal_mulai, tanggal_selesai)
VALUES 
    ('ay-sma-2026-ganjil', 'unit-sma-001', '2026/2027', 'Ganjil', true, '2026-07-15', '2026-12-20')
ON CONFLICT (id) DO NOTHING;

-- 4. Backfill data riwayat siswa existing ke class_enrollments
INSERT INTO class_enrollments (id, school_unit_id, academic_year_id, kelas_id, siswa_id, status)
SELECT 
    'enr-' || s.id,
    s.school_unit_id,
    'ay-sma-2026-ganjil',
    s.kelas_id,
    s.id,
    'Aktif'
FROM siswa s
WHERE s.school_unit_id = 'unit-sma-001' AND s.kelas_id IS NOT NULL AND s.delete_at IS NULL
ON CONFLICT (id) DO NOTHING;
