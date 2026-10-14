package service

import (
	"errors"
	"fmt"
	academicyears "go_rest_native_sekolah/features/academic_years"
	"strings"
)

type AcademicYearService struct {
	repo academicyears.DataAcademicYearInterface
}

func NewAcademicYearService(repo academicyears.DataAcademicYearInterface) academicyears.ServiceAcademicYearInterface {
	if repo == nil {
		panic("nil repo on academic year service")
	}
	return &AcademicYearService{repo: repo}
}

func validateSemester(sem string) (string, error) {
	norm := strings.TrimSpace(sem)
	if strings.EqualFold(norm, "Ganjil") {
		return "Ganjil", nil
	}
	if strings.EqualFold(norm, "Genap") {
		return "Genap", nil
	}
	return "", errors.New("semester harus 'Ganjil' atau 'Genap'")
}

func validateEnrollmentStatus(status string) (string, error) {
	norm := strings.TrimSpace(status)
	if norm == "" {
		return "Aktif", nil
	}
	allowed := []string{"Aktif", "Lulus", "Pindah", "Tinggal Kelas"}
	for _, val := range allowed {
		if strings.EqualFold(norm, val) {
			return val, nil
		}
	}
	return "", fmt.Errorf("status pendaftaran tidak valid: harus salah satu dari %s", strings.Join(allowed, ", "))
}

func (s *AcademicYearService) GetAllYears(schoolUnitID string) ([]academicyears.AcademicYearCore, error) {
	return s.repo.SelectAllYears(schoolUnitID)
}

func (s *AcademicYearService) GetYearByID(id string, schoolUnitID ...string) (*academicyears.AcademicYearCore, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id tahun ajaran wajib diisi")
	}
	return s.repo.SelectYearByID(id, schoolUnitID...)
}

func (s *AcademicYearService) GetActiveYear(schoolUnitID string) (*academicyears.AcademicYearCore, error) {
	return s.repo.GetActiveYear(schoolUnitID)
}

func (s *AcademicYearService) CreateYear(input *academicyears.AcademicYearCore) (*academicyears.AcademicYearCore, error) {
	if input == nil {
		return nil, errors.New("data tahun ajaran tidak boleh kosong")
	}
	if strings.TrimSpace(input.SchoolUnitID) == "" {
		return nil, errors.New("unit sekolah (school_unit_id) wajib diisi")
	}
	if strings.TrimSpace(input.TahunAjaran) == "" {
		return nil, errors.New("nama tahun ajaran (contoh: 2026/2027) wajib diisi")
	}
	sem, err := validateSemester(input.Semester)
	if err != nil {
		return nil, err
	}
	input.Semester = sem

	// Jika ditandai aktif saat pembuatan
	shouldActivate := input.IsActive
	input.IsActive = false // insert awal sebagai false agar aman dari constraint unik

	if err := s.repo.InsertYear(input); err != nil {
		return nil, err
	}

	if shouldActivate {
		if err := s.repo.SetActiveYear(input.ID, input.SchoolUnitID); err != nil {
			return nil, fmt.Errorf("tahun ajaran tersimpan namun gagal diaktifkan: %w", err)
		}
		input.IsActive = true
	}

	return input, nil
}

func (s *AcademicYearService) UpdateYear(id string, input *academicyears.AcademicYearCore, schoolUnitID ...string) (*academicyears.AcademicYearCore, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id tahun ajaran wajib diisi")
	}
	if input == nil {
		return nil, errors.New("data update tahun ajaran tidak boleh kosong")
	}
	if strings.TrimSpace(input.TahunAjaran) == "" {
		return nil, errors.New("nama tahun ajaran wajib diisi")
	}
	sem, err := validateSemester(input.Semester)
	if err != nil {
		return nil, err
	}
	input.Semester = sem

	if err := s.repo.UpdateYear(id, input, schoolUnitID...); err != nil {
		return nil, err
	}
	input.ID = id
	return input, nil
}

func (s *AcademicYearService) SetActiveYear(id string, schoolUnitID string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id tahun ajaran wajib diisi")
	}
	if strings.TrimSpace(schoolUnitID) == "" {
		return errors.New("unit sekolah wajib diisi untuk aktivasi tahun ajaran")
	}
	return s.repo.SetActiveYear(id, schoolUnitID)
}

func (s *AcademicYearService) DeleteYear(id string, schoolUnitID ...string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id tahun ajaran wajib diisi")
	}

	// Cek apakah tahun ajaran sedang aktif
	existing, err := s.repo.SelectYearByID(id, schoolUnitID...)
	if err != nil {
		return err
	}
	if existing.IsActive {
		return errors.New("tahun ajaran yang sedang aktif tidak dapat dihapus; aktifkan tahun ajaran lain terlebih dahulu")
	}

	return s.repo.DeleteYear(id, schoolUnitID...)
}

// Enrollments

func (s *AcademicYearService) GetEnrollments(schoolUnitID string, academicYearID string, kelasID string) ([]academicyears.ClassEnrollmentCore, error) {
	return s.repo.SelectEnrollments(schoolUnitID, academicYearID, kelasID)
}

func (s *AcademicYearService) GetEnrollmentByID(id string, schoolUnitID ...string) (*academicyears.ClassEnrollmentCore, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id penempatan kelas wajib diisi")
	}
	return s.repo.SelectEnrollmentByID(id, schoolUnitID...)
}

func (s *AcademicYearService) EnrollStudent(input *academicyears.ClassEnrollmentCore) (*academicyears.ClassEnrollmentCore, error) {
	if input == nil {
		return nil, errors.New("data penempatan kelas tidak boleh kosong")
	}
	if strings.TrimSpace(input.SchoolUnitID) == "" {
		return nil, errors.New("unit sekolah (school_unit_id) wajib diisi")
	}
	if strings.TrimSpace(input.AcademicYearID) == "" {
		return nil, errors.New("tahun ajaran (academic_year_id) wajib diisi")
	}
	if strings.TrimSpace(input.KelasID) == "" {
		return nil, errors.New("kelas/rombel (kelas_id) wajib diisi")
	}
	if strings.TrimSpace(input.SiswaID) == "" {
		return nil, errors.New("siswa (siswa_id) wajib diisi")
	}

	status, err := validateEnrollmentStatus(input.Status)
	if err != nil {
		return nil, err
	}
	input.Status = status

	if err := s.repo.InsertEnrollment(input); err != nil {
		return nil, err
	}
	return input, nil
}

func (s *AcademicYearService) UpdateEnrollment(id string, input *academicyears.ClassEnrollmentCore, schoolUnitID ...string) (*academicyears.ClassEnrollmentCore, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id penempatan kelas wajib diisi")
	}
	if input == nil {
		return nil, errors.New("data update penempatan kelas tidak boleh kosong")
	}
	if strings.TrimSpace(input.KelasID) == "" {
		return nil, errors.New("kelas/rombel (kelas_id) wajib diisi")
	}
	status, err := validateEnrollmentStatus(input.Status)
	if err != nil {
		return nil, err
	}
	input.Status = status

	if err := s.repo.UpdateEnrollment(id, input, schoolUnitID...); err != nil {
		return nil, err
	}
	input.ID = id
	return input, nil
}

func (s *AcademicYearService) DeleteEnrollment(id string, schoolUnitID ...string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id penempatan kelas wajib diisi")
	}
	return s.repo.DeleteEnrollment(id, schoolUnitID...)
}

func (s *AcademicYearService) GetStudentEnrollmentHistory(siswaID string) ([]academicyears.ClassEnrollmentCore, error) {
	if strings.TrimSpace(siswaID) == "" {
		return nil, errors.New("id siswa wajib diisi")
	}
	return s.repo.GetStudentEnrollmentHistory(siswaID)
}
