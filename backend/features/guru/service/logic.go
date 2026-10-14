package service

import (
	"context"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/guru"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errGuruNotFound = errors.New("guru service: Data tidak ditemukan")

type guruService struct {
	guruData guru.DataGuruInterface
	db       *pgxpool.Pool
}

func NewServiceGuru(repo guru.DataGuruInterface, db *pgxpool.Pool) guru.ServiceGuruInterface {
	if repo == nil {
		panic("guru service: Nil repository")
	}
	return &guruService{guruData: repo, db: db}
}

func (s *guruService) GetAllGuru(schoolUnitID ...string) ([]guru.GuruCore, error) {
	if s.guruData == nil {
		return nil, errors.New("guru service: Nil repository")
	}
	gurus, err := s.guruData.SelectAllGuru(schoolUnitID...)
	if err != nil {
		return nil, errors.New("guru service: gagal mengambil data")
	}
	return gurus, nil
}

func validateTanggal(dateStr string) error {
	trimmed := strings.TrimSpace(dateStr)
	if trimmed == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", trimmed)
	if err != nil {
		return errors.New("format tanggal lahir harus YYYY-MM-DD")
	}
	if t.After(time.Now()) {
		return errors.New("tanggal lahir tidak boleh di masa depan")
	}
	return nil
}

func validateJenisKelamin(jk string) error {
	trimmed := strings.TrimSpace(jk)
	if trimmed == "" {
		return nil
	}
	lower := strings.ToLower(trimmed)
	if lower != "laki-laki" && lower != "perempuan" && lower != "l" && lower != "p" {
		return errors.New("jenis kelamin harus 'Laki-laki' atau 'Perempuan'")
	}
	return nil
}

func (s *guruService) InsertGuru(insert *guru.GuruCore) error {
	if s.guruData == nil {
		return errors.New("guru service: Repository kosong")
	}

	insert.Nama = strings.TrimSpace(insert.Nama)
	insert.Email = strings.TrimSpace(insert.Email)
	insert.Alamat = strings.TrimSpace(insert.Alamat)

	if insert.Nama == "" || insert.Email == "" || insert.Alamat == "" {
		return errors.New("validasi error: nama, email, dan alamat harus diisi")
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(strings.ToLower(insert.Email)) {
		return errors.New("validasi error: email tidak valid")
	}

	if err := validateTanggal(insert.TanggalLahir); err != nil {
		return err
	}
	if err := validateJenisKelamin(insert.JenisKelamin); err != nil {
		return err
	}

	if insert.ID_User == "" {
		var idUser string
		err := s.db.QueryRow(context.Background(),
			"SELECT id FROM users WHERE email = $1", insert.Email).Scan(&idUser)
		if err == nil {
			insert.ID_User = idUser
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("gagal cek users: %w", err)
		}
	}

	return s.guruData.InsertGuru(insert)
}

func (s *guruService) UpdateGuru(insert *guru.GuruCore, id string, schoolUnitID ...string) error {
	if s == nil || s.guruData == nil {
		return errors.New("guru service: Nil repository")
	}
	if insert == nil {
		return errors.New("guru service: input is nil")
	}
	if id == "" {
		return errors.New("validation error: id harus diisi")
	}

	existingData, err := s.guruData.SelectById(id, schoolUnitID...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errGuruNotFound
		}
		return fmt.Errorf("Id salah atau gagal mengambil data lama: %w", err)
	}

	if insert.Nama == "" {
		insert.Nama = existingData.Nama
	}
	if insert.Email == "" {
		insert.Email = existingData.Email
	}
	if insert.Alamat == "" {
		insert.Alamat = existingData.Alamat
	}

	if err := validateTanggal(insert.TanggalLahir); err != nil {
		return err
	}
	if err := validateJenisKelamin(insert.JenisKelamin); err != nil {
		return err
	}

	if err := s.guruData.Update(insert, id, schoolUnitID...); err != nil {
		return err
	}
	return nil
}

func (s *guruService) SelectById(id string, schoolUnitID ...string) (*guru.GuruCore, error) {
	guruData, err := s.guruData.SelectById(id, schoolUnitID...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errGuruNotFound
		}
		return nil, fmt.Errorf("gagal mengambil data guru: %w", err)
	}
	return guruData, nil
}

func (s *guruService) DeleteById(id string, schoolUnitID ...string) error {
	if s == nil || s.guruData == nil {
		return errors.New("guru service: Nil repository")
	}
	if id == "" {
		return errors.New("validation error: id harus diisi")
	}
	if err := s.guruData.DeleteById(id, schoolUnitID...); err != nil {
		return fmt.Errorf("gagal menghapus data guru: %w", err)
	}
	return nil
}
