package service

import (
	"errors"
	"fmt"
	schoolunits "go_rest_native_sekolah/features/school_units"
	"strings"
)

type SchoolUnitService struct {
	repo schoolunits.DataSchoolUnitInterface
}

func NewSchoolUnitService(repo schoolunits.DataSchoolUnitInterface) schoolunits.ServiceSchoolUnitInterface {
	if repo == nil {
		panic("nil school unit repo")
	}
	return &SchoolUnitService{repo: repo}
}

var validJenjang = map[string]bool{
	"SD":  true,
	"SMP": true,
	"SMA": true,
	"SMK": true,
}

var validStatus = map[string]bool{
	"aktif":              true,
	"nonaktif":           true,
	"dalam_pembangunan": true,
}

func (s *SchoolUnitService) validateInput(input *schoolunits.SchoolUnitCore, isCreate bool, currentID string) error {
	input.NamaSekolah = strings.TrimSpace(input.NamaSekolah)
	if input.NamaSekolah == "" {
		return errors.New("nama sekolah wajib diisi")
	}

	input.Jenjang = strings.ToUpper(strings.TrimSpace(input.Jenjang))
	if input.Jenjang == "" {
		return errors.New("jenjang wajib diisi")
	}
	if !validJenjang[input.Jenjang] {
		return errors.New("jenjang tidak valid: harus SD, SMP, SMA, atau SMK")
	}

	input.Status = strings.ToLower(strings.TrimSpace(input.Status))
	if input.Status == "" {
		input.Status = "aktif"
	}
	if !validStatus[input.Status] {
		return errors.New("status tidak valid: harus aktif, nonaktif, atau dalam_pembangunan")
	}

	input.NPSN = strings.TrimSpace(input.NPSN)
	if input.NPSN != "" {
		exists, err := s.repo.CheckNPSNExists(input.NPSN, currentID)
		if err != nil {
			return fmt.Errorf("gagal memvalidasi keunikan NPSN: %w", err)
		}
		if exists {
			return errors.New("NPSN sudah digunakan oleh unit sekolah lain")
		}
	}

	return nil
}

func (s *SchoolUnitService) GetAll(jenjang, status string) ([]schoolunits.SchoolUnitCore, error) {
	return s.repo.SelectAll(jenjang, status)
}

func (s *SchoolUnitService) GetByID(id string) (*schoolunits.SchoolUnitCore, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("parameter id wajib diisi")
	}
	return s.repo.SelectByID(strings.TrimSpace(id))
}

func (s *SchoolUnitService) Create(input *schoolunits.SchoolUnitCore) (*schoolunits.SchoolUnitCore, error) {
	if input == nil {
		return nil, errors.New("payload data unit sekolah tidak boleh kosong")
	}

	if err := s.validateInput(input, true, ""); err != nil {
		return nil, err
	}

	if err := s.repo.Insert(input); err != nil {
		return nil, err
	}

	return s.repo.SelectByID(input.ID)
}

func (s *SchoolUnitService) Update(id string, input *schoolunits.SchoolUnitCore) (*schoolunits.SchoolUnitCore, error) {
	trimID := strings.TrimSpace(id)
	if trimID == "" {
		return nil, errors.New("parameter id wajib diisi")
	}
	if input == nil {
		return nil, errors.New("payload data unit sekolah tidak boleh kosong")
	}

	existing, err := s.repo.SelectByID(trimID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(input.NamaSekolah) == "" {
		input.NamaSekolah = existing.NamaSekolah
	}
	if strings.TrimSpace(input.Jenjang) == "" {
		input.Jenjang = existing.Jenjang
	}
	if strings.TrimSpace(input.Status) == "" {
		input.Status = existing.Status
	}
	if strings.TrimSpace(input.NPSN) == "" {
		input.NPSN = existing.NPSN
	}
	if strings.TrimSpace(input.Alamat) == "" {
		input.Alamat = existing.Alamat
	}
	if strings.TrimSpace(input.Telepon) == "" {
		input.Telepon = existing.Telepon
	}
	if strings.TrimSpace(input.Email) == "" {
		input.Email = existing.Email
	}
	if strings.TrimSpace(input.Website) == "" {
		input.Website = existing.Website
	}

	if err := s.validateInput(input, false, trimID); err != nil {
		return nil, err
	}

	if err := s.repo.Update(trimID, input); err != nil {
		return nil, err
	}

	return s.repo.SelectByID(trimID)
}

func (s *SchoolUnitService) Delete(id string) error {
	trimID := strings.TrimSpace(id)
	if trimID == "" {
		return errors.New("parameter id wajib diisi")
	}

	if _, err := s.repo.SelectByID(trimID); err != nil {
		return err
	}

	return s.repo.Delete(trimID)
}
