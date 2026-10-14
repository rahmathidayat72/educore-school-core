-- Migration Up: 000001_create_school_units_table.up.sql
-- Description: Membuat tabel master school_units untuk multi-unit yayasan

CREATE TABLE IF NOT EXISTS school_units (
    id VARCHAR(36) PRIMARY KEY,
    nama_sekolah VARCHAR(150) NOT NULL,
    jenjang VARCHAR(10) NOT NULL, -- 'SD', 'SMP', 'SMA', 'SMK'
    npsn VARCHAR(20),
    alamat TEXT,
    telepon VARCHAR(30),
    email VARCHAR(100),
    website VARCHAR(100),
    status VARCHAR(20) DEFAULT 'aktif', -- 'aktif', 'nonaktif', 'dalam_pembangunan'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP WITH TIME ZONE NULL
);

-- Index keunikan NPSN hanya untuk record yang belum di-soft delete
CREATE UNIQUE INDEX IF NOT EXISTS idx_school_units_npsn_unique 
    ON school_units (npsn) 
    WHERE delete_at IS NULL AND npsn IS NOT NULL;

-- Index performa query pencarian dan filter
CREATE INDEX IF NOT EXISTS idx_school_units_jenjang ON school_units (jenjang);
CREATE INDEX IF NOT EXISTS idx_school_units_status ON school_units (status);
CREATE INDEX IF NOT EXISTS idx_school_units_delete_at ON school_units (delete_at);
