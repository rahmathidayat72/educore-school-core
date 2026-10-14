package controllers

import (
	"bytes"
	"encoding/json"
	"go_rest_native_sekolah/features/guardians"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGuardianService struct {
	mock.Mock
}

func (m *mockGuardianService) GetAll(search string) ([]guardians.GuardianCore, error) {
	args := m.Called(search)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]guardians.GuardianCore), args.Error(1)
}

func (m *mockGuardianService) GetByID(id string) (*guardians.GuardianCore, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*guardians.GuardianCore), args.Error(1)
}

func (m *mockGuardianService) GetByUserID(userID string) (*guardians.GuardianCore, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*guardians.GuardianCore), args.Error(1)
}

func (m *mockGuardianService) Create(input *guardians.GuardianCore) (*guardians.GuardianCore, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*guardians.GuardianCore), args.Error(1)
}

func (m *mockGuardianService) Update(id string, input *guardians.GuardianCore) (*guardians.GuardianCore, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*guardians.GuardianCore), args.Error(1)
}

func (m *mockGuardianService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockGuardianService) LinkStudent(guardianID string, siswaID string, isPrimary bool) (*guardians.StudentGuardianCore, error) {
	args := m.Called(guardianID, siswaID, isPrimary)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*guardians.StudentGuardianCore), args.Error(1)
}

func (m *mockGuardianService) UnlinkStudent(guardianID string, siswaID string) error {
	args := m.Called(guardianID, siswaID)
	return args.Error(0)
}

func (m *mockGuardianService) GetStudentsByGuardianID(guardianID string) ([]guardians.StudentGuardianCore, error) {
	args := m.Called(guardianID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]guardians.StudentGuardianCore), args.Error(1)
}

func (m *mockGuardianService) GetGuardiansBySiswaID(siswaID string) ([]guardians.StudentGuardianCore, error) {
	args := m.Called(siswaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]guardians.StudentGuardianCore), args.Error(1)
}

func TestGuardianController(t *testing.T) {
	mockSvc := new(mockGuardianService)
	controller := NewGuardianController(mockSvc)

	t.Run("GetAll success", func(t *testing.T) {
		mockSvc.On("GetAll", "").Return([]guardians.GuardianCore{
			{
				ID:      "g-1",
				Nama:    "Ayah Budi",
				Telepon: "0811111",
			},
		}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/guardians", nil)
		rec := httptest.NewRecorder()

		err := controller.GetAll(rec, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Create success", func(t *testing.T) {
		body := GuardianRequest{
			Nama:    "Ayah Budi",
			Telepon: "0811111",
		}
		jsonBytes, _ := json.Marshal(body)

		mockSvc.On("Create", mock.Anything).Return(&guardians.GuardianCore{
			ID:      "g-1",
			Nama:    "Ayah Budi",
			Telepon: "0811111",
		}, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/guardians/tambah", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		err := controller.Create(rec, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("GetByID empty param", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/guardians/detail", nil)
		rec := httptest.NewRecorder()

		err := controller.GetByID(rec, req)
		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("UnlinkStudent missing params", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/guardians/unlink-student", nil)
		rec := httptest.NewRecorder()

		err := controller.UnlinkStudent(rec, req)
		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
