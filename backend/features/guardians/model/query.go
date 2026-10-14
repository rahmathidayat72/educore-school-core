package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/guardians"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GuardianModel struct {
	db *pgxpool.Pool
}

func NewGuardianModel(db *pgxpool.Pool) guardians.DataGuardianInterface {
	if db == nil {
		panic("nil db pool on guardian model")
	}
	return &GuardianModel{db: db}
}

func toNullStr(s string) sql.NullString {
	trim := strings.TrimSpace(s)
	if trim == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: trim, Valid: true}
}

func (m *GuardianModel) SelectAll(search string) ([]guardians.GuardianCore, error) {
	ctx := context.Background()
	queryBuilder := `
		SELECT 
			id, user_id, nama, nik, hubungan, pekerjaan, penghasilan,
			telepon, email, alamat, created_at, update_at, delete_at
		FROM guardians
		WHERE delete_at IS NULL
	`
	var args []interface{}
	if strings.TrimSpace(search) != "" {
		queryBuilder += ` AND (
			LOWER(nama) LIKE LOWER($1) OR 
			telepon LIKE $1 OR 
			nik LIKE $1
		)`
		args = append(args, "%"+strings.TrimSpace(search)+"%")
	}
	queryBuilder += " ORDER BY nama ASC"

	rows, err := m.db.Query(ctx, queryBuilder, args...)
	if err != nil {
		return nil, fmt.Errorf("gagal query daftar wali murid: %w", err)
	}
	defer rows.Close()

	var list []Guardian
	for rows.Next() {
		var item Guardian
		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Nama,
			&item.NIK,
			&item.Hubungan,
			&item.Pekerjaan,
			&item.Penghasilan,
			&item.Telepon,
			&item.Email,
			&item.Alamat,
			&item.CreatedAt,
			&item.UpdateAt,
			&item.DeleteAt,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan data wali: %w", err)
		}
		list = append(list, item)
	}

	return toGuardianCoreList(list), nil
}

func (m *GuardianModel) SelectByID(id string) (*guardians.GuardianCore, error) {
	ctx := context.Background()
	query := `
		SELECT 
			id, user_id, nama, nik, hubungan, pekerjaan, penghasilan,
			telepon, email, alamat, created_at, update_at, delete_at
		FROM guardians
		WHERE id = $1 AND delete_at IS NULL
	`
	var item Guardian
	err := m.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.UserID,
		&item.Nama,
		&item.NIK,
		&item.Hubungan,
		&item.Pekerjaan,
		&item.Penghasilan,
		&item.Telepon,
		&item.Email,
		&item.Alamat,
		&item.CreatedAt,
		&item.UpdateAt,
		&item.DeleteAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("data wali murid tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal query detail wali: %w", err)
	}

	core := item.toCore()
	return &core, nil
}

func (m *GuardianModel) SelectByUserID(userID string) (*guardians.GuardianCore, error) {
	ctx := context.Background()
	query := `
		SELECT 
			id, user_id, nama, nik, hubungan, pekerjaan, penghasilan,
			telepon, email, alamat, created_at, update_at, delete_at
		FROM guardians
		WHERE user_id = $1 AND delete_at IS NULL
		LIMIT 1
	`
	var item Guardian
	err := m.db.QueryRow(ctx, query, userID).Scan(
		&item.ID,
		&item.UserID,
		&item.Nama,
		&item.NIK,
		&item.Hubungan,
		&item.Pekerjaan,
		&item.Penghasilan,
		&item.Telepon,
		&item.Email,
		&item.Alamat,
		&item.CreatedAt,
		&item.UpdateAt,
		&item.DeleteAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("profil wali murid untuk user ini tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal query profil wali user: %w", err)
	}

	core := item.toCore()
	return &core, nil
}

func (m *GuardianModel) Insert(input *guardians.GuardianCore) error {
	ctx := context.Background()
	if strings.TrimSpace(input.ID) == "" {
		input.ID = uuid.New().String()
	}

	query := `
		INSERT INTO guardians (
			id, user_id, nama, nik, hubungan, pekerjaan, penghasilan,
			telepon, email, alamat, created_at, update_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	var uid sql.NullString
	if input.UserID != nil && strings.TrimSpace(*input.UserID) != "" {
		uid = sql.NullString{String: strings.TrimSpace(*input.UserID), Valid: true}
	}

	_, err := m.db.Exec(
		ctx,
		query,
		input.ID,
		uid,
		input.Nama,
		toNullStr(input.NIK),
		input.Hubungan,
		toNullStr(input.Pekerjaan),
		toNullStr(input.Penghasilan),
		input.Telepon,
		toNullStr(input.Email),
		toNullStr(input.Alamat),
	)
	if err != nil {
		return fmt.Errorf("gagal insert data wali: %w", err)
	}
	return nil
}

func (m *GuardianModel) Update(id string, input *guardians.GuardianCore) error {
	ctx := context.Background()
	query := `
		UPDATE guardians
		SET nama = $1, nik = $2, hubungan = $3, pekerjaan = $4, penghasilan = $5,
		    telepon = $6, email = $7, alamat = $8, user_id = $9, update_at = CURRENT_TIMESTAMP
		WHERE id = $10 AND delete_at IS NULL
	`
	var uid sql.NullString
	if input.UserID != nil && strings.TrimSpace(*input.UserID) != "" {
		uid = sql.NullString{String: strings.TrimSpace(*input.UserID), Valid: true}
	}

	cmd, err := m.db.Exec(
		ctx,
		query,
		input.Nama,
		toNullStr(input.NIK),
		input.Hubungan,
		toNullStr(input.Pekerjaan),
		toNullStr(input.Penghasilan),
		input.Telepon,
		toNullStr(input.Email),
		toNullStr(input.Alamat),
		uid,
		id,
	)
	if err != nil {
		return fmt.Errorf("gagal update data wali: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("data wali murid tidak ditemukan")
	}
	return nil
}

func (m *GuardianModel) Delete(id string) error {
	ctx := context.Background()
	query := `
		UPDATE guardians
		SET delete_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND delete_at IS NULL
	`
	cmd, err := m.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus wali murid: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("data wali murid tidak ditemukan")
	}
	return nil
}

// StudentGuardian Relations

func (m *GuardianModel) LinkStudent(input *guardians.StudentGuardianCore) error {
	ctx := context.Background()
	if strings.TrimSpace(input.ID) == "" {
		input.ID = uuid.New().String()
	}

	query := `
		INSERT INTO student_guardians (
			id, guardian_id, siswa_id, is_primary, created_at, update_at
		) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	_, err := m.db.Exec(ctx, query, input.ID, input.GuardianID, input.SiswaID, input.IsPrimary)
	if err != nil {
		if strings.Contains(err.Error(), "idx_student_guardians_unique") {
			return errors.New("siswa sudah ditautkan dengan wali murid ini")
		}
		return fmt.Errorf("gagal menautkan siswa dan wali: %w", err)
	}
	return nil
}

func (m *GuardianModel) UnlinkStudent(guardianID string, siswaID string) error {
	ctx := context.Background()
	query := `
		UPDATE student_guardians
		SET delete_at = CURRENT_TIMESTAMP
		WHERE guardian_id = $1 AND siswa_id = $2 AND delete_at IS NULL
	`
	cmd, err := m.db.Exec(ctx, query, guardianID, siswaID)
	if err != nil {
		return fmt.Errorf("gagal memutuskan relasi siswa dan wali: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("relasi siswa dan wali tidak ditemukan")
	}
	return nil
}

func (m *GuardianModel) SelectStudentsByGuardianID(guardianID string) ([]guardians.StudentGuardianCore, error) {
	ctx := context.Background()
	query := `
		SELECT 
			sg.id, sg.guardian_id, sg.siswa_id, sg.is_primary, sg.created_at, sg.update_at, sg.delete_at,
			s.nama as nama_siswa, s.nisn,
			k.nama_kelas,
			u.nama_sekolah
		FROM student_guardians sg
		JOIN siswa s ON sg.siswa_id = s.id
		LEFT JOIN kelas k ON s.kelas_id = k.id
		LEFT JOIN school_units u ON s.school_unit_id = u.id
		WHERE sg.guardian_id = $1 AND sg.delete_at IS NULL AND s.delete_at IS NULL
		ORDER BY s.nama ASC
	`
	rows, err := m.db.Query(ctx, query, guardianID)
	if err != nil {
		return nil, fmt.Errorf("gagal query daftar anak terhubung: %w", err)
	}
	defer rows.Close()

	var list []StudentGuardian
	for rows.Next() {
		var item StudentGuardian
		err := rows.Scan(
			&item.ID,
			&item.GuardianID,
			&item.SiswaID,
			&item.IsPrimary,
			&item.CreatedAt,
			&item.UpdateAt,
			&item.DeleteAt,
			&item.NamaSiswa,
			&item.NISN,
			&item.NamaKelas,
			&item.NamaSekolah,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan anak terhubung: %w", err)
		}
		list = append(list, item)
	}

	return toStudentGuardianCoreList(list), nil
}

func (m *GuardianModel) SelectGuardiansBySiswaID(siswaID string) ([]guardians.StudentGuardianCore, error) {
	ctx := context.Background()
	query := `
		SELECT 
			sg.id, sg.guardian_id, sg.siswa_id, sg.is_primary, sg.created_at, sg.update_at, sg.delete_at,
			s.nama as nama_siswa, s.nisn,
			k.nama_kelas,
			u.nama_sekolah,
			g.nama as nama_wali, g.hubungan, g.telepon as telepon_wali
		FROM student_guardians sg
		JOIN guardians g ON sg.guardian_id = g.id
		JOIN siswa s ON sg.siswa_id = s.id
		LEFT JOIN kelas k ON s.kelas_id = k.id
		LEFT JOIN school_units u ON s.school_unit_id = u.id
		WHERE sg.siswa_id = $1 AND sg.delete_at IS NULL AND g.delete_at IS NULL
		ORDER BY sg.is_primary DESC, g.nama ASC
	`
	rows, err := m.db.Query(ctx, query, siswaID)
	if err != nil {
		return nil, fmt.Errorf("gagal query wali siswa: %w", err)
	}
	defer rows.Close()

	var list []StudentGuardian
	for rows.Next() {
		var item StudentGuardian
		err := rows.Scan(
			&item.ID,
			&item.GuardianID,
			&item.SiswaID,
			&item.IsPrimary,
			&item.CreatedAt,
			&item.UpdateAt,
			&item.DeleteAt,
			&item.NamaSiswa,
			&item.NISN,
			&item.NamaKelas,
			&item.NamaSekolah,
			&item.NamaWali,
			&item.Hubungan,
			&item.TeleponWali,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan wali siswa: %w", err)
		}
		list = append(list, item)
	}

	return toStudentGuardianCoreList(list), nil
}
