-- Migration Down: 000001_create_school_units_table.down.sql
-- Description: Rollback tabel school_units

DROP TABLE IF EXISTS school_units CASCADE;
