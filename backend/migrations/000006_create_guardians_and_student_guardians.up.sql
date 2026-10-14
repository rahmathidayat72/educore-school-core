-- Migration Up: 000006_create_guardians_and_student_guardians.up.sql
-- Description: Membuat tabel master guardians dan relasi student_guardians (many-to-many)

-- 1. Tabel Master Wali Murid / Orang Tua
CREATE TABLE IF NOT EXISTS guardians (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NULL,
    nama VARCHAR(150) NOT NULL,
    nik VARCHAR(20) NULL,
    hubungan VARCHAR(30) NOT NULL DEFAULT 'Wali', -- 'Ayah', 'Ibu', 'Wali', 'Kakak'
    pekerjaan VARCHAR(100) NULL,
    penghasilan VARCHAR(50) NULL,
    telepon VARCHAR(30) NOT NULL,
    email VARCHAR(100) NULL,
    alamat TEXT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP WITH TIME ZONE NULL,
    CONSTRAINT fk_guardians_user 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Index pencarian dan filter
CREATE INDEX IF NOT EXISTS idx_guardians_telepon ON guardians (telepon);
CREATE INDEX IF NOT EXISTS idx_guardians_nik ON guardians (nik);
CREATE INDEX IF NOT EXISTS idx_guardians_user_id ON guardians (user_id);
CREATE INDEX IF NOT EXISTS idx_guardians_delete_at ON guardians (delete_at);

-- 2. Tabel Relasi Siswa & Wali Murid (Many-to-Many)
CREATE TABLE IF NOT EXISTS student_guardians (
    id VARCHAR(36) PRIMARY KEY,
    guardian_id VARCHAR(36) NOT NULL,
    siswa_id VARCHAR(36) NOT NULL,
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP WITH TIME ZONE NULL,
    CONSTRAINT fk_student_guardians_guardian 
        FOREIGN KEY (guardian_id) REFERENCES guardians(id) ON DELETE CASCADE,
    CONSTRAINT fk_student_guardians_siswa 
        FOREIGN KEY (siswa_id) REFERENCES siswa(id) ON DELETE CASCADE
);

-- Mencegah duplikasi tautan wali ke siswa yang sama
CREATE UNIQUE INDEX IF NOT EXISTS idx_student_guardians_unique
    ON student_guardians (guardian_id, siswa_id)
    WHERE delete_at IS NULL;

-- Index performa query pencarian
CREATE INDEX IF NOT EXISTS idx_student_guardians_guardian ON student_guardians (guardian_id);
CREATE INDEX IF NOT EXISTS idx_student_guardians_siswa ON student_guardians (siswa_id);
CREATE INDEX IF NOT EXISTS idx_student_guardians_delete_at ON student_guardians (delete_at);

-- 3. Backfill data dari siswa existing yang memiliki nama_wali
INSERT INTO guardians (id, nama, hubungan, telepon, created_at, update_at)
SELECT 
    'gdn-' || s.id,
    s.nama_wali,
    'Wali',
    COALESCE(NULLIF(s.no_telepon_wali, ''), '-'),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
FROM siswa s
WHERE s.nama_wali IS NOT NULL AND s.nama_wali != '' AND s.delete_at IS NULL
ON CONFLICT (id) DO NOTHING;

INSERT INTO student_guardians (id, guardian_id, siswa_id, is_primary, created_at, update_at)
SELECT 
    'sg-' || s.id,
    'gdn-' || s.id,
    s.id,
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
FROM siswa s
WHERE s.nama_wali IS NOT NULL AND s.nama_wali != '' AND s.delete_at IS NULL
ON CONFLICT (id) DO NOTHING;
