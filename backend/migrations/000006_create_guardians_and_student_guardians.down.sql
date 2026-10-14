-- Migration Down: 000006_create_guardians_and_student_guardians.down.sql
-- Description: Rollback tabel student_guardians dan guardians

DROP TABLE IF EXISTS student_guardians CASCADE;
DROP TABLE IF EXISTS guardians CASCADE;
