package service

import (
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/siswa"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type siswaService struct {
	siswaData siswa.DataSiswaInterface
}

func NewServiceSiswa(repo siswa.DataSiswaInterface) siswa.ServiceSiswaInterface {
	if repo == nil {
		panic("Nil repository")
	}
	return &siswaService{siswaData: repo}
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

func (s *siswaService) InsertSiswa(insert *siswa.SiswaCore) error {
	if s.siswaData == nil {
		return errors.New("Nil repository")
	}

	insert.Nama = strings.TrimSpace(insert.Nama)
	insert.Alamat = strings.TrimSpace(insert.Alamat)
	insert.Email = strings.TrimSpace(insert.Email)

	if insert.Nama == "" {
		return errors.New("Validation error: nama siswa tidak boleh kosong")
	}
	if insert.Alamat == "" {
		return errors.New("Validation error: alamat siswa tidak boleh kosong")
	}
	if insert.Email == "" {
		return errors.New("Validation error: email siswa tidak boleh kosong")
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(strings.ToLower(insert.Email)) {
		return errors.New("validation error: email tidak valid")
	}

	if err := validateTanggal(insert.TanggalLahir); err != nil {
		return err
	}
	if err := validateJenisKelamin(insert.JenisKelamin); err != nil {
		return err
	}

	return s.siswaData.InsertSiswa(insert)
}

func (s *siswaService) SelectAllSiswa(schoolUnitID ...string) ([]siswa.SiswaCore, error) {
	if s.siswaData == nil {
		return nil, errors.New("SiswaService: Nil repository")
	}
	siswas, err := s.siswaData.SelectAllSiswa(schoolUnitID...)
	if err != nil {
		return nil, errors.New("SiswaService: gagal mengambil data")
	}
	return siswas, nil
}

func (s *siswaService) SelectById(id string, schoolUnitID ...string) (*siswa.SiswaCore, error) {
	siswa, err := s.siswaData.SelectById(id, schoolUnitID...)
	if err != nil {
		return nil, errors.New("SiswaService: gagal mengambil data")
	}
	return siswa, nil
}

func (s *siswaService) Update(insert *siswa.SiswaCore, id string, schoolUnitID ...string) error {
	if s == nil || s.siswaData == nil {
		return errors.New("Nil repository")
	}
	if id == "" {
		return errors.New("Validation error: id is nil")
	}

	existingData, err := s.siswaData.SelectById(id, schoolUnitID...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("siswa service: Data tidak ditemukan")
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
	if insert.Kelas_ID == "" {
		insert.Kelas_ID = existingData.Kelas_ID
	}

	if err := validateTanggal(insert.TanggalLahir); err != nil {
		return err
	}
	if err := validateJenisKelamin(insert.JenisKelamin); err != nil {
		return err
	}

	if err := s.siswaData.Update(insert, id, schoolUnitID...); err != nil {
		return err
	}
	return nil
}

func (s *siswaService) DeleteById(id string, schoolUnitID ...string) error {
	if s == nil || s.siswaData == nil {
		return errors.New("Nil repository")
	}
	if id == "" {
		return errors.New("validation error: id harus diisi")
	}
	if err := s.siswaData.DeleteById(id, schoolUnitID...); err != nil {
		return fmt.Errorf("gagal menghapus data siswa: %w", err)
	}
	return nil
}
