-- Migration Down: 000005_create_academic_years_and_enrollments.down.sql
-- Description: Rollback tabel class_enrollments dan academic_years

DROP TABLE IF EXISTS class_enrollments CASCADE;
DROP TABLE IF EXISTS academic_years CASCADE;
