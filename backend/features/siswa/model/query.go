package model

import (
	"context"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/siswa"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type siswaQuery struct {
	db *pgxpool.Pool
}

func NewSiswaData(db *pgxpool.Pool) siswa.DataSiswaInterface {
	if db == nil {
		panic("Nil database")
	}
	return &siswaQuery{db: db}
}

func toNullString(val string) interface{} {
	trimmed := strings.TrimSpace(val)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

func toNullDate(val string) (interface{}, error) {
	trimmed := strings.TrimSpace(val)
	if trimmed == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", trimmed)
	if err != nil {
		return nil, fmt.Errorf("format tanggal harus YYYY-MM-DD: %w", err)
	}
	return t, nil
}

func (s *siswaQuery) InsertSiswa(insert *siswa.SiswaCore) error {
	if s.db == nil {
		return errors.New("Nil database")
	}

	if insert == nil {
		return errors.New("insert data is nil")
	}

	if insert.ID == "" {
		insert.ID = uuid.New().String()
	}

	// Sinkronisasi Nama_Kelas & Kelas_ID
	switch {
	case insert.Nama_Kelas != "" && insert.Kelas_ID == "":
		insert.Nama_Kelas = strings.TrimSpace(insert.Nama_Kelas)
		var kelasID string
		err := s.db.QueryRow(context.Background(),
			"SELECT id FROM kelas WHERE TRIM(kelas) ILIKE TRIM($1)", insert.Nama_Kelas).Scan(&kelasID)
		if err != nil {
			log.Printf("InsertKelas: nama kelas '%s' tidak ditemukan", insert.Nama_Kelas)
			return fmt.Errorf("kelas dengan nama '%s' tidak ditemukan", insert.Nama_Kelas)
		}
		insert.Kelas_ID = kelasID

	case insert.Kelas_ID != "" && insert.Nama_Kelas == "":
		var namaKelas string
		err := s.db.QueryRow(context.Background(),
			"SELECT kelas FROM kelas WHERE id = $1", insert.Kelas_ID).Scan(&namaKelas)
		if err != nil {
			log.Printf("InsertKelas: ID kelas '%s' tidak ditemukan", insert.Kelas_ID)
			return fmt.Errorf("kelas dengan ID '%s' tidak ditemukan", insert.Kelas_ID)
		}
		insert.Nama_Kelas = namaKelas

	case insert.Kelas_ID != "" && insert.Nama_Kelas != "":
		insert.Nama_Kelas = strings.TrimSpace(insert.Nama_Kelas)
		var existingName string
		err := s.db.QueryRow(context.Background(),
			"SELECT kelas FROM kelas WHERE id = $1", insert.Kelas_ID).Scan(&existingName)
		if err != nil {
			log.Printf("InsertKelas: ID kelas '%s' tidak ditemukan", insert.Kelas_ID)
			return fmt.Errorf("kelas dengan ID '%s' tidak ditemukan", insert.Kelas_ID)
		}
		if strings.TrimSpace(strings.ToLower(existingName)) != strings.ToLower(insert.Nama_Kelas) {
			log.Printf("InsertKelas: Nama kelas tidak cocok dengan ID kelas. Dapat: '%s', seharusnya: '%s'",
				insert.Nama_Kelas, existingName)
			return fmt.Errorf("nama kelas '%s' tidak cocok dengan ID kelas '%s'", insert.Nama_Kelas, insert.Kelas_ID)
		}
	}

	var idKelasParam interface{}
	if strings.TrimSpace(insert.Kelas_ID) == "" {
		idKelasParam = nil
	} else {
		idKelasParam = strings.TrimSpace(insert.Kelas_ID)
	}

	tglLahirParam, err := toNullDate(insert.TanggalLahir)
	if err != nil {
		return err
	}

	statusVal := strings.TrimSpace(insert.Status)
	if statusVal == "" {
		statusVal = "Aktif"
	}

	query := `INSERT INTO siswa (
		id, kelas_id, nama, nis, nisn, jenis_kelamin, tempat_lahir, tanggal_lahir,
		agama, no_telepon, nama_wali, no_telepon_wali, tahun_masuk, status,
		email, alamat, school_unit_id, update_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8,
		$9, $10, $11, $12, $13, $14,
		$15, $16, $17, CURRENT_TIMESTAMP
	)`

	var unitIDParam interface{}
	if insert.SchoolUnitID != nil && strings.TrimSpace(*insert.SchoolUnitID) != "" {
		unitIDParam = strings.TrimSpace(*insert.SchoolUnitID)
	}

	_, err = s.db.Exec(context.Background(), query,
		insert.ID,
		idKelasParam,
		insert.Nama,
		toNullString(insert.NIS),
		toNullString(insert.NISN),
		toNullString(insert.JenisKelamin),
		toNullString(insert.TempatLahir),
		tglLahirParam,
		toNullString(insert.Agama),
		toNullString(insert.NoTelepon),
		toNullString(insert.NamaWali),
		toNullString(insert.NoTeleponWali),
		toNullString(insert.TahunMasuk),
		statusVal,
		toNullString(insert.Email),
		toNullString(insert.Alamat),
		unitIDParam,
	)
	if err != nil {
		log.Printf("InsertSiswa error exec: %v", err)
		if strings.Contains(err.Error(), "idx_siswa_nis_unique") {
			return errors.New("NIS sudah terdaftar untuk siswa lain")
		}
		if strings.Contains(err.Error(), "idx_siswa_nisn_unique") {
			return errors.New("NISN sudah terdaftar untuk siswa lain")
		}
		if strings.Contains(err.Error(), "siswa_email_key") {
			return errors.New("email siswa sudah digunakan")
		}
		return fmt.Errorf("insert failed: %w", err)
	}

	return nil
}

func (s *siswaQuery) SelectAllSiswa(schoolUnitID ...string) ([]siswa.SiswaCore, error) {
	if s.db == nil {
		return nil, errors.New("Nil database")
	}

	query := `SELECT 
		s.id, 
		s.kelas_id, 
		k.kelas AS nama_kelas,
		s.nama, 
		s.nis,
		s.nisn,
		s.jenis_kelamin,
		s.tempat_lahir,
		s.tanggal_lahir,
		s.agama,
		s.no_telepon,
		s.nama_wali,
		s.no_telepon_wali,
		s.tahun_masuk,
		s.status,
		s.email, 
		s.alamat,
		s.school_unit_id
	FROM 
		siswa s
	LEFT JOIN 
		kelas k ON s.kelas_id = k.id
	WHERE 
		s.delete_at IS NULL`

	var args []interface{}
	if len(schoolUnitID) > 0 && strings.TrimSpace(schoolUnitID[0]) != "" {
		query += " AND s.school_unit_id = $1"
		args = append(args, strings.TrimSpace(schoolUnitID[0]))
	}

	query += " ORDER BY s.nama ASC"

	rows, err := s.db.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("SelectAllSiswa error query: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}
	defer rows.Close()

	var result []siswa.SiswaCore

	for rows.Next() {
		var sw Siswa

		err := rows.Scan(
			&sw.ID,
			&sw.Kelas_ID,
			&sw.Nama_Kelas,
			&sw.Nama,
			&sw.NIS,
			&sw.NISN,
			&sw.JenisKelamin,
			&sw.TempatLahir,
			&sw.TanggalLahir,
			&sw.Agama,
			&sw.NoTelepon,
			&sw.NamaWali,
			&sw.NoTeleponWali,
			&sw.TahunMasuk,
			&sw.Status,
			&sw.Email,
			&sw.Alamat,
			&sw.SchoolUnitID,
		)
		if err != nil {
			log.Printf("SelectAllSiswa error scan: %v", err)
			return nil, fmt.Errorf("select failed: %w", err)
		}

		core := FormatterResponse(sw)
		result = append(result, core)
	}

	if err := rows.Err(); err != nil {
		log.Printf("SelectAll error rows: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}

	log.Printf("Successfully fetched %d siswa from database", len(result))
	return result, nil
}

func (s *siswaQuery) SelectById(id string, schoolUnitID ...string) (*siswa.SiswaCore, error) {
	if s == nil || s.db == nil {
		return nil, errors.New("Nil database")
	}
	if id == "" {
		return nil, errors.New("ID cannot be empty")
	}

	query := `SELECT 
		s.id, 
		s.kelas_id, 
		k.kelas AS nama_kelas, 
		s.nama, 
		s.nis,
		s.nisn,
		s.jenis_kelamin,
		s.tempat_lahir,
		s.tanggal_lahir,
		s.agama,
		s.no_telepon,
		s.nama_wali,
		s.no_telepon_wali,
		s.tahun_masuk,
		s.status,
		s.email, 
		s.alamat,
		s.school_unit_id
	FROM siswa s 
	LEFT JOIN kelas k ON s.kelas_id = k.id 
	WHERE s.id = $1 AND s.delete_at IS NULL`

	args := []interface{}{id}
	if len(schoolUnitID) > 0 && strings.TrimSpace(schoolUnitID[0]) != "" {
		query += " AND s.school_unit_id = $2"
		args = append(args, strings.TrimSpace(schoolUnitID[0]))
	}

	row := s.db.QueryRow(context.Background(), query, args...)

	var sw Siswa
	err := row.Scan(
		&sw.ID,
		&sw.Kelas_ID,
		&sw.Nama_Kelas,
		&sw.Nama,
		&sw.NIS,
		&sw.NISN,
		&sw.JenisKelamin,
		&sw.TempatLahir,
		&sw.TanggalLahir,
		&sw.Agama,
		&sw.NoTelepon,
		&sw.NamaWali,
		&sw.NoTeleponWali,
		&sw.TahunMasuk,
		&sw.Status,
		&sw.Email,
		&sw.Alamat,
		&sw.SchoolUnitID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		log.Printf("SelectById error scan: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}

	result := FormatterResponse(sw)
	return &result, nil
}

func (s *siswaQuery) Update(insert *siswa.SiswaCore, id string, schoolUnitID ...string) error {
	if s == nil || s.db == nil {
		return errors.New("Nil database")
	}
	if id == "" {
		return errors.New("ID tidak boleh kosong")
	}

	tglLahirParam, err := toNullDate(insert.TanggalLahir)
	if err != nil {
		return err
	}

	statusVal := strings.TrimSpace(insert.Status)
	if statusVal == "" {
		statusVal = "Aktif"
	}

	var idKelasParam interface{}
	if strings.TrimSpace(insert.Kelas_ID) == "" {
		idKelasParam = nil
	} else {
		idKelasParam = strings.TrimSpace(insert.Kelas_ID)
	}

	query := `UPDATE siswa SET 
		nama = $1, 
		email = $2, 
		alamat = $3, 
		kelas_id = $4,
		nis = $5,
		nisn = $6,
		jenis_kelamin = $7,
		tempat_lahir = $8,
		tanggal_lahir = $9,
		agama = $10,
		no_telepon = $11,
		nama_wali = $12,
		no_telepon_wali = $13,
		tahun_masuk = $14,
		status = $15,
		update_at = CURRENT_TIMESTAMP
	WHERE id = $16 AND delete_at IS NULL`

	args := []interface{}{
		insert.Nama,
		toNullString(insert.Email),
		toNullString(insert.Alamat),
		idKelasParam,
		toNullString(insert.NIS),
		toNullString(insert.NISN),
		toNullString(insert.JenisKelamin),
		toNullString(insert.TempatLahir),
		tglLahirParam,
		toNullString(insert.Agama),
		toNullString(insert.NoTelepon),
		toNullString(insert.NamaWali),
		toNullString(insert.NoTeleponWali),
		toNullString(insert.TahunMasuk),
		statusVal,
		id,
	}

	if len(schoolUnitID) > 0 && strings.TrimSpace(schoolUnitID[0]) != "" {
		query += fmt.Sprintf(" AND school_unit_id = $%d", len(args)+1)
		args = append(args, strings.TrimSpace(schoolUnitID[0]))
	}

	res, err := s.db.Exec(context.Background(), query, args...)
	if err != nil {
		log.Printf("Update error exec: %v", err)
		if strings.Contains(err.Error(), "idx_siswa_nis_unique") {
			return errors.New("NIS sudah terdaftar untuk siswa lain")
		}
		if strings.Contains(err.Error(), "idx_siswa_nisn_unique") {
			return errors.New("NISN sudah terdaftar untuk siswa lain")
		}
		if strings.Contains(err.Error(), "siswa_email_key") {
			return errors.New("email siswa sudah digunakan")
		}
		return fmt.Errorf("update failed: %w", err)
	}

	if res.RowsAffected() == 0 {
		log.Printf("UpdateSiswa: no rows updated for id %s", id)
		return errors.New("update failed: no rows affected")
	}

	return nil
}

func (s *siswaQuery) DeleteById(id string, schoolUnitID ...string) error {
	if s.db == nil {
		return errors.New("Koneksi database tidak ada")
	}

	query := "UPDATE siswa SET delete_at = NOW() WHERE id = $1 AND delete_at IS NULL"
	args := []interface{}{id}
	if len(schoolUnitID) > 0 && strings.TrimSpace(schoolUnitID[0]) != "" {
		query += " AND school_unit_id = $2"
		args = append(args, strings.TrimSpace(schoolUnitID[0]))
	}

	res, err := s.db.Exec(context.Background(), query, args...)
	if err != nil {
		log.Printf("DeleteById error exec: %v", err)
		return fmt.Errorf("hapus gagal: %w", err)
	}

	if res.RowsAffected() == 0 {
		log.Printf("DeleteById: tidak ada baris yang dihapus untuk id %s", id)
		return errors.New("hapus gagal: tidak ada baris yang terpengaruh")
	}

	return nil
}
