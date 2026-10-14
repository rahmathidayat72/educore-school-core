package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	schoolunits "go_rest_native_sekolah/features/school_units"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SchoolUnitModel struct {
	db *pgxpool.Pool
}

func NewSchoolUnitModel(db *pgxpool.Pool) schoolunits.DataSchoolUnitInterface {
	if db == nil {
		panic("nil db pool on school unit model")
	}
	return &SchoolUnitModel{db: db}
}

func toNullString(s string) sql.NullString {
	trim := strings.TrimSpace(s)
	if trim == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: trim, Valid: true}
}

func (m *SchoolUnitModel) SelectAll(jenjang, status string) ([]schoolunits.SchoolUnitCore, error) {
	ctx := context.Background()
	queryBuilder := `
		SELECT 
			id, nama_sekolah, jenjang, npsn, alamat, telepon, email, website, status,
			created_at, update_at, delete_at
		FROM school_units
		WHERE delete_at IS NULL
	`
	var args []interface{}
	argIdx := 1

	if strings.TrimSpace(jenjang) != "" {
		queryBuilder += fmt.Sprintf(" AND UPPER(jenjang) = UPPER($%d)", argIdx)
		args = append(args, strings.TrimSpace(jenjang))
		argIdx++
	}

	if strings.TrimSpace(status) != "" {
		queryBuilder += fmt.Sprintf(" AND LOWER(status) = LOWER($%d)", argIdx)
		args = append(args, strings.TrimSpace(status))
		argIdx++
	}

	queryBuilder += " ORDER BY created_at ASC"

	rows, err := m.db.Query(ctx, queryBuilder, args...)
	if err != nil {
		return nil, fmt.Errorf("gagal query daftar unit sekolah: %w", err)
	}
	defer rows.Close()

	var list []SchoolUnit
	for rows.Next() {
		var item SchoolUnit
		err := rows.Scan(
			&item.ID,
			&item.NamaSekolah,
			&item.Jenjang,
			&item.NPSN,
			&item.Alamat,
			&item.Telepon,
			&item.Email,
			&item.Website,
			&item.Status,
			&item.CreatedAt,
			&item.UpdateAt,
			&item.DeleteAt,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan data unit sekolah: %w", err)
		}
		list = append(list, item)
	}

	return toCoreList(list), nil
}

func (m *SchoolUnitModel) SelectByID(id string) (*schoolunits.SchoolUnitCore, error) {
	ctx := context.Background()
	query := `
		SELECT 
			id, nama_sekolah, jenjang, npsn, alamat, telepon, email, website, status,
			created_at, update_at, delete_at
		FROM school_units
		WHERE id = $1 AND delete_at IS NULL
	`
	var item SchoolUnit
	err := m.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.NamaSekolah,
		&item.Jenjang,
		&item.NPSN,
		&item.Alamat,
		&item.Telepon,
		&item.Email,
		&item.Website,
		&item.Status,
		&item.CreatedAt,
		&item.UpdateAt,
		&item.DeleteAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("unit sekolah tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil detail unit sekolah: %w", err)
	}

	core := item.toCore()
	return &core, nil
}

func (m *SchoolUnitModel) Insert(input *schoolunits.SchoolUnitCore) error {
	ctx := context.Background()
	if strings.TrimSpace(input.ID) == "" {
		input.ID = uuid.New().String()
	}

	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = "aktif"
	}

	query := `
		INSERT INTO school_units (
			id, nama_sekolah, jenjang, npsn, alamat, telepon, email, website, status,
			created_at, update_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9,
			CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
	`
	_, err := m.db.Exec(ctx, query,
		input.ID,
		strings.TrimSpace(input.NamaSekolah),
		strings.ToUpper(strings.TrimSpace(input.Jenjang)),
		toNullString(input.NPSN),
		toNullString(input.Alamat),
		toNullString(input.Telepon),
		toNullString(input.Email),
		toNullString(input.Website),
		status,
	)
	if err != nil {
		if strings.Contains(err.Error(), "idx_school_units_npsn_unique") {
			return errors.New("NPSN sudah terdaftar pada unit sekolah lain")
		}
		return fmt.Errorf("gagal menyimpan data unit sekolah: %w", err)
	}

	return nil
}

func (m *SchoolUnitModel) Update(id string, input *schoolunits.SchoolUnitCore) error {
	ctx := context.Background()

	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = "aktif"
	}

	query := `
		UPDATE school_units
		SET nama_sekolah = $1,
		    jenjang = $2,
		    npsn = $3,
		    alamat = $4,
		    telepon = $5,
		    email = $6,
		    website = $7,
		    status = $8,
		    update_at = CURRENT_TIMESTAMP
		WHERE id = $9 AND delete_at IS NULL
	`
	res, err := m.db.Exec(ctx, query,
		strings.TrimSpace(input.NamaSekolah),
		strings.ToUpper(strings.TrimSpace(input.Jenjang)),
		toNullString(input.NPSN),
		toNullString(input.Alamat),
		toNullString(input.Telepon),
		toNullString(input.Email),
		toNullString(input.Website),
		status,
		id,
	)
	if err != nil {
		if strings.Contains(err.Error(), "idx_school_units_npsn_unique") {
			return errors.New("NPSN sudah terdaftar pada unit sekolah lain")
		}
		return fmt.Errorf("gagal memperbarui data unit sekolah: %w", err)
	}

	if res.RowsAffected() == 0 {
		return errors.New("unit sekolah tidak ditemukan atau tidak ada perubahan")
	}

	return nil
}

func (m *SchoolUnitModel) Delete(id string) error {
	ctx := context.Background()
	query := `
		UPDATE school_units
		SET delete_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND delete_at IS NULL
	`
	res, err := m.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus unit sekolah: %w", err)
	}
	if res.RowsAffected() == 0 {
		return errors.New("unit sekolah tidak ditemukan atau sudah dihapus")
	}
	return nil
}

func (m *SchoolUnitModel) CheckNPSNExists(npsn string, excludeID string) (bool, error) {
	ctx := context.Background()
	trimNPSN := strings.TrimSpace(npsn)
	if trimNPSN == "" {
		return false, nil
	}

	query := `
		SELECT COUNT(1)
		FROM school_units
		WHERE npsn = $1 
		  AND delete_at IS NULL 
		  AND ($2 = '' OR id != $2)
	`
	var count int
	err := m.db.QueryRow(ctx, query, trimNPSN, strings.TrimSpace(excludeID)).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
