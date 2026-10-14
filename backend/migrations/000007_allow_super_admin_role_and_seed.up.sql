-- Update users role check constraint to support all EduCore roles
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role::text = ANY (ARRAY['super_admin'::text, 'admin'::text, 'guru'::text, 'wali'::text, 'siswa'::text, 'alumni'::text, 'user'::text]));

-- Seed default super_admin account for foundation-level management
INSERT INTO users (id, username, email, password, role, school_unit_id, update_at)
VALUES (
    'usr-superadmin-001',
    'superadmin',
    'superadmin@sekolah.id',
    '$2a$10$dROJaAgusWYhQg61zx31r.ZwdEI2oLBY8LXIEwg7fvSMfrqMW0Gpy',
    'super_admin',
    NULL,
    CURRENT_TIMESTAMP
) ON CONFLICT (email) DO UPDATE 
SET role = 'super_admin', school_unit_id = NULL;
