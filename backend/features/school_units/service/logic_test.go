package service

import (
	"errors"
	schoolunits "go_rest_native_sekolah/features/school_units"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSchoolUnitRepo struct {
	mock.Mock
}

func (m *mockSchoolUnitRepo) SelectAll(jenjang, status string) ([]schoolunits.SchoolUnitCore, error) {
	args := m.Called(jenjang, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]schoolunits.SchoolUnitCore), args.Error(1)
}

func (m *mockSchoolUnitRepo) SelectByID(id string) (*schoolunits.SchoolUnitCore, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schoolunits.SchoolUnitCore), args.Error(1)
}

func (m *mockSchoolUnitRepo) Insert(input *schoolunits.SchoolUnitCore) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *mockSchoolUnitRepo) Update(id string, input *schoolunits.SchoolUnitCore) error {
	args := m.Called(id, input)
	return args.Error(0)
}

func (m *mockSchoolUnitRepo) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockSchoolUnitRepo) CheckNPSNExists(npsn string, excludeID string) (bool, error) {
	args := m.Called(npsn, excludeID)
	return args.Bool(0), args.Error(1)
}

func TestGetAll(t *testing.T) {
	mockRepo := new(mockSchoolUnitRepo)
	service := NewSchoolUnitService(mockRepo)

	expected := []schoolunits.SchoolUnitCore{
		{
			ID:          "unit-sma-001",
			NamaSekolah: "SMA EduCore 1",
			Jenjang:     "SMA",
			NPSN:        "20109901",
			Status:      "aktif",
		},
	}

	t.Run("success get all", func(t *testing.T) {
		mockRepo.On("SelectAll", "", "").Return(expected, nil).Once()

		res, err := service.GetAll("", "")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(res))
		assert.Equal(t, "unit-sma-001", res[0].ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error get all", func(t *testing.T) {
		mockRepo.On("SelectAll", "SMA", "aktif").Return(nil, errors.New("db error")).Once()

		res, err := service.GetAll("SMA", "aktif")
		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	mockRepo := new(mockSchoolUnitRepo)
	service := NewSchoolUnitService(mockRepo)

	expected := &schoolunits.SchoolUnitCore{
		ID:          "unit-sma-001",
		NamaSekolah: "SMA EduCore 1",
		Jenjang:     "SMA",
		NPSN:        "20109901",
		Status:      "aktif",
	}

	t.Run("success get by id", func(t *testing.T) {
		mockRepo.On("SelectByID", "unit-sma-001").Return(expected, nil).Once()

		res, err := service.GetByID("unit-sma-001")
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "SMA EduCore 1", res.NamaSekolah)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty id error", func(t *testing.T) {
		res, err := service.GetByID("")
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestCreate(t *testing.T) {
	mockRepo := new(mockSchoolUnitRepo)
	service := NewSchoolUnitService(mockRepo)

	t.Run("success create", func(t *testing.T) {
		input := &schoolunits.SchoolUnitCore{
			NamaSekolah: "SMP EduCore 2",
			Jenjang:     "smp",
			NPSN:        "20109920",
			Status:      "aktif",
		}
		expectedResult := &schoolunits.SchoolUnitCore{
			ID:          "unit-smp-002",
			NamaSekolah: "SMP EduCore 2",
			Jenjang:     "SMP",
			NPSN:        "20109920",
			Status:      "aktif",
		}

		mockRepo.On("CheckNPSNExists", "20109920", "").Return(false, nil).Once()
		mockRepo.On("Insert", mock.AnythingOfType("*schoolunits.SchoolUnitCore")).Return(nil).Once()
		mockRepo.On("SelectByID", mock.Anything).Return(expectedResult, nil).Once()

		res, err := service.Create(input)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "SMP EduCore 2", res.NamaSekolah)
		assert.Equal(t, "SMP", res.Jenjang)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error - empty nama sekolah", func(t *testing.T) {
		input := &schoolunits.SchoolUnitCore{
			NamaSekolah: "",
			Jenjang:     "SMP",
		}

		res, err := service.Create(input)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "nama sekolah wajib diisi")
	})

	t.Run("validation error - invalid jenjang", func(t *testing.T) {
		input := &schoolunits.SchoolUnitCore{
			NamaSekolah: "Universitas EduCore",
			Jenjang:     "UNIVERSITAS",
		}

		res, err := service.Create(input)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "jenjang tidak valid")
	})

	t.Run("validation error - duplicate npsn", func(t *testing.T) {
		input := &schoolunits.SchoolUnitCore{
			NamaSekolah: "SMA EduCore 3",
			Jenjang:     "SMA",
			NPSN:        "20109901",
		}

		mockRepo.On("CheckNPSNExists", "20109901", "").Return(true, nil).Once()

		res, err := service.Create(input)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "NPSN sudah digunakan")
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockRepo := new(mockSchoolUnitRepo)
	service := NewSchoolUnitService(mockRepo)

	existing := &schoolunits.SchoolUnitCore{
		ID:          "unit-sma-001",
		NamaSekolah: "SMA EduCore 1",
		Jenjang:     "SMA",
		NPSN:        "20109901",
		Status:      "aktif",
	}

	t.Run("success update", func(t *testing.T) {
		updateInput := &schoolunits.SchoolUnitCore{
			NamaSekolah: "SMA EduCore Unggulan 1",
		}

		mockRepo.On("SelectByID", "unit-sma-001").Return(existing, nil).Once()
		mockRepo.On("CheckNPSNExists", "20109901", "unit-sma-001").Return(false, nil).Once()
		mockRepo.On("Update", "unit-sma-001", mock.AnythingOfType("*schoolunits.SchoolUnitCore")).Return(nil).Once()
		mockRepo.On("SelectByID", "unit-sma-001").Return(&schoolunits.SchoolUnitCore{
			ID:          "unit-sma-001",
			NamaSekolah: "SMA EduCore Unggulan 1",
			Jenjang:     "SMA",
			NPSN:        "20109901",
			Status:      "aktif",
		}, nil).Once()

		res, err := service.Update("unit-sma-001", updateInput)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "SMA EduCore Unggulan 1", res.NamaSekolah)
		mockRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockRepo := new(mockSchoolUnitRepo)
	service := NewSchoolUnitService(mockRepo)

	existing := &schoolunits.SchoolUnitCore{
		ID:          "unit-sma-001",
		NamaSekolah: "SMA EduCore 1",
		Jenjang:     "SMA",
	}

	t.Run("success delete", func(t *testing.T) {
		mockRepo.On("SelectByID", "unit-sma-001").Return(existing, nil).Once()
		mockRepo.On("Delete", "unit-sma-001").Return(nil).Once()

		err := service.Delete("unit-sma-001")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("delete non existing unit", func(t *testing.T) {
		mockRepo.On("SelectByID", "non-existent").Return(nil, errors.New("unit sekolah tidak ditemukan")).Once()

		err := service.Delete("non-existent")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
