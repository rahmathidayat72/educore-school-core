package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	academicyears "go_rest_native_sekolah/features/academic_years"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AcademicYearModel struct {
	db *pgxpool.Pool
}

func NewAcademicYearModel(db *pgxpool.Pool) academicyears.DataAcademicYearInterface {
	if db == nil {
		panic("nil db pool on academic year model")
	}
	return &AcademicYearModel{db: db}
}

func toNullDate(s *string) sql.NullString {
	if s == nil || strings.TrimSpace(*s) == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: strings.TrimSpace(*s), Valid: true}
}

func (m *AcademicYearModel) SelectAllYears(schoolUnitID string) ([]academicyears.AcademicYearCore, error) {
	ctx := context.Background()
	query := `
		SELECT 
			id, school_unit_id, tahun_ajaran, semester, is_active,
			tanggal_mulai::text, tanggal_selesai::text,
			created_at, update_at, delete_at
		FROM academic_years
		WHERE delete_at IS NULL
	`
	var args []interface{}
	if strings.TrimSpace(schoolUnitID) != "" {
		query += " AND school_unit_id = $1"
		args = append(args, strings.TrimSpace(schoolUnitID))
	}
	query += " ORDER BY created_at DESC"

	rows, err := m.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("gagal query tahun ajaran: %w", err)
	}
	defer rows.Close()

	var list []AcademicYear
	for rows.Next() {
		var item AcademicYear
		err := rows.Scan(
			&item.ID,
			&item.SchoolUnitID,
			&item.TahunAjaran,
			&item.Semester,
			&item.IsActive,
			&item.TanggalMulai,
			&item.TanggalSelesai,
			&item.CreatedAt,
			&item.UpdateAt,
			&item.DeleteAt,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan tahun ajaran: %w", err)
		}
		list = append(list, item)
	}

	return toAcademicYearCoreList(list), nil
}

func (m *AcademicYearModel) SelectYearByID(id string, schoolUnitID ...string) (*academicyears.AcademicYearCore, error) {
	ctx := context.Background()
	query := `
		SELECT 
			id, school_unit_id, tahun_ajaran, semester, is_active,
			tanggal_mulai::text, tanggal_selesai::text,
			created_at, update_at, delete_at
		FROM academic_years
		WHERE id = $1 AND delete_at IS NULL
	`
	args := []interface{}{id}
	if len(schoolUnitID) > 0 && strings.TrimSpace(schoolUnitID[0]) != "" {
		query += " AND school_unit_id = $2"
		args = append(args, strings.TrimSpace(schoolUnitID[0]))
	}

	var item AcademicYear
	err := m.db.QueryRow(ctx, query, args...).Scan(
		&item.ID,
		&item.SchoolUnitID,
		&item.TahunAjaran,
		&item.Semester,
		&item.IsActive,
		&item.TanggalMulai,
		&item.TanggalSelesai,
		&item.CreatedAt,
		&item.UpdateAt,
		&item.DeleteAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("data tahun ajaran tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal query detail tahun ajaran: %w", err)
	}

	core := item.toCore()
	return &core, nil
}

func (m *AcademicYearModel) GetActiveYear(schoolUnitID string) (*academicyears.AcademicYearCore, error) {
	ctx := context.Background()
	query := `
		SELECT 
			id, school_unit_id, tahun_ajaran, semester, is_active,
			tanggal_mulai::text, tanggal_selesai::text,
			created_at, update_at, delete_at
		FROM academic_years
		WHERE is_active = true AND delete_at IS NULL
	`
	var args []interface{}
	if strings.TrimSpace(schoolUnitID) != "" {
		query += " AND school_unit_id = $1"
		args = append(args, strings.TrimSpace(schoolUnitID))
	}
	query += " LIMIT 1"

	var item AcademicYear
	err := m.db.QueryRow(ctx, query, args...).Scan(
		&item.ID,
		&item.SchoolUnitID,
		&item.TahunAjaran,
		&item.Semester,
		&item.IsActive,
		&item.TanggalMulai,
		&item.TanggalSelesai,
		&item.CreatedAt,
		&item.UpdateAt,
		&item.DeleteAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("tidak ada tahun ajaran aktif untuk unit ini")
		}
		return nil, fmt.Errorf("gagal query tahun ajaran aktif: %w", err)
	}

	core := item.toCore()
	return &core, nil
}

func (m *AcademicYearModel) InsertYear(input *academicyears.AcademicYearCore) error {
	ctx := context.Background()
	if strings.TrimSpace(input.ID) == "" {
		input.ID = uuid.New().String()
	}

	query := `
		INSERT INTO academic_years (
			id, school_unit_id, tahun_ajaran, semester, is_active,
			tanggal_mulai, tanggal_selesai, created_at, update_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	_, err := m.db.Exec(
		ctx,
		query,
		input.ID,
		input.SchoolUnitID,
		input.TahunAjaran,
		input.Semester,
		input.IsActive,
		toNullDate(input.TanggalMulai),
		toNullDate(input.TanggalSelesai),
	)
	if err != nil {
		return fmt.Errorf("gagal insert tahun ajaran: %w", err)
	}
	return nil
}

func (m *AcademicYearModel) UpdateYear(id string, input *academicyears.AcademicYearCore, schoolUnitID ...string) error {
	ctx := context.Background()
	query := `
		UPDATE academic_years
		SET tahun_ajaran = $1, semester = $2, tanggal_mulai = $3, tanggal_selesai = $4, update_at = CURRENT_TIMESTAMP
		WHERE id = $5 AND delete_at IS NULL
	`
	args := []interface{}{
		input.TahunAjaran,
		input.Semester,
		toNullDate(input.TanggalMulai),
		toNullDate(input.TanggalSelesai),
		id,
	}
	if len(schoolUnitID) > 0 && strings.TrimSpace(schoolUnitID[0]) != "" {
		query += " AND school_unit_id = $6"
		args = append(args, strings.TrimSpace(schoolUnitID[0]))
	}

	cmd, err := m.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("gagal update tahun ajaran: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("data tahun ajaran tidak ditemukan atau tidak berubah")
	}
	return nil
}

func (m *AcademicYearModel) SetActiveYear(id string, schoolUnitID string) error {
	ctx := context.Background()
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi set active year: %w", err)
	}
	defer tx.Rollback(ctx)

	// Nonaktifkan semua tahun ajaran di unit sekolah tersebut
	deactivateQuery := `
		UPDATE academic_years
		SET is_active = false, update_at = CURRENT_TIMESTAMP
		WHERE school_unit_id = $1 AND is_active = true AND delete_at IS NULL
	`
	if _, err := tx.Exec(ctx, deactivateQuery, schoolUnitID); err != nil {
		return fmt.Errorf("gagal menonaktifkan tahun ajaran aktif sebelumnya: %w", err)
	}

	// Aktifkan tahun ajaran yang dipilih
	activateQuery := `
		UPDATE academic_years
		SET is_active = true, update_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND school_unit_id = $2 AND delete_at IS NULL
	`
	cmd, err := tx.Exec(ctx, activateQuery, id, schoolUnitID)
	if err != nil {
		return fmt.Errorf("gagal mengaktifkan tahun ajaran: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("tahun ajaran yang akan diaktifkan tidak ditemukan pada unit ini")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("gagal commit aktivasi tahun ajaran: %w", err)
	}
	return nil
}

func (m *AcademicYearModel) DeleteYear(id string, schoolUnitID ...string) error {
	ctx := context.Background()
	query := `
		UPDATE academic_years
		SET delete_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND delete_at IS NULL
	`
	args := []interface{}{id}
	if len(schoolUnitID) > 0 && strings.TrimSpace(schoolUnitID[0]) != "" {
		query += " AND school_unit_id = $2"
		args = append(args, strings.TrimSpace(schoolUnitID[0]))
	}

	cmd, err := m.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("gagal menghapus tahun ajaran: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("data tahun ajaran tidak ditemukan")
	}
	return nil
}

// Enrollments Implementation

func (m *AcademicYearModel) SelectEnrollments(schoolUnitID string, academicYearID string, kelasID string) ([]academicyears.ClassEnrollmentCore, error) {
	ctx := context.Background()
	queryBuilder := `
		SELECT 
			e.id, e.school_unit_id, e.academic_year_id, e.kelas_id, e.siswa_id, e.status, e.keterangan,
			e.created_at, e.update_at, e.delete_at,
			s.nama as nama_siswa, s.nis, s.nisn,
			k.nama_kelas, k.tingkat,
			ay.tahun_ajaran, ay.semester
		FROM class_enrollments e
		LEFT JOIN siswa s ON e.siswa_id = s.id
		LEFT JOIN kelas k ON e.kelas_id = k.id
		LEFT JOIN academic_years ay ON e.academic_year_id = ay.id
		WHERE e.delete_at IS NULL
	`
	var args []interface{}
	argIdx := 1

	if strings.TrimSpace(schoolUnitID) != "" {
		queryBuilder += fmt.Sprintf(" AND e.school_unit_id = $%d", argIdx)
		args = append(args, strings.TrimSpace(schoolUnitID))
		argIdx++
	}

	if strings.TrimSpace(academicYearID) != "" {
		queryBuilder += fmt.Sprintf(" AND e.academic_year_id = $%d", argIdx)
		args = append(args, strings.TrimSpace(academicYearID))
		argIdx++
	}

	if strings.TrimSpace(kelasID) != "" {
		queryBuilder += fmt.Sprintf(" AND e.kelas_id = $%d", argIdx)
		args = append(args, strings.TrimSpace(kelasID))
		argIdx++
	}

	queryBuilder += " ORDER BY k.nama_kelas ASC, s.nama ASC"

	rows, err := m.db.Query(ctx, queryBuilder, args...)
	if err != nil {
		return nil, fmt.Errorf("gagal query daftar penempatan kelas: %w", err)
	}
	defer rows.Close()

	var list []ClassEnrollment
	for rows.Next() {
		var item ClassEnrollment
		err := rows.Scan(
			&item.ID,
			&item.SchoolUnitID,
			&item.AcademicYearID,
			&item.KelasID,
			&item.SiswaID,
			&item.Status,
			&item.Keterangan,
			&item.CreatedAt,
			&item.UpdateAt,
			&item.DeleteAt,
			&item.NamaSiswa,
			&item.NIS,
			&item.NISN,
			&item.NamaKelas,
			&item.Tingkat,
			&item.TahunAjaran,
			&item.Semester,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan data penempatan kelas: %w", err)
		}
		list = append(list, item)
	}

	return toEnrollmentCoreList(list), nil
}

func (m *AcademicYearModel) SelectEnrollmentByID(id string, schoolUnitID ...string) (*academicyears.ClassEnrollmentCore, error) {
	ctx := context.Background()
	query := `
		SELECT 
			e.id, e.school_unit_id, e.academic_year_id, e.kelas_id, e.siswa_id, e.status, e.keterangan,
			e.created_at, e.update_at, e.delete_at,
			s.nama as nama_siswa, s.nis, s.nisn,
			k.nama_kelas, k.tingkat,
			ay.tahun_ajaran, ay.semester
		FROM class_enrollments e
		LEFT JOIN siswa s ON e.siswa_id = s.id
		LEFT JOIN kelas k ON e.kelas_id = k.id
		LEFT JOIN academic_years ay ON e.academic_year_id = ay.id
		WHERE e.id = $1 AND e.delete_at IS NULL
	`
	args := []interface{}{id}
	if len(schoolUnitID) > 0 && strings.TrimSpace(schoolUnitID[0]) != "" {
		query += " AND e.school_unit_id = $2"
		args = append(args, strings.TrimSpace(schoolUnitID[0]))
	}

	var item ClassEnrollment
	err := m.db.QueryRow(ctx, query, args...).Scan(
		&item.ID,
		&item.SchoolUnitID,
		&item.AcademicYearID,
		&item.KelasID,
		&item.SiswaID,
		&item.Status,
		&item.Keterangan,
		&item.CreatedAt,
		&item.UpdateAt,
		&item.DeleteAt,
		&item.NamaSiswa,
		&item.NIS,
		&item.NISN,
		&item.NamaKelas,
		&item.Tingkat,
		&item.TahunAjaran,
		&item.Semester,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("data penempatan kelas tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal query penempatan kelas: %w", err)
	}

	core := item.toCore()
	return &core, nil
}

func (m *AcademicYearModel) InsertEnrollment(input *academicyears.ClassEnrollmentCore) error {
	ctx := context.Background()
	if strings.TrimSpace(input.ID) == "" {
		input.ID = uuid.New().String()
	}

	query := `
		INSERT INTO class_enrollments (
			id, school_unit_id, academic_year_id, kelas_id, siswa_id,
			status, keterangan, created_at, update_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	var ket sql.NullString
	if strings.TrimSpace(input.Keterangan) != "" {
		ket = sql.NullString{String: strings.TrimSpace(input.Keterangan), Valid: true}
	}

	_, err := m.db.Exec(
		ctx,
		query,
		input.ID,
		input.SchoolUnitID,
		input.AcademicYearID,
		input.KelasID,
		input.SiswaID,
		input.Status,
		ket,
	)
	if err != nil {
		if strings.Contains(err.Error(), "idx_class_enrollments_unique_student_year") {
			return errors.New("siswa sudah terdaftar pada rombel di tahun ajaran ini")
		}
		return fmt.Errorf("gagal insert penempatan kelas: %w", err)
	}
	return nil
}

func (m *AcademicYearModel) UpdateEnrollment(id string, input *academicyears.ClassEnrollmentCore, schoolUnitID ...string) error {
	ctx := context.Background()
	query := `
		UPDATE class_enrollments
		SET kelas_id = $1, status = $2, keterangan = $3, update_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND delete_at IS NULL
	`
	var ket sql.NullString
	if strings.TrimSpace(input.Keterangan) != "" {
		ket = sql.NullString{String: strings.TrimSpace(input.Keterangan), Valid: true}
	}

	args := []interface{}{
		input.KelasID,
		input.Status,
		ket,
		id,
	}
	if len(schoolUnitID) > 0 && strings.TrimSpace(schoolUnitID[0]) != "" {
		query += " AND school_unit_id = $5"
		args = append(args, strings.TrimSpace(schoolUnitID[0]))
	}

	cmd, err := m.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("gagal update penempatan kelas: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("data penempatan kelas tidak ditemukan")
	}
	return nil
}

func (m *AcademicYearModel) DeleteEnrollment(id string, schoolUnitID ...string) error {
	ctx := context.Background()
	query := `
		UPDATE class_enrollments
		SET delete_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND delete_at IS NULL
	`
	args := []interface{}{id}
	if len(schoolUnitID) > 0 && strings.TrimSpace(schoolUnitID[0]) != "" {
		query += " AND school_unit_id = $2"
		args = append(args, strings.TrimSpace(schoolUnitID[0]))
	}

	cmd, err := m.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("gagal menghapus penempatan kelas: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("data penempatan kelas tidak ditemukan")
	}
	return nil
}

func (m *AcademicYearModel) GetStudentEnrollmentHistory(siswaID string) ([]academicyears.ClassEnrollmentCore, error) {
	ctx := context.Background()
	query := `
		SELECT 
			e.id, e.school_unit_id, e.academic_year_id, e.kelas_id, e.siswa_id, e.status, e.keterangan,
			e.created_at, e.update_at, e.delete_at,
			s.nama as nama_siswa, s.nis, s.nisn,
			k.nama_kelas, k.tingkat,
			ay.tahun_ajaran, ay.semester
		FROM class_enrollments e
		LEFT JOIN siswa s ON e.siswa_id = s.id
		LEFT JOIN kelas k ON e.kelas_id = k.id
		LEFT JOIN academic_years ay ON e.academic_year_id = ay.id
		WHERE e.siswa_id = $1 AND e.delete_at IS NULL
		ORDER BY ay.tahun_ajaran DESC, ay.semester DESC, e.created_at DESC
	`
	rows, err := m.db.Query(ctx, query, siswaID)
	if err != nil {
		return nil, fmt.Errorf("gagal query riwayat penempatan siswa: %w", err)
	}
	defer rows.Close()

	var list []ClassEnrollment
	for rows.Next() {
		var item ClassEnrollment
		err := rows.Scan(
			&item.ID,
			&item.SchoolUnitID,
			&item.AcademicYearID,
			&item.KelasID,
			&item.SiswaID,
			&item.Status,
			&item.Keterangan,
			&item.CreatedAt,
			&item.UpdateAt,
			&item.DeleteAt,
			&item.NamaSiswa,
			&item.NIS,
			&item.NISN,
			&item.NamaKelas,
			&item.Tingkat,
			&item.TahunAjaran,
			&item.Semester,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan riwayat penempatan siswa: %w", err)
		}
		list = append(list, item)
	}

	return toEnrollmentCoreList(list), nil
}
