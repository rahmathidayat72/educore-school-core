package controllers

import (
	"bytes"
	"encoding/json"
	academicyears "go_rest_native_sekolah/features/academic_years"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAcademicYearService struct {
	mock.Mock
}

func (m *mockAcademicYearService) GetAllYears(schoolUnitID string) ([]academicyears.AcademicYearCore, error) {
	args := m.Called(schoolUnitID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]academicyears.AcademicYearCore), args.Error(1)
}

func (m *mockAcademicYearService) GetYearByID(id string, schoolUnitID ...string) (*academicyears.AcademicYearCore, error) {
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

func (m *mockAcademicYearService) GetActiveYear(schoolUnitID string) (*academicyears.AcademicYearCore, error) {
	args := m.Called(schoolUnitID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*academicyears.AcademicYearCore), args.Error(1)
}

func (m *mockAcademicYearService) CreateYear(input *academicyears.AcademicYearCore) (*academicyears.AcademicYearCore, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*academicyears.AcademicYearCore), args.Error(1)
}

func (m *mockAcademicYearService) UpdateYear(id string, input *academicyears.AcademicYearCore, schoolUnitID ...string) (*academicyears.AcademicYearCore, error) {
	var unit string
	if len(schoolUnitID) > 0 {
		unit = schoolUnitID[0]
	}
	args := m.Called(id, input, unit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*academicyears.AcademicYearCore), args.Error(1)
}

func (m *mockAcademicYearService) SetActiveYear(id string, schoolUnitID string) error {
	args := m.Called(id, schoolUnitID)
	return args.Error(0)
}

func (m *mockAcademicYearService) DeleteYear(id string, schoolUnitID ...string) error {
	var unit string
	if len(schoolUnitID) > 0 {
		unit = schoolUnitID[0]
	}
	args := m.Called(id, unit)
	return args.Error(0)
}

func (m *mockAcademicYearService) GetEnrollments(schoolUnitID string, academicYearID string, kelasID string) ([]academicyears.ClassEnrollmentCore, error) {
	args := m.Called(schoolUnitID, academicYearID, kelasID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]academicyears.ClassEnrollmentCore), args.Error(1)
}

func (m *mockAcademicYearService) GetEnrollmentByID(id string, schoolUnitID ...string) (*academicyears.ClassEnrollmentCore, error) {
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

func (m *mockAcademicYearService) EnrollStudent(input *academicyears.ClassEnrollmentCore) (*academicyears.ClassEnrollmentCore, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*academicyears.ClassEnrollmentCore), args.Error(1)
}

func (m *mockAcademicYearService) UpdateEnrollment(id string, input *academicyears.ClassEnrollmentCore, schoolUnitID ...string) (*academicyears.ClassEnrollmentCore, error) {
	var unit string
	if len(schoolUnitID) > 0 {
		unit = schoolUnitID[0]
	}
	args := m.Called(id, input, unit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*academicyears.ClassEnrollmentCore), args.Error(1)
}

func (m *mockAcademicYearService) DeleteEnrollment(id string, schoolUnitID ...string) error {
	var unit string
	if len(schoolUnitID) > 0 {
		unit = schoolUnitID[0]
	}
	args := m.Called(id, unit)
	return args.Error(0)
}

func (m *mockAcademicYearService) GetStudentEnrollmentHistory(siswaID string) ([]academicyears.ClassEnrollmentCore, error) {
	args := m.Called(siswaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]academicyears.ClassEnrollmentCore), args.Error(1)
}

func TestAcademicYearController(t *testing.T) {
	mockSvc := new(mockAcademicYearService)
	controller := NewAcademicYearController(mockSvc)

	t.Run("GetAllYears success", func(t *testing.T) {
		mockSvc.On("GetAllYears", "").Return([]academicyears.AcademicYearCore{
			{
				ID:          "ay-1",
				TahunAjaran: "2026/2027",
				Semester:    "Ganjil",
				IsActive:    true,
			},
		}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/academic-years", nil)
		rec := httptest.NewRecorder()

		err := controller.GetAllYears(rec, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("CreateYear validation empty unit", func(t *testing.T) {
		body := AcademicYearRequest{
			TahunAjaran: "2026/2027",
			Semester:    "Ganjil",
		}
		jsonBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/academic-years/tambah", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		err := controller.CreateYear(rec, req)
		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("GetStudentEnrollmentHistory missing param", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/class-enrollments/student-history", nil)
		rec := httptest.NewRecorder()

		err := controller.GetStudentEnrollmentHistory(rec, req)
		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
