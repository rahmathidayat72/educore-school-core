DELETE FROM users WHERE id = 'usr-superadmin-001';
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role::text = ANY (ARRAY['admin'::character varying, 'user'::character varying, 'guru'::character varying]::text[]));
