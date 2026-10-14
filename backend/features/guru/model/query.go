package gurumodels

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/guru"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type guruQuery struct {
	db *pgxpool.Pool
}

func NewDataGuru(db *pgxpool.Pool) guru.DataGuruInterface {
	if db == nil {
		panic("guru model: Nil database")
	}
	return &guruQuery{db: db}
}

func (r *guruQuery) SelectAllGuru(schoolUnitID ...string) ([]guru.GuruCore, error) {
	if r.db == nil {
		return nil, errors.New("guru model: Nil database")
	}

	query := `SELECT 
		id, 
		id_user, 
		nama, 
		nip, 
		nuptk, 
		jenis_kelamin, 
		tempat_lahir, 
		tanggal_lahir, 
		no_telepon, 
		agama, 
		pendidikan_terakhir, 
		status_kepegawaian, 
		email, 
		alamat,
		school_unit_id
	FROM guru 
	WHERE delete_at IS NULL`

	var args []interface{}
	if len(schoolUnitID) > 0 && strings.TrimSpace(schoolUnitID[0]) != "" {
		query += " AND school_unit_id = $1"
		args = append(args, strings.TrimSpace(schoolUnitID[0]))
	}

	query += " ORDER BY nama ASC"

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []guru.GuruCore

	for rows.Next() {
		var g Guru
		var idUser sql.NullString

		err := rows.Scan(
			&g.ID,
			&idUser,
			&g.Nama,
			&g.NIP,
			&g.NUPTK,
			&g.JenisKelamin,
			&g.TempatLahir,
			&g.TanggalLahir,
			&g.NoTelepon,
			&g.Agama,
			&g.PendidikanTerakhir,
			&g.StatusKepegawaian,
			&g.Email,
			&g.Alamat,
			&g.SchoolUnitID,
		)
		if err != nil {
			log.Printf("SelectAll error scan: %v", err)
			return nil, fmt.Errorf("select failed: %w", err)
		}

		if idUser.Valid {
			g.ID_User = idUser.String
		} else {
			g.ID_User = ""
		}

		core := FormatterResponse(g)
		result = append(result, core)
	}

	if err := rows.Err(); err != nil {
		log.Printf("SelectAll error rows: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}

	log.Printf("Successfully fetched %d guru from database", len(result))
	return result, nil
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

func (r *guruQuery) InsertGuru(insert *guru.GuruCore) error {
	if r.db == nil {
		return errors.New("guru model: nil database connection")
	}

	if insert.ID == "" {
		insert.ID = uuid.New().String()
	}

	var idUserParam interface{}
	if strings.TrimSpace(insert.ID_User) == "" {
		idUserParam = nil
	} else {
		idUserParam = strings.TrimSpace(insert.ID_User)
	}

	tglLahirParam, err := toNullDate(insert.TanggalLahir)
	if err != nil {
		return err
	}

	query := `INSERT INTO guru (
		id, id_user, nama, nip, nuptk, jenis_kelamin, tempat_lahir, tanggal_lahir,
		no_telepon, agama, pendidikan_terakhir, status_kepegawaian, email, alamat,
		school_unit_id, update_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8,
		$9, $10, $11, $12, $13, $14,
		$15, CURRENT_TIMESTAMP
	)`

	var unitIDParam interface{}
	if insert.SchoolUnitID != nil && strings.TrimSpace(*insert.SchoolUnitID) != "" {
		unitIDParam = strings.TrimSpace(*insert.SchoolUnitID)
	}

	_, err = r.db.Exec(context.Background(), query,
		insert.ID,
		idUserParam,
		insert.Nama,
		toNullString(insert.NIP),
		toNullString(insert.NUPTK),
		toNullString(insert.JenisKelamin),
		toNullString(insert.TempatLahir),
		tglLahirParam,
		toNullString(insert.NoTelepon),
		toNullString(insert.Agama),
		toNullString(insert.PendidikanTerakhir),
		toNullString(insert.StatusKepegawaian),
		toNullString(insert.Email),
		toNullString(insert.Alamat),
		unitIDParam,
	)
	if err != nil {
		log.Printf("InsertGuru error exec: %v", err)
		if strings.Contains(err.Error(), "idx_guru_nip_unique") {
			return errors.New("NIP sudah terdaftar untuk guru lain")
		}
		if strings.Contains(err.Error(), "idx_guru_nuptk_unique") {
			return errors.New("NUPTK sudah terdaftar untuk guru lain")
		}
		if strings.Contains(err.Error(), "guru_email_key") {
			return errors.New("email guru sudah digunakan")
		}
		return fmt.Errorf("insert failed: %w", err)
	}

	return nil
}

func (r *guruQuery) Update(insert *guru.GuruCore, id string, schoolUnitID ...string) error {
	if r.db == nil {
		return errors.New("koneksi database nil")
	}
	if insert == nil {
		return errors.New("input data nil")
	}
	if id == "" {
		return errors.New("validation error: id harus diisi")
	}

	tglLahirParam, err := toNullDate(insert.TanggalLahir)
	if err != nil {
		return err
	}

	query := `UPDATE guru SET 
		nama = $2, 
		nip = $3,
		nuptk = $4,
		jenis_kelamin = $5,
		tempat_lahir = $6,
		tanggal_lahir = $7,
		no_telepon = $8,
		agama = $9,
		pendidikan_terakhir = $10,
		status_kepegawaian = $11,
		email = $12, 
		alamat = $13,
		update_at = CURRENT_TIMESTAMP
	WHERE id = $1 AND delete_at IS NULL`

	args := []interface{}{
		id,
		insert.Nama,
		toNullString(insert.NIP),
		toNullString(insert.NUPTK),
		toNullString(insert.JenisKelamin),
		toNullString(insert.TempatLahir),
		tglLahirParam,
		toNullString(insert.NoTelepon),
		toNullString(insert.Agama),
		toNullString(insert.PendidikanTerakhir),
		toNullString(insert.StatusKepegawaian),
		toNullString(insert.Email),
		toNullString(insert.Alamat),
	}

	if len(schoolUnitID) > 0 && strings.TrimSpace(schoolUnitID[0]) != "" {
		query += fmt.Sprintf(" AND school_unit_id = $%d", len(args)+1)
		args = append(args, strings.TrimSpace(schoolUnitID[0]))
	}

	res, err := r.db.Exec(context.Background(), query, args...)
	if err != nil {
		log.Printf("UpdateGuru error exec: %v", err)
		if strings.Contains(err.Error(), "idx_guru_nip_unique") {
			return errors.New("NIP sudah terdaftar untuk guru lain")
		}
		if strings.Contains(err.Error(), "idx_guru_nuptk_unique") {
			return errors.New("NUPTK sudah terdaftar untuk guru lain")
		}
		if strings.Contains(err.Error(), "guru_email_key") {
			return errors.New("email guru sudah digunakan")
		}
		return fmt.Errorf("update failed: %w", err)
	}

	if res.RowsAffected() == 0 {
		log.Printf("UpdateGuru: no rows updated for id %s", id)
		return errors.New("update failed: no rows affected")
	}

	return nil
}

func (r *guruQuery) SelectById(id string, schoolUnitID ...string) (*guru.GuruCore, error) {
	if r.db == nil {
		return nil, errors.New("guru query: koneksi database nil")
	}

	query := `SELECT 
		id, 
		id_user, 
		nama, 
		nip, 
		nuptk, 
		jenis_kelamin, 
		tempat_lahir, 
		tanggal_lahir, 
		no_telepon, 
		agama, 
		pendidikan_terakhir, 
		status_kepegawaian, 
		email, 
		alamat,
		school_unit_id
	FROM guru 
	WHERE id = $1 AND delete_at IS NULL`

	args := []interface{}{id}
	if len(schoolUnitID) > 0 && strings.TrimSpace(schoolUnitID[0]) != "" {
		query += " AND school_unit_id = $2"
		args = append(args, strings.TrimSpace(schoolUnitID[0]))
	}

	row := r.db.QueryRow(context.Background(), query, args...)

	var g Guru
	var idUser sql.NullString

	err := row.Scan(
		&g.ID,
		&idUser,
		&g.Nama,
		&g.NIP,
		&g.NUPTK,
		&g.JenisKelamin,
		&g.TempatLahir,
		&g.TanggalLahir,
		&g.NoTelepon,
		&g.Agama,
		&g.PendidikanTerakhir,
		&g.StatusKepegawaian,
		&g.Email,
		&g.Alamat,
		&g.SchoolUnitID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("gagal mengambil data guru: %w", err)
	}

	if idUser.Valid {
		g.ID_User = idUser.String
	} else {
		g.ID_User = ""
	}

	result := FormatterResponse(g)
	return &result, nil
}

func (r *guruQuery) DeleteById(id string, schoolUnitID ...string) error {
	if r.db == nil {
		return errors.New("guru query: koneksi database nil")
	}

	query := "UPDATE guru SET delete_at = NOW() WHERE id = $1 AND delete_at IS NULL"
	args := []interface{}{id}
	if len(schoolUnitID) > 0 && strings.TrimSpace(schoolUnitID[0]) != "" {
		query += " AND school_unit_id = $2"
		args = append(args, strings.TrimSpace(schoolUnitID[0]))
	}

	res, err := r.db.Exec(context.Background(), query, args...)
	if err != nil {
		log.Printf("DeleteById error exec: %v", err)
		return fmt.Errorf("delete failed: %w", err)
	}

	if res.RowsAffected() == 0 {
		log.Printf("DeleteById: no rows deleted for id %s", id)
		return errors.New("delete failed: no rows affected")
	}

	return nil
}
