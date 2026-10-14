package service

import (
	academicyears "go_rest_native_sekolah/features/academic_years"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAcademicYearRepo struct {
	mock.Mock
}

func (m *mockAcademicYearRepo) SelectAllYears(schoolUnitID string) ([]academicyears.AcademicYearCore, error) {
	args := m.Called(schoolUnitID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]academicyears.AcademicYearCore), args.Error(1)
}

func (m *mockAcademicYearRepo) SelectYearByID(id string, schoolUnitID ...string) (*academicyears.AcademicYearCore, error) {
	var unit string
	if len(schoolUnitID) > 0 {
		unit = schoolUnitID[0]
	}
	args := m.Called(id, unit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*academicyears.AcademicYearCore), args.Error(1)
}

func (m *mockAcademicYearRepo) GetActiveYear(schoolUnitID string) (*academicyears.AcademicYearCore, error) {
	args := m.Called(schoolUnitID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*academicyears.AcademicYearCore), args.Error(1)
}

func (m *mockAcademicYearRepo) InsertYear(input *academicyears.AcademicYearCore) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *mockAcademicYearRepo) UpdateYear(id string, input *academicyears.AcademicYearCore, schoolUnitID ...string) error {
	var unit string
	if len(schoolUnitID) > 0 {
		unit = schoolUnitID[0]
	}
	args := m.Called(id, input, unit)
	return args.Error(0)
}

func (m *mockAcademicYearRepo) SetActiveYear(id string, schoolUnitID string) error {
	args := m.Called(id, schoolUnitID)
	return args.Error(0)
}

func (m *mockAcademicYearRepo) DeleteYear(id string, schoolUnitID ...string) error {
	var unit string
	if len(schoolUnitID) > 0 {
		unit = schoolUnitID[0]
	}
	args := m.Called(id, unit)
	return args.Error(0)
}

func (m *mockAcademicYearRepo) SelectEnrollments(schoolUnitID string, academicYearID string, kelasID string) ([]academicyears.ClassEnrollmentCore, error) {
	args := m.Called(schoolUnitID, academicYearID, kelasID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]academicyears.ClassEnrollmentCore), args.Error(1)
}

func (m *mockAcademicYearRepo) SelectEnrollmentByID(id string, schoolUnitID ...string) (*academicyears.ClassEnrollmentCore, error) {
	var unit string
	if len(schoolUnitID) > 0 {
		unit = schoolUnitID[0]
	}
	args := m.Called(id, unit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*academicyears.ClassEnrollmentCore), args.Error(1)
}

func (m *mockAcademicYearRepo) InsertEnrollment(input *academicyears.ClassEnrollmentCore) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *mockAcademicYearRepo) UpdateEnrollment(id string, input *academicyears.ClassEnrollmentCore, schoolUnitID ...string) error {
	var unit string
	if len(schoolUnitID) > 0 {
		unit = schoolUnitID[0]
	}
	args := m.Called(id, input, unit)
	return args.Error(0)
}

func (m *mockAcademicYearRepo) DeleteEnrollment(id string, schoolUnitID ...string) error {
	var unit string
	if len(schoolUnitID) > 0 {
		unit = schoolUnitID[0]
	}
	args := m.Called(id, unit)
	return args.Error(0)
}

func (m *mockAcademicYearRepo) GetStudentEnrollmentHistory(siswaID string) ([]academicyears.ClassEnrollmentCore, error) {
	args := m.Called(siswaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]academicyears.ClassEnrollmentCore), args.Error(1)
}

func TestAcademicYearService_Years(t *testing.T) {
	mockRepo := new(mockAcademicYearRepo)
	svc := NewAcademicYearService(mockRepo)

	t.Run("CreateYear validation errors", func(t *testing.T) {
		_, err := svc.CreateYear(nil)
		assert.Error(t, err)

		_, err = svc.CreateYear(&academicyears.AcademicYearCore{})
		assert.ErrorContains(t, err, "unit sekolah")

		_, err = svc.CreateYear(&academicyears.AcademicYearCore{SchoolUnitID: "unit-1"})
		assert.ErrorContains(t, err, "nama tahun ajaran")

		_, err = svc.CreateYear(&academicyears.AcademicYearCore{
			SchoolUnitID: "unit-1",
			TahunAjaran:  "2026/2027",
			Semester:     "Caturwulan",
		})
		assert.ErrorContains(t, err, "semester harus 'Ganjil' atau 'Genap'")
	})

	t.Run("CreateYear success with auto activate", func(t *testing.T) {
		input := &academicyears.AcademicYearCore{
			ID:           "ay-1",
			SchoolUnitID: "unit-1",
			TahunAjaran:  "2026/2027",
			Semester:     "ganjil",
			IsActive:     true,
		}

		mockRepo.On("InsertYear", mock.MatchedBy(func(item *academicyears.AcademicYearCore) bool {
			return item.SchoolUnitID == "unit-1" && item.Semester == "Ganjil" && item.IsActive == false
		})).Return(nil).Once()

		mockRepo.On("SetActiveYear", "ay-1", "unit-1").Return(nil).Once()

		res, err := svc.CreateYear(input)
		assert.NoError(t, err)
		assert.True(t, res.IsActive)
		assert.Equal(t, "Ganjil", res.Semester)
	})

	t.Run("DeleteYear fails if active", func(t *testing.T) {
		mockRepo.On("SelectYearByID", "ay-active", "unit-1").Return(&academicyears.AcademicYearCore{
			ID:       "ay-active",
			IsActive: true,
		}, nil).Once()

		err := svc.DeleteYear("ay-active", "unit-1")
		assert.ErrorContains(t, err, "tahun ajaran yang sedang aktif tidak dapat dihapus")
	})

	t.Run("DeleteYear success if inactive", func(t *testing.T) {
		mockRepo.On("SelectYearByID", "ay-inactive", "unit-1").Return(&academicyears.AcademicYearCore{
			ID:       "ay-inactive",
			IsActive: false,
		}, nil).Once()

		mockRepo.On("DeleteYear", "ay-inactive", "unit-1").Return(nil).Once()

		err := svc.DeleteYear("ay-inactive", "unit-1")
		assert.NoError(t, err)
	})
}

func TestAcademicYearService_Enrollments(t *testing.T) {
	mockRepo := new(mockAcademicYearRepo)
	svc := NewAcademicYearService(mockRepo)

	t.Run("EnrollStudent validation errors", func(t *testing.T) {
		_, err := svc.EnrollStudent(nil)
		assert.Error(t, err)

		_, err = svc.EnrollStudent(&academicyears.ClassEnrollmentCore{})
		assert.ErrorContains(t, err, "unit sekolah")

		_, err = svc.EnrollStudent(&academicyears.ClassEnrollmentCore{
			SchoolUnitID: "unit-1",
		})
		assert.ErrorContains(t, err, "tahun ajaran")

		_, err = svc.EnrollStudent(&academicyears.ClassEnrollmentCore{
			SchoolUnitID:   "unit-1",
			AcademicYearID: "ay-1",
			KelasID:        "k-1",
			SiswaID:        "s-1",
			Status:         "DO",
		})
		assert.ErrorContains(t, err, "status pendaftaran tidak valid")
	})

	t.Run("EnrollStudent success", func(t *testing.T) {
		input := &academicyears.ClassEnrollmentCore{
			SchoolUnitID:   "unit-1",
			AcademicYearID: "ay-1",
			KelasID:        "k-1",
			SiswaID:        "s-1",
			Status:         "aktif",
		}

		mockRepo.On("InsertEnrollment", mock.MatchedBy(func(item *academicyears.ClassEnrollmentCore) bool {
			return item.Status == "Aktif"
		})).Return(nil).Once()

		res, err := svc.EnrollStudent(input)
		assert.NoError(t, err)
		assert.Equal(t, "Aktif", res.Status)
	})

	t.Run("GetStudentEnrollmentHistory success", func(t *testing.T) {
		mockRepo.On("GetStudentEnrollmentHistory", "s-1").Return([]academicyears.ClassEnrollmentCore{
			{
				ID:          "enr-1",
				SiswaID:     "s-1",
				NamaKelas:   "X-A",
				TahunAjaran: "2025/2026",
				Semester:    "Ganjil",
			},
		}, nil).Once()

		res, err := svc.GetStudentEnrollmentHistory("s-1")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "X-A", res[0].NamaKelas)
	})
}
