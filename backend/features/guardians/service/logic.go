package service

import (
	"errors"
	"go_rest_native_sekolah/features/guardians"
	"strings"
)

type GuardianService struct {
	repo guardians.DataGuardianInterface
}

func NewGuardianService(repo guardians.DataGuardianInterface) guardians.ServiceGuardianInterface {
	if repo == nil {
		panic("nil repo on guardian service")
	}
	return &GuardianService{repo: repo}
}

func normalizeHubungan(h string) string {
	trim := strings.TrimSpace(h)
	if trim == "" {
		return "Wali"
	}
	// Normalisasi kapitalisasi awal
	switch strings.ToLower(trim) {
	case "ayah", "bapak":
		return "Ayah"
	case "ibu", "mama":
		return "Ibu"
	case "kakak":
		return "Kakak"
	default:
		return trim
	}
}

func (s *GuardianService) GetAll(search string) ([]guardians.GuardianCore, error) {
	return s.repo.SelectAll(search)
}

func (s *GuardianService) GetByID(id string) (*guardians.GuardianCore, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id wali murid wajib diisi")
	}
	data, err := s.repo.SelectByID(id)
	if err != nil {
		return nil, err
	}

	// Enreach dengan daftar anak terhubung (BUS-005, FR-008)
	students, _ := s.repo.SelectStudentsByGuardianID(id)
	data.SiswaList = students

	return data, nil
}

func (s *GuardianService) GetByUserID(userID string) (*guardians.GuardianCore, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("id user wajib diisi")
	}
	data, err := s.repo.SelectByUserID(userID)
	if err != nil {
		return nil, err
	}

	students, _ := s.repo.SelectStudentsByGuardianID(data.ID)
	data.SiswaList = students

	return data, nil
}

func (s *GuardianService) Create(input *guardians.GuardianCore) (*guardians.GuardianCore, error) {
	if input == nil {
		return nil, errors.New("data wali murid tidak boleh kosong")
	}
	if strings.TrimSpace(input.Nama) == "" {
		return nil, errors.New("nama wali murid wajib diisi")
	}
	if strings.TrimSpace(input.Telepon) == "" {
		return nil, errors.New("nomor telepon wali murid wajib diisi")
	}

	input.Hubungan = normalizeHubungan(input.Hubungan)

	if err := s.repo.Insert(input); err != nil {
		return nil, err
	}
	return input, nil
}

func (s *GuardianService) Update(id string, input *guardians.GuardianCore) (*guardians.GuardianCore, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id wali murid wajib diisi")
	}
	if input == nil {
		return nil, errors.New("data update wali murid tidak boleh kosong")
	}
	if strings.TrimSpace(input.Nama) == "" {
		return nil, errors.New("nama wali murid wajib diisi")
	}
	if strings.TrimSpace(input.Telepon) == "" {
		return nil, errors.New("nomor telepon wali murid wajib diisi")
	}

	input.Hubungan = normalizeHubungan(input.Hubungan)

	if err := s.repo.Update(id, input); err != nil {
		return nil, err
	}
	input.ID = id
	return input, nil
}

func (s *GuardianService) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id wali murid wajib diisi")
	}
	return s.repo.Delete(id)
}

func (s *GuardianService) LinkStudent(guardianID string, siswaID string, isPrimary bool) (*guardians.StudentGuardianCore, error) {
	if strings.TrimSpace(guardianID) == "" {
		return nil, errors.New("id wali murid wajib diisi")
	}
	if strings.TrimSpace(siswaID) == "" {
		return nil, errors.New("id siswa wajib diisi")
	}

	link := &guardians.StudentGuardianCore{
		GuardianID: guardianID,
		SiswaID:    siswaID,
		IsPrimary:  isPrimary,
	}

	if err := s.repo.LinkStudent(link); err != nil {
		return nil, err
	}
	return link, nil
}

func (s *GuardianService) UnlinkStudent(guardianID string, siswaID string) error {
	if strings.TrimSpace(guardianID) == "" {
		return errors.New("id wali murid wajib diisi")
	}
	if strings.TrimSpace(siswaID) == "" {
		return errors.New("id siswa wajib diisi")
	}
	return s.repo.UnlinkStudent(guardianID, siswaID)
}

func (s *GuardianService) GetStudentsByGuardianID(guardianID string) ([]guardians.StudentGuardianCore, error) {
	if strings.TrimSpace(guardianID) == "" {
		return nil, errors.New("id wali murid wajib diisi")
	}
	return s.repo.SelectStudentsByGuardianID(guardianID)
}

func (s *GuardianService) GetGuardiansBySiswaID(siswaID string) ([]guardians.StudentGuardianCore, error) {
	if strings.TrimSpace(siswaID) == "" {
		return nil, errors.New("id siswa wajib diisi")
	}
	return s.repo.SelectGuardiansBySiswaID(siswaID)
}
