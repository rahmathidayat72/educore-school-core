package model

import (
	"context"
	"errors"
	"fmt"
	matapelajaran "go_rest_native_sekolah/features/mata_pelajaran"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type mataPelajaranQuery struct {
	db *pgxpool.Pool
}

func NewDataMataPelajaran(db *pgxpool.Pool) matapelajaran.DataMataPelajaranInterface {
	if db == nil {
		panic("Nil database")
	}
	return &mataPelajaranQuery{db: db}
}

func toNullString(val string) interface{} {
	trimmed := strings.TrimSpace(val)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

func (m *mataPelajaranQuery) InsertMapel(insert *matapelajaran.MataPelajaranCore) error {
	if m.db == nil {
		return errors.New("Nil database")
	}
	if insert == nil {
		return errors.New("insert data is nil")
	}

	if insert.ID == "" {
		insert.ID = uuid.New().String()
	}

	// Validasi dan sinkronisasi Nama_Guru & ID_Guru.
	switch {
	case insert.Guru != "" && insert.ID_Guru == "":
		insert.Guru = strings.TrimSpace(insert.Guru)
		var guruID string
		err := m.db.QueryRow(context.Background(),
			"SELECT id FROM guru WHERE TRIM(nama) ILIKE TRIM($1)", insert.Guru).Scan(&guruID)
		if err != nil {
			log.Printf("InsertMapel: nama guru '%s' tidak ditemukan", insert.Guru)
			return fmt.Errorf("guru dengan nama '%s' tidak ditemukan", insert.Guru)
		}
		insert.ID_Guru = guruID

	case insert.ID_Guru != "" && insert.Guru == "":
		var namaGuru string
		err := m.db.QueryRow(context.Background(),
			"SELECT nama FROM guru WHERE id = $1", insert.ID_Guru).Scan(&namaGuru)
		if err != nil {
			log.Printf("InsertMapel: ID guru '%s' tidak ditemukan", insert.ID_Guru)
			return fmt.Errorf("guru dengan ID '%s' tidak ditemukan", insert.ID_Guru)
		}
		insert.Guru = namaGuru

	case insert.ID_Guru != "" && insert.Guru != "":
		insert.Guru = strings.TrimSpace(insert.Guru)
		var existingName string
		err := m.db.QueryRow(context.Background(),
			"SELECT nama FROM guru WHERE id = $1", insert.ID_Guru).Scan(&existingName)
		if err != nil {
			log.Printf("InsertMapel: ID guru '%s' tidak ditemukan", insert.ID_Guru)
			return fmt.Errorf("guru dengan ID '%s' tidak ditemukan", insert.ID_Guru)
		}
		if strings.ToLower(strings.TrimSpace(existingName)) != strings.ToLower(insert.Guru) {
			log.Printf("InsertMapel: Nama guru tidak cocok. Dapat: '%s', seharusnya: '%s'", insert.Guru, existingName)
			return fmt.Errorf("nama guru '%s' tidak cocok dengan ID guru '%s'", insert.Guru, insert.ID_Guru)
		}
	}

	// Validasi dan sinkronisasi Nama_Kelas & Kelas_ID.
	switch {
	case insert.Nama_Kelas != "" && insert.Kelas_ID == "":
		insert.Nama_Kelas = strings.TrimSpace(insert.Nama_Kelas)
		var kelasID string
		err := m.db.QueryRow(context.Background(),
			"SELECT id FROM kelas WHERE TRIM(kelas) ILIKE TRIM($1)", insert.Nama_Kelas).Scan(&kelasID)
		if err != nil {
			log.Printf("InsertMapel: nama kelas '%s' tidak ditemukan", insert.Nama_Kelas)
			return fmt.Errorf("kelas dengan nama '%s' tidak ditemukan", insert.Nama_Kelas)
		}
		insert.Kelas_ID = kelasID

	case insert.Kelas_ID != "" && insert.Nama_Kelas == "":
		var namaKelas string
		err := m.db.QueryRow(context.Background(),
			"SELECT kelas FROM kelas WHERE id = $1", insert.Kelas_ID).Scan(&namaKelas)
		if err != nil {
			log.Printf("InsertMapel: ID kelas '%s' tidak ditemukan", insert.Kelas_ID)
			return fmt.Errorf("kelas dengan ID '%s' tidak ditemukan", insert.Kelas_ID)
		}
		insert.Nama_Kelas = namaKelas

	case insert.Kelas_ID != "" && insert.Nama_Kelas != "":
		insert.Nama_Kelas = strings.TrimSpace(insert.Nama_Kelas)
		var existingKelas string
		err := m.db.QueryRow(context.Background(),
			"SELECT kelas FROM kelas WHERE id = $1", insert.Kelas_ID).Scan(&existingKelas)
		if err != nil {
			log.Printf("InsertMapel: ID kelas '%s' tidak ditemukan", insert.Kelas_ID)
			return fmt.Errorf("kelas dengan ID '%s' tidak ditemukan", insert.Kelas_ID)
		}
		if strings.ToLower(strings.TrimSpace(existingKelas)) != strings.ToLower(insert.Nama_Kelas) {
			log.Printf("InsertMapel: Nama kelas tidak cocok. Dapat: '%s', seharusnya: '%s'", insert.Nama_Kelas, existingKelas)
			return fmt.Errorf("nama kelas '%s' tidak cocok dengan ID kelas '%s'", insert.Nama_Kelas, insert.Kelas_ID)
		}
	}

	var idGuruParam interface{}
	if strings.TrimSpace(insert.ID_Guru) == "" {
		idGuruParam = nil
	} else {
		idGuruParam = strings.TrimSpace(insert.ID_Guru)
	}

	var idKelasParam interface{}
	if strings.TrimSpace(insert.Kelas_ID) == "" {
		idKelasParam = nil
	} else {
		idKelasParam = strings.TrimSpace(insert.Kelas_ID)
	}

	kurikulum := strings.TrimSpace(insert.Kurikulum)
	if kurikulum == "" {
		kurikulum = "Kurikulum Merdeka"
	}

	status := strings.TrimSpace(insert.Status)
	if status == "" {
		status = "Aktif"
	}

	jam := insert.JamPerMinggu
	if jam <= 0 {
		jam = 2
	}

	kkm := insert.KKM
	if kkm <= 0 {
		kkm = 75.0
	}

	query := `INSERT INTO mata_pelajaran (
		id, nama_pelajaran, kode_mapel, kelompok, kurikulum, tingkat, fase,
		jam_per_minggu, kkm, status, id_guru, kelas_id, deskripsi, update_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7,
		$8, $9, $10, $11, $12, $13, CURRENT_TIMESTAMP
	)`

	_, err := m.db.Exec(context.Background(), query,
		insert.ID,
		insert.Nama_Pelajaran,
		toNullString(insert.KodeMapel),
		toNullString(insert.Kelompok),
		kurikulum,
		toNullString(insert.Tingkat),
		toNullString(insert.Fase),
		jam,
		kkm,
		status,
		idGuruParam,
		idKelasParam,
		toNullString(insert.Deskripsi),
	)
	if err != nil {
		log.Printf("InsertMapel error exec: %v", err)
		return fmt.Errorf("insert failed: %w", err)
	}

	return nil
}

func (m *mataPelajaranQuery) SelectAllMapel() ([]matapelajaran.MataPelajaranCore, error) {
	if m.db == nil {
		return nil, errors.New("Nil database")
	}

	query := `SELECT 
		mp.id,
		COALESCE(mp.nama_pelajaran, ''),
		mp.kode_mapel,
		mp.kelompok,
		mp.kurikulum,
		mp.tingkat,
		mp.fase,
		mp.jam_per_minggu,
		mp.kkm,
		mp.status,
		mp.id_guru,
		g.nama AS nama_guru,
		mp.kelas_id,
		k.kelas AS nama_kelas,
		mp.deskripsi
	FROM mata_pelajaran mp
	LEFT JOIN guru g ON mp.id_guru = g.id
	LEFT JOIN kelas k ON mp.kelas_id = k.id
	WHERE mp.delete_at IS NULL
	ORDER BY mp.nama_pelajaran ASC;`

	rows, err := m.db.Query(context.Background(), query)
	if err != nil {
		log.Printf("SelectAllMapel error exec: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}
	defer rows.Close()

	var result []matapelajaran.MataPelajaranCore

	for rows.Next() {
		var mp MataPelajaran

		err = rows.Scan(
			&mp.ID,
			&mp.Nama_Pelajaran,
			&mp.KodeMapel,
			&mp.Kelompok,
			&mp.Kurikulum,
			&mp.Tingkat,
			&mp.Fase,
			&mp.JamPerMinggu,
			&mp.KKM,
			&mp.Status,
			&mp.ID_Guru,
			&mp.Guru,
			&mp.Kelas_ID,
			&mp.Nama_Kelas,
			&mp.Deskripsi,
		)
		if err != nil {
			log.Printf("SelectAllMapel error scan: %v", err)
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		core := FormatterResponse(mp)
		result = append(result, core)
	}

	if err := rows.Err(); err != nil {
		log.Printf("SelectAll error rows: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}

	log.Printf("Successfully fetched %d mata pelajaran from database", len(result))
	return result, nil
}

func (m *mataPelajaranQuery) SelectMapelById(id string) (*matapelajaran.MataPelajaranCore, error) {
	if m.db == nil {
		return nil, errors.New("Nil database")
	}
	if id == "" {
		return nil, errors.New("ID cannot be empty")
	}

	query := `SELECT 
		mp.id,
		COALESCE(mp.nama_pelajaran, ''),
		mp.kode_mapel,
		mp.kelompok,
		mp.kurikulum,
		mp.tingkat,
		mp.fase,
		mp.jam_per_minggu,
		mp.kkm,
		mp.status,
		mp.id_guru,
		g.nama AS nama_guru,
		mp.kelas_id,
		k.kelas AS nama_kelas,
		mp.deskripsi
	FROM mata_pelajaran mp
	LEFT JOIN guru g ON mp.id_guru = g.id
	LEFT JOIN kelas k ON mp.kelas_id = k.id
	WHERE mp.id = $1 AND mp.delete_at IS NULL`

	var mp MataPelajaran
	err := m.db.QueryRow(context.Background(), query, id).Scan(
		&mp.ID,
		&mp.Nama_Pelajaran,
		&mp.KodeMapel,
		&mp.Kelompok,
		&mp.Kurikulum,
		&mp.Tingkat,
		&mp.Fase,
		&mp.JamPerMinggu,
		&mp.KKM,
		&mp.Status,
		&mp.ID_Guru,
		&mp.Guru,
		&mp.Kelas_ID,
		&mp.Nama_Kelas,
		&mp.Deskripsi,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("ID not found")
		}
		log.Printf("QueryRow error: %v", err)
		return nil, fmt.Errorf("select failed: %w", err)
	}

	core := FormatterResponse(mp)
	return &core, nil
}

func (m *mataPelajaranQuery) UpdateMapel(update *matapelajaran.MataPelajaranCore, id string) error {
	if m.db == nil {
		return errors.New("Nil database")
	}
	if id == "" {
		return errors.New("ID cannot be empty")
	}

	var idGuruParam interface{}
	if strings.TrimSpace(update.ID_Guru) == "" {
		idGuruParam = nil
	} else {
		idGuruParam = strings.TrimSpace(update.ID_Guru)
	}

	var idKelasParam interface{}
	if strings.TrimSpace(update.Kelas_ID) == "" {
		idKelasParam = nil
	} else {
		idKelasParam = strings.TrimSpace(update.Kelas_ID)
	}

	kurikulum := strings.TrimSpace(update.Kurikulum)
	if kurikulum == "" {
		kurikulum = "Kurikulum Merdeka"
	}

	status := strings.TrimSpace(update.Status)
	if status == "" {
		status = "Aktif"
	}

	jam := update.JamPerMinggu
	if jam <= 0 {
		jam = 2
	}

	kkm := update.KKM
	if kkm <= 0 {
		kkm = 75.0
	}

	query := `
	UPDATE mata_pelajaran 
	SET nama_pelajaran = $1,
		kode_mapel = $2,
		kelompok = $3,
		kurikulum = $4,
		tingkat = $5,
		fase = $6,
		jam_per_minggu = $7,
		kkm = $8,
		status = $9,
		id_guru = $10,
		kelas_id = $11,
		deskripsi = $12,
		update_at = CURRENT_TIMESTAMP
	WHERE id = $13 AND delete_at IS NULL;
	`

	res, err := m.db.Exec(
		context.Background(), query,
		update.Nama_Pelajaran,
		toNullString(update.KodeMapel),
		toNullString(update.Kelompok),
		kurikulum,
		toNullString(update.Tingkat),
		toNullString(update.Fase),
		jam,
		kkm,
		status,
		idGuruParam,
		idKelasParam,
		toNullString(update.Deskripsi),
		id,
	)
	if err != nil {
		log.Printf("UpdateMapel error exec: %v", err)
		return fmt.Errorf("update failed: %w", err)
	}
	if res.RowsAffected() == 0 {
		log.Printf("UpdateMapel: no rows updated for id %s", id)
		return errors.New("update failed: no rows affected")
	}

	log.Printf("Successfully updated mata_pelajaran with id: %s", id)
	return nil
}

func (m *mataPelajaranQuery) DeleteMapel(id string) error {
	if m.db == nil {
		return errors.New("Nil database")
	}
	if id == "" {
		return errors.New("ID tidak boleh kosong")
	}

	query := "UPDATE mata_pelajaran SET delete_at = NOW() WHERE id = $1 AND delete_at IS NULL"
	res, err := m.db.Exec(context.Background(), query, id)
	if err != nil {
		log.Printf("DeleteMapel error exec: %v", err)
		return fmt.Errorf("delete failed: %w", err)
	}
	if res.RowsAffected() == 0 {
		log.Printf("DeleteMapel: no rows deleted for id %s", id)
		return errors.New("delete failed: no rows affected")
	}

	log.Printf("Successfully deleted mata_pelajaran with id: %s", id)
	return nil
}
