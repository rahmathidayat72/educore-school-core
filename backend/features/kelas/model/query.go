package model

import (
	"context"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/kelas"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type kelasQuery struct {
	db *pgxpool.Pool
}

func NewDataKelas(db *pgxpool.Pool) kelas.DataKelasInterface {
	if db == nil {
		panic("Nil database")
	}
	return &kelasQuery{db: db}
}

func toNullString(val string) interface{} {
	trimmed := strings.TrimSpace(val)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

func (k *kelasQuery) Insert(insert *kelas.KelasCore) error {
	if k.db == nil {
		return errors.New("kelas model: Nil database connection")
	}
	if insert == nil {
		return errors.New("insert data is nil")
	}

	if insert.ID == "" {
		insert.ID = uuid.New().String()
	}

	// Sinkronisasi Nama_Guru & ID_Guru
	switch {
	case insert.Nama_Guru != "" && insert.ID_Guru == "":
		var idGuru string
		err := k.db.QueryRow(context.Background(), "SELECT id FROM guru WHERE TRIM(nama) ILIKE TRIM($1)", insert.Nama_Guru).Scan(&idGuru)
		if err != nil {
			log.Printf("InsertKelas: nama guru '%s' tidak ditemukan", insert.Nama_Guru)
			return fmt.Errorf("guru dengan nama '%s' tidak ditemukan", insert.Nama_Guru)
		}
		insert.ID_Guru = idGuru

	case insert.ID_Guru != "" && insert.Nama_Guru == "":
		var namaGuru string
		err := k.db.QueryRow(context.Background(), "SELECT nama FROM guru WHERE id = $1", insert.ID_Guru).Scan(&namaGuru)
		if err != nil {
			log.Printf("InsertKelas: ID guru '%s' tidak ditemukan", insert.ID_Guru)
			return fmt.Errorf("guru dengan ID '%s' tidak ditemukan", insert.ID_Guru)
		}
		insert.Nama_Guru = namaGuru

	case insert.ID_Guru != "" && insert.Nama_Guru != "":
		var existingName string
		err := k.db.QueryRow(context.Background(), "SELECT nama FROM guru WHERE id = $1", insert.ID_Guru).Scan(&existingName)
		if err != nil {
			log.Printf("InsertKelas: ID guru '%s' tidak ditemukan", insert.ID_Guru)
			return fmt.Errorf("guru dengan ID '%s' tidak ditemukan", insert.ID_Guru)
		}
		if strings.TrimSpace(strings.ToLower(existingName)) != strings.TrimSpace(strings.ToLower(insert.Nama_Guru)) {
			log.Printf("InsertKelas: Nama guru tidak cocok dengan ID guru. Dapat: '%s', seharusnya: '%s'", insert.Nama_Guru, existingName)
			return fmt.Errorf("nama guru '%s' tidak cocok dengan ID guru '%s'", insert.Nama_Guru, insert.ID_Guru)
		}
	}

	var idGuruParam interface{}
	if strings.TrimSpace(insert.ID_Guru) == "" {
		idGuruParam = nil
	} else {
		idGuruParam = strings.TrimSpace(insert.ID_Guru)
	}

	kapasitas := insert.Kapasitas
	if kapasitas <= 0 {
		kapasitas = 36
	}

	status := strings.TrimSpace(insert.Status)
	if status == "" {
		status = "Aktif"
	}

	tahunAjaran := strings.TrimSpace(insert.TahunAjaran)
	if tahunAjaran == "" {
		tahunAjaran = "2026/2027"
	}

	query := `INSERT INTO kelas (
		id, kelas, kode_kelas, tingkat, fase, jurusan, tahun_ajaran, ruangan,
		kapasitas, status, deskripsi, id_guru, school_unit_id, update_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8,
		$9, $10, $11, $12, $13, CURRENT_TIMESTAMP
	)`

	var unitIDParam interface{}
	if insert.SchoolUnitID != nil && strings.TrimSpace(*insert.SchoolUnitID) != "" {
		unitIDParam = strings.TrimSpace(*insert.SchoolUnitID)
	}

	_, err := k.db.Exec(context.Background(), query,
		insert.ID,
		insert.Kelas,
		toNullString(insert.KodeKelas),
		toNullString(insert.Tingkat),
		toNullString(insert.Fase),
		toNullString(insert.Jurusan),
		tahunAjaran,
		toNullString(insert.Ruangan),
		kapasitas,
		status,
		toNullString(insert.Deskripsi),
		idGuruParam,
		unitIDParam,
	)
	if err != nil {
		log.Printf("InsertKelas error exec: %v", err)
		if strings.Contains(err.Error(), "idx_kelas_kode_unique") {
			return errors.New("kode kelas sudah terdaftar untuk rombel lain")
		}
		return fmt.Errorf("insert failed: %w", err)
	}

	return nil
}

func (k *kelasQuery) SelectAll(schoolUnitID ...string) ([]kelas.KelasCore, error) {
	if k.db == nil {
		return nil, errors.New("Nil database connection")
	}

	whereClause := "WHERE k.delete_at IS NULL"
	var args []interface{}
	if len(schoolUnitID) > 0 && strings.TrimSpace(schoolUnitID[0]) != "" {
		whereClause += " AND k.school_unit_id = $1"
		args = append(args, strings.TrimSpace(schoolUnitID[0]))
	}

	query := fmt.Sprintf(`
		SELECT 
			k.id, 
			k.kelas, 
			k.kode_kelas,
			k.tingkat,
			k.fase,
			k.jurusan,
			k.tahun_ajaran,
			k.ruangan,
			k.kapasitas,
			k.status,
			k.deskripsi,
			k.id_guru, 
			g.nama AS nama_guru,
			COUNT(s.id) AS total_siswa,
			k.school_unit_id
		FROM 
			kelas k
		LEFT JOIN 
			guru g ON k.id_guru = g.id
		LEFT JOIN
			siswa s ON s.kelas_id = k.id AND s.delete_at IS NULL
		%s
		GROUP BY 
			k.id, g.nama
		ORDER BY 
			k.kelas ASC
	`, whereClause)

	rows, err := k.db.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("SelectAll error exec: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}
	defer rows.Close()

	var result []kelas.KelasCore

	for rows.Next() {
		var kl Kelas

		err := rows.Scan(
			&kl.ID,
			&kl.Kelas,
			&kl.KodeKelas,
			&kl.Tingkat,
			&kl.Fase,
			&kl.Jurusan,
			&kl.TahunAjaran,
			&kl.Ruangan,
			&kl.Kapasitas,
			&kl.Status,
			&kl.Deskripsi,
			&kl.ID_Guru,
			&kl.Nama_Guru,
			&kl.TotalSiswa,
			&kl.SchoolUnitID,
		)
		if err != nil {
			log.Printf("SelectAll error scan: %v", err)
			return nil, fmt.Errorf("select failed: %w", err)
		}

		core := FormatterResponse(kl)
		result = append(result, core)
	}

	if err := rows.Err(); err != nil {
		log.Printf("SelectAll error rows: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}

	log.Printf("Successfully fetched %d kelas from database", len(result))
	return result, nil
}

func (k *kelasQuery) SelectById(id string, schoolUnitID ...string) (*kelas.KelasCore, error) {
	if k.db == nil {
		return nil, errors.New("Nil database connection")
	}

	whereClause := "WHERE k.id = $1 AND k.delete_at IS NULL"
	args := []interface{}{id}
	if len(schoolUnitID) > 0 && strings.TrimSpace(schoolUnitID[0]) != "" {
		whereClause += " AND k.school_unit_id = $2"
		args = append(args, strings.TrimSpace(schoolUnitID[0]))
	}

	query := fmt.Sprintf(`
		SELECT 
			k.id, 
			k.kelas, 
			k.kode_kelas,
			k.tingkat,
			k.fase,
			k.jurusan,
			k.tahun_ajaran,
			k.ruangan,
			k.kapasitas,
			k.status,
			k.deskripsi,
			k.id_guru, 
			g.nama AS nama_guru,
			COUNT(s.id) AS total_siswa,
			k.school_unit_id
		FROM 
			kelas k
		LEFT JOIN 
			guru g ON k.id_guru = g.id 
		LEFT JOIN
			siswa s ON s.kelas_id = k.id AND s.delete_at IS NULL
		%s
		GROUP BY
			k.id, g.nama
	`, whereClause)

	var kl Kelas
	err := k.db.QueryRow(context.Background(), query, args...).Scan(
		&kl.ID,
		&kl.Kelas,
		&kl.KodeKelas,
		&kl.Tingkat,
		&kl.Fase,
		&kl.Jurusan,
		&kl.TahunAjaran,
		&kl.Ruangan,
		&kl.Kapasitas,
		&kl.Status,
		&kl.Deskripsi,
		&kl.ID_Guru,
		&kl.Nama_Guru,
		&kl.TotalSiswa,
		&kl.SchoolUnitID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		log.Printf("SelectById error exec: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}

	core := FormatterResponse(kl)
	return &core, nil
}

func (k *kelasQuery) Update(insert *kelas.KelasCore, id string, schoolUnitID ...string) error {
	if k.db == nil {
		return errors.New("Nil database connection")
	}

	var idGuruParam interface{}
	if strings.TrimSpace(insert.ID_Guru) == "" {
		idGuruParam = nil
	} else {
		idGuruParam = strings.TrimSpace(insert.ID_Guru)
	}

	kapasitas := insert.Kapasitas
	if kapasitas <= 0 {
		kapasitas = 36
	}

	status := strings.TrimSpace(insert.Status)
	if status == "" {
		status = "Aktif"
	}

	tahunAjaran := strings.TrimSpace(insert.TahunAjaran)
	if tahunAjaran == "" {
		tahunAjaran = "2026/2027"
	}

	query := `UPDATE kelas SET 
		kelas = $1, 
		kode_kelas = $2,
		tingkat = $3,
		fase = $4,
		jurusan = $5,
		tahun_ajaran = $6,
		ruangan = $7,
		kapasitas = $8,
		status = $9,
		deskripsi = $10,
		id_guru = $11,
		update_at = CURRENT_TIMESTAMP
	WHERE id = $12 AND delete_at IS NULL`

	args := []interface{}{
		insert.Kelas,
		toNullString(insert.KodeKelas),
		toNullString(insert.Tingkat),
		toNullString(insert.Fase),
		toNullString(insert.Jurusan),
		tahunAjaran,
		toNullString(insert.Ruangan),
		kapasitas,
		status,
		toNullString(insert.Deskripsi),
		idGuruParam,
		id,
	}

	if len(schoolUnitID) > 0 && strings.TrimSpace(schoolUnitID[0]) != "" {
		query += fmt.Sprintf(" AND school_unit_id = $%d", len(args)+1)
		args = append(args, strings.TrimSpace(schoolUnitID[0]))
	}

	res, err := k.db.Exec(context.Background(), query, args...)
	if err != nil {
		log.Printf("UpdateKelas error exec: %v", err)
		if strings.Contains(err.Error(), "idx_kelas_kode_unique") {
			return errors.New("kode kelas sudah terdaftar untuk rombel lain")
		}
		return fmt.Errorf("update failed: %w", err)
	}

	if res.RowsAffected() == 0 {
		log.Printf("Updatekelas: no rows updated for id %s", id)
		return errors.New("update failed: no rows affected")
	}

	return nil
}

func (k *kelasQuery) DeleteById(id string, schoolUnitID ...string) error {
	if k.db == nil {
		return errors.New("Nil database connection")
	}

	query := "UPDATE kelas SET delete_at = NOW() WHERE id = $1 AND delete_at IS NULL"
	args := []interface{}{id}
	if len(schoolUnitID) > 0 && strings.TrimSpace(schoolUnitID[0]) != "" {
		query += " AND school_unit_id = $2"
		args = append(args, strings.TrimSpace(schoolUnitID[0]))
	}

	res, err := k.db.Exec(context.Background(), query, args...)
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
