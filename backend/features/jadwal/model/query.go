package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/jadwal"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JadwalModel struct {
	db *pgxpool.Pool
}

func NewJadwalModel(db *pgxpool.Pool) jadwal.DataJadwalInterface {
	if db == nil {
		panic("jadwal model: nil database")
	}
	return &JadwalModel{db: db}
}

func toNullString(val string) interface{} {
	trimmed := strings.TrimSpace(val)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

func (m *JadwalModel) GetGuruIDByUserID(userID string) (string, error) {
	if userID == "" {
		return "", errors.New("userID kosong")
	}
	var guruID string
	err := m.db.QueryRow(context.Background(),
		"SELECT id FROM guru WHERE id_user = $1 AND delete_at IS NULL LIMIT 1", userID).Scan(&guruID)
	if err != nil {
		return "", err
	}
	return guruID, nil
}

func (m *JadwalModel) SelectAll(role, userID, hariFilter, kelasIDFilter, guruIDFilter, tahunFilter, semesterFilter, statusFilter string) ([]jadwal.JadwalCore, error) {
	ctx := context.Background()
	queryBuilder := `
		SELECT 
			j.id,
			COALESCE(j.kode_jadwal, '') AS kode_jadwal,
			j.mata_pelajaran_id,
			COALESCE(mp.kode_mapel, '') AS kode_mapel,
			COALESCE(mp.nama_pelajaran, '') AS nama_pelajaran,
			COALESCE(j.guru_id, mp.id_guru, '') AS id_guru,
			COALESCE(g.nama, '') AS nama_guru,
			COALESCE(g.nip, '') AS nip_guru,
			COALESCE(j.kelas_id, mp.kelas_id, '') AS kelas_id,
			COALESCE(k.kelas, '') AS nama_kelas,
			COALESCE(k.kode_kelas, '') AS kode_kelas,
			COALESCE(k.tingkat, '') AS tingkat_kelas,
			COALESCE(k.jurusan, '') AS jurusan_kelas,
			j.hari,
			COALESCE(j.jam_ke_mulai, 0) AS jam_ke_mulai,
			COALESCE(j.jam_ke_selesai, 0) AS jam_ke_selesai,
			to_char(j.jam_mulai, 'HH24:MI') AS jam_mulai,
			to_char(j.jam_selesai, 'HH24:MI') AS jam_selesai,
			COALESCE(j.ruang, '') AS ruang,
			COALESCE(j.jenis_kegiatan, 'Tatap Muka') AS jenis_kegiatan,
			COALESCE(j.tahun_ajaran, '') AS tahun_ajaran,
			COALESCE(j.semester, '') AS semester,
			COALESCE(j.keterangan, '') AS keterangan,
			j.aktif
		FROM jadwal_pelajaran j
		JOIN mata_pelajaran mp ON mp.id = j.mata_pelajaran_id AND mp.delete_at IS NULL
		LEFT JOIN guru g ON g.id = COALESCE(j.guru_id, mp.id_guru) AND g.delete_at IS NULL
		LEFT JOIN kelas k ON k.id = COALESCE(j.kelas_id, mp.kelas_id) AND k.delete_at IS NULL
		WHERE j.delete_at IS NULL
	`

	var args []interface{}
	argIdx := 1

	// Jika role guru, otomatis filter ke jadwal guru tersebut
	if strings.ToLower(role) == "guru" {
		guruID, err := m.GetGuruIDByUserID(userID)
		if err != nil {
			return []jadwal.JadwalCore{}, nil
		}
		queryBuilder += fmt.Sprintf(" AND COALESCE(j.guru_id, mp.id_guru) = $%d", argIdx)
		args = append(args, guruID)
		argIdx++
	} else if strings.TrimSpace(guruIDFilter) != "" {
		queryBuilder += fmt.Sprintf(" AND COALESCE(j.guru_id, mp.id_guru) = $%d", argIdx)
		args = append(args, strings.TrimSpace(guruIDFilter))
		argIdx++
	}

	if strings.TrimSpace(hariFilter) != "" {
		if h, err := strconv.Atoi(strings.TrimSpace(hariFilter)); err == nil && h >= 1 && h <= 7 {
			queryBuilder += fmt.Sprintf(" AND j.hari = $%d", argIdx)
			args = append(args, h)
			argIdx++
		}
	}

	if strings.TrimSpace(kelasIDFilter) != "" {
		queryBuilder += fmt.Sprintf(" AND COALESCE(j.kelas_id, mp.kelas_id) = $%d", argIdx)
		args = append(args, strings.TrimSpace(kelasIDFilter))
		argIdx++
	}

	if strings.TrimSpace(tahunFilter) != "" {
		queryBuilder += fmt.Sprintf(" AND LOWER(TRIM(j.tahun_ajaran)) = LOWER(TRIM($%d))", argIdx)
		args = append(args, strings.TrimSpace(tahunFilter))
		argIdx++
	}

	if strings.TrimSpace(semesterFilter) != "" {
		queryBuilder += fmt.Sprintf(" AND LOWER(TRIM(j.semester)) = LOWER(TRIM($%d))", argIdx)
		args = append(args, strings.TrimSpace(semesterFilter))
		argIdx++
	}

	if strings.TrimSpace(statusFilter) != "" {
		isAktif := strings.ToLower(strings.TrimSpace(statusFilter)) == "aktif" || strings.TrimSpace(statusFilter) == "true" || strings.TrimSpace(statusFilter) == "1"
		queryBuilder += fmt.Sprintf(" AND j.aktif = $%d", argIdx)
		args = append(args, isAktif)
		argIdx++
	}

	queryBuilder += " ORDER BY j.hari ASC, j.jam_mulai ASC, k.kelas ASC, j.kode_jadwal ASC"

	rows, err := m.db.Query(ctx, queryBuilder, args...)
	if err != nil {
		return nil, fmt.Errorf("gagal query jadwal: %w", err)
	}
	defer rows.Close()

	list := make([]jadwal.JadwalCore, 0)
	for rows.Next() {
		var item jadwal.JadwalCore
		err := rows.Scan(
			&item.ID,
			&item.KodeJadwal,
			&item.MataPelajaranID,
			&item.KodeMapel,
			&item.NamaPelajaran,
			&item.IDGuru,
			&item.NamaGuru,
			&item.NIPGuru,
			&item.KelasID,
			&item.NamaKelas,
			&item.KodeKelas,
			&item.TingkatKelas,
			&item.JurusanKelas,
			&item.Hari,
			&item.JamKeMulai,
			&item.JamKeSelesai,
			&item.JamMulai,
			&item.JamSelesai,
			&item.Ruang,
			&item.JenisKegiatan,
			&item.TahunAjaran,
			&item.Semester,
			&item.Keterangan,
			&item.Aktif,
		)
		if err != nil {
			return nil, fmt.Errorf("scan jadwal error: %w", err)
		}
		item.NamaHari = jadwal.NamaHariIndo(item.Hari)
		list = append(list, item)
	}

	return list, nil
}

func (m *JadwalModel) SelectByID(id string) (*jadwal.JadwalCore, error) {
	ctx := context.Background()
	query := `
		SELECT 
			j.id,
			COALESCE(j.kode_jadwal, '') AS kode_jadwal,
			j.mata_pelajaran_id,
			COALESCE(mp.kode_mapel, '') AS kode_mapel,
			COALESCE(mp.nama_pelajaran, '') AS nama_pelajaran,
			COALESCE(j.guru_id, mp.id_guru, '') AS id_guru,
			COALESCE(g.nama, '') AS nama_guru,
			COALESCE(g.nip, '') AS nip_guru,
			COALESCE(j.kelas_id, mp.kelas_id, '') AS kelas_id,
			COALESCE(k.kelas, '') AS nama_kelas,
			COALESCE(k.kode_kelas, '') AS kode_kelas,
			COALESCE(k.tingkat, '') AS tingkat_kelas,
			COALESCE(k.jurusan, '') AS jurusan_kelas,
			j.hari,
			COALESCE(j.jam_ke_mulai, 0) AS jam_ke_mulai,
			COALESCE(j.jam_ke_selesai, 0) AS jam_ke_selesai,
			to_char(j.jam_mulai, 'HH24:MI') AS jam_mulai,
			to_char(j.jam_selesai, 'HH24:MI') AS jam_selesai,
			COALESCE(j.ruang, '') AS ruang,
			COALESCE(j.jenis_kegiatan, 'Tatap Muka') AS jenis_kegiatan,
			COALESCE(j.tahun_ajaran, '') AS tahun_ajaran,
			COALESCE(j.semester, '') AS semester,
			COALESCE(j.keterangan, '') AS keterangan,
			j.aktif
		FROM jadwal_pelajaran j
		JOIN mata_pelajaran mp ON mp.id = j.mata_pelajaran_id AND mp.delete_at IS NULL
		LEFT JOIN guru g ON g.id = COALESCE(j.guru_id, mp.id_guru) AND g.delete_at IS NULL
		LEFT JOIN kelas k ON k.id = COALESCE(j.kelas_id, mp.kelas_id) AND k.delete_at IS NULL
		WHERE j.id = $1 AND j.delete_at IS NULL
	`

	var item jadwal.JadwalCore
	err := m.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.KodeJadwal,
		&item.MataPelajaranID,
		&item.KodeMapel,
		&item.NamaPelajaran,
		&item.IDGuru,
		&item.NamaGuru,
		&item.NIPGuru,
		&item.KelasID,
		&item.NamaKelas,
		&item.KodeKelas,
		&item.TingkatKelas,
		&item.JurusanKelas,
		&item.Hari,
		&item.JamKeMulai,
		&item.JamKeSelesai,
		&item.JamMulai,
		&item.JamSelesai,
		&item.Ruang,
		&item.JenisKegiatan,
		&item.TahunAjaran,
		&item.Semester,
		&item.Keterangan,
		&item.Aktif,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("jadwal tidak ditemukan")
		}
		return nil, err
	}
	item.NamaHari = jadwal.NamaHariIndo(item.Hari)
	return &item, nil
}

func (m *JadwalModel) CheckConflict(excludeID, mapelID, guruID, kelasID, hari, jamMulai, jamSelesai, ruang, tahunAjaran, semester string, aktif bool) (*jadwal.ConflictInfo, error) {
	// Jika status jadwal tidak aktif, tidak perlu memicu bentrok operasional
	if !aktif {
		return &jadwal.ConflictInfo{HasConflict: false}, nil
	}

	ctx := context.Background()

	// 1. Ambil detail mapel default
	var defaultGuruID, defaultKelasID, defaultNamaMapel string
	err := m.db.QueryRow(ctx, `
		SELECT 
			COALESCE(mp.id_guru, ''),
			COALESCE(mp.kelas_id, ''),
			COALESCE(mp.nama_pelajaran, '')
		FROM mata_pelajaran mp
		WHERE mp.id = $1 AND mp.delete_at IS NULL
	`, mapelID).Scan(&defaultGuruID, &defaultKelasID, &defaultNamaMapel)
	if err != nil {
		return nil, fmt.Errorf("mata pelajaran tidak ditemukan: %w", err)
	}

	actualGuruID := strings.TrimSpace(guruID)
	if actualGuruID == "" {
		actualGuruID = defaultGuruID
	}

	actualKelasID := strings.TrimSpace(kelasID)
	if actualKelasID == "" {
		actualKelasID = defaultKelasID
	}

	h, _ := strconv.Atoi(hari)
	tahun := strings.TrimSpace(tahunAjaran)
	if tahun == "" {
		tahun = "2026/2027"
	}
	sem := strings.TrimSpace(semester)
	if sem == "" {
		sem = "Ganjil"
	}

	// 2. Cek konflik Guru: Apakah guru sudah mengajar pada jam dan semester/tahun yang sama?
	if actualGuruID != "" {
		guruCheckQuery := `
			SELECT 
				COALESCE(g2.nama, 'Guru'),
				COALESCE(k2.kelas, 'Lain'), 
				COALESCE(mp2.nama_pelajaran, ''), 
				to_char(j2.jam_mulai, 'HH24:MI'), 
				to_char(j2.jam_selesai, 'HH24:MI')
			FROM jadwal_pelajaran j2
			JOIN mata_pelajaran mp2 ON mp2.id = j2.mata_pelajaran_id AND mp2.delete_at IS NULL
			LEFT JOIN guru g2 ON g2.id = COALESCE(j2.guru_id, mp2.id_guru) AND g2.delete_at IS NULL
			LEFT JOIN kelas k2 ON k2.id = COALESCE(j2.kelas_id, mp2.kelas_id) AND k2.delete_at IS NULL
			WHERE j2.delete_at IS NULL
			  AND j2.aktif = TRUE
			  AND COALESCE(j2.guru_id, mp2.id_guru) = $1
			  AND j2.hari = $2
			  AND ($3 = '' OR j2.id != $3)
			  AND LOWER(TRIM(j2.tahun_ajaran)) = LOWER(TRIM($4))
			  AND LOWER(TRIM(j2.semester)) = LOWER(TRIM($5))
			  AND (j2.jam_mulai < $7::time AND j2.jam_selesai > $6::time)
			LIMIT 1
		`
		var confGuru, confKelas, confMapel, confMulai, confSelesai string
		err := m.db.QueryRow(ctx, guruCheckQuery, actualGuruID, h, excludeID, tahun, sem, jamMulai, jamSelesai).
			Scan(&confGuru, &confKelas, &confMapel, &confMulai, &confSelesai)
		if err == nil {
			return &jadwal.ConflictInfo{
				HasConflict: true,
				Message: fmt.Sprintf("Guru %s sudah memiliki jadwal mengajar %s di kelas %s pukul %s-%s.",
					confGuru, confMapel, confKelas, confMulai, confSelesai),
			}, nil
		}
	}

	// 3. Cek konflik Kelas: Apakah kelas sudah ada pelajaran lain pada jam bertabrakan?
	if actualKelasID != "" {
		kelasCheckQuery := `
			SELECT 
				COALESCE(k2.kelas, 'Kelas'),
				COALESCE(mp2.nama_pelajaran, ''), 
				COALESCE(g2.nama, 'Guru'), 
				to_char(j2.jam_mulai, 'HH24:MI'), 
				to_char(j2.jam_selesai, 'HH24:MI')
			FROM jadwal_pelajaran j2
			JOIN mata_pelajaran mp2 ON mp2.id = j2.mata_pelajaran_id AND mp2.delete_at IS NULL
			LEFT JOIN guru g2 ON g2.id = COALESCE(j2.guru_id, mp2.id_guru) AND g2.delete_at IS NULL
			LEFT JOIN kelas k2 ON k2.id = COALESCE(j2.kelas_id, mp2.kelas_id) AND k2.delete_at IS NULL
			WHERE j2.delete_at IS NULL
			  AND j2.aktif = TRUE
			  AND COALESCE(j2.kelas_id, mp2.kelas_id) = $1
			  AND j2.hari = $2
			  AND ($3 = '' OR j2.id != $3)
			  AND LOWER(TRIM(j2.tahun_ajaran)) = LOWER(TRIM($4))
			  AND LOWER(TRIM(j2.semester)) = LOWER(TRIM($5))
			  AND (j2.jam_mulai < $7::time AND j2.jam_selesai > $6::time)
			LIMIT 1
		`
		var confKelas, confMapel, confGuru, confMulai, confSelesai string
		err := m.db.QueryRow(ctx, kelasCheckQuery, actualKelasID, h, excludeID, tahun, sem, jamMulai, jamSelesai).
			Scan(&confKelas, &confMapel, &confGuru, &confMulai, &confSelesai)
		if err == nil {
			return &jadwal.ConflictInfo{
				HasConflict: true,
				Message: fmt.Sprintf("Kelas %s sudah memiliki pelajaran %s (%s) pukul %s-%s.",
					confKelas, confMapel, confGuru, confMulai, confSelesai),
			}, nil
		}
	}

	// 4. Cek konflik Ruangan jika ruang diisi
	ruangTrim := strings.TrimSpace(ruang)
	if ruangTrim != "" {
		ruangCheckQuery := `
			SELECT 
				COALESCE(k2.kelas, 'Lain'), 
				COALESCE(mp2.nama_pelajaran, ''), 
				to_char(j2.jam_mulai, 'HH24:MI'), 
				to_char(j2.jam_selesai, 'HH24:MI')
			FROM jadwal_pelajaran j2
			JOIN mata_pelajaran mp2 ON mp2.id = j2.mata_pelajaran_id AND mp2.delete_at IS NULL
			LEFT JOIN kelas k2 ON k2.id = COALESCE(j2.kelas_id, mp2.kelas_id) AND k2.delete_at IS NULL
			WHERE j2.delete_at IS NULL
			  AND j2.aktif = TRUE
			  AND LOWER(TRIM(j2.ruang)) = LOWER($1)
			  AND j2.hari = $2
			  AND ($3 = '' OR j2.id != $3)
			  AND LOWER(TRIM(j2.tahun_ajaran)) = LOWER(TRIM($4))
			  AND LOWER(TRIM(j2.semester)) = LOWER(TRIM($5))
			  AND (j2.jam_mulai < $7::time AND j2.jam_selesai > $6::time)
			LIMIT 1
		`
		var confKelas, confMapel, confMulai, confSelesai string
		err := m.db.QueryRow(ctx, ruangCheckQuery, ruangTrim, h, excludeID, tahun, sem, jamMulai, jamSelesai).
			Scan(&confKelas, &confMapel, &confMulai, &confSelesai)
		if err == nil {
			return &jadwal.ConflictInfo{
				HasConflict: true,
				Message: fmt.Sprintf("Ruangan %s sudah digunakan kelas %s untuk pelajaran %s pukul %s-%s.",
					ruangTrim, confKelas, confMapel, confMulai, confSelesai),
			}, nil
		}
	}

	return &jadwal.ConflictInfo{HasConflict: false}, nil
}

func (m *JadwalModel) Insert(input *jadwal.JadwalCore) error {
	ctx := context.Background()
	if input.ID == "" {
		input.ID = uuid.New().String()
	}

	tahun := strings.TrimSpace(input.TahunAjaran)
	if tahun == "" {
		tahun = "2026/2027"
	}
	sem := strings.TrimSpace(input.Semester)
	if sem == "" {
		sem = "Ganjil"
	}

	jenis := strings.TrimSpace(input.JenisKegiatan)
	if jenis == "" {
		jenis = "Tatap Muka"
	}

	var jamKeMulai, jamKeSelesai sql.NullInt32
	if input.JamKeMulai > 0 {
		jamKeMulai = sql.NullInt32{Int32: int32(input.JamKeMulai), Valid: true}
	}
	if input.JamKeSelesai > 0 {
		jamKeSelesai = sql.NullInt32{Int32: int32(input.JamKeSelesai), Valid: true}
	}

	query := `
		INSERT INTO jadwal_pelajaran (
			id, kode_jadwal, mata_pelajaran_id, guru_id, kelas_id,
			hari, jam_ke_mulai, jam_ke_selesai, jam_mulai, jam_selesai,
			ruang, jenis_kegiatan, tahun_ajaran, semester, keterangan,
			aktif, created_at, update_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9::time, $10::time,
			$11, $12, $13, $14, $15,
			$16, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
	`
	_, err := m.db.Exec(ctx, query,
		input.ID,
		toNullString(input.KodeJadwal),
		input.MataPelajaranID,
		toNullString(input.IDGuru),
		toNullString(input.KelasID),
		input.Hari,
		jamKeMulai,
		jamKeSelesai,
		input.JamMulai,
		input.JamSelesai,
		toNullString(input.Ruang),
		jenis,
		tahun,
		sem,
		toNullString(input.Keterangan),
		input.Aktif,
	)
	if err != nil {
		if strings.Contains(err.Error(), "idx_jadwal_kode_unique") {
			return errors.New("kode jadwal sudah digunakan pada sesi jadwal lain")
		}
		return err
	}
	return nil
}

func (m *JadwalModel) Update(id string, input *jadwal.JadwalCore) error {
	ctx := context.Background()

	tahun := strings.TrimSpace(input.TahunAjaran)
	if tahun == "" {
		tahun = "2026/2027"
	}
	sem := strings.TrimSpace(input.Semester)
	if sem == "" {
		sem = "Ganjil"
	}

	jenis := strings.TrimSpace(input.JenisKegiatan)
	if jenis == "" {
		jenis = "Tatap Muka"
	}

	var jamKeMulai, jamKeSelesai sql.NullInt32
	if input.JamKeMulai > 0 {
		jamKeMulai = sql.NullInt32{Int32: int32(input.JamKeMulai), Valid: true}
	}
	if input.JamKeSelesai > 0 {
		jamKeSelesai = sql.NullInt32{Int32: int32(input.JamKeSelesai), Valid: true}
	}

	query := `
		UPDATE jadwal_pelajaran
		SET kode_jadwal = $1,
		    mata_pelajaran_id = $2,
		    guru_id = $3,
		    kelas_id = $4,
		    hari = $5,
		    jam_ke_mulai = $6,
		    jam_ke_selesai = $7,
		    jam_mulai = $8::time,
		    jam_selesai = $9::time,
		    ruang = $10,
		    jenis_kegiatan = $11,
		    tahun_ajaran = $12,
		    semester = $13,
		    keterangan = $14,
		    aktif = $15,
		    update_at = CURRENT_TIMESTAMP
		WHERE id = $16 AND delete_at IS NULL
	`
	res, err := m.db.Exec(ctx, query,
		toNullString(input.KodeJadwal),
		input.MataPelajaranID,
		toNullString(input.IDGuru),
		toNullString(input.KelasID),
		input.Hari,
		jamKeMulai,
		jamKeSelesai,
		input.JamMulai,
		input.JamSelesai,
		toNullString(input.Ruang),
		jenis,
		tahun,
		sem,
		toNullString(input.Keterangan),
		input.Aktif,
		id,
	)
	if err != nil {
		if strings.Contains(err.Error(), "idx_jadwal_kode_unique") {
			return errors.New("kode jadwal sudah digunakan pada sesi jadwal lain")
		}
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("tidak ada data jadwal yang diperbarui")
	}
	return nil
}

func (m *JadwalModel) Delete(id string) error {
	ctx := context.Background()
	query := `
		UPDATE jadwal_pelajaran
		SET delete_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND delete_at IS NULL
	`
	res, err := m.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("jadwal tidak ditemukan atau sudah dihapus")
	}
	return nil
}

var _ = time.Now
