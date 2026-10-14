package service

import (
	"go_rest_native_sekolah/features/guardians"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGuardianRepo struct {
	mock.Mock
}

func (m *mockGuardianRepo) SelectAll(search string) ([]guardians.GuardianCore, error) {
	args := m.Called(search)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]guardians.GuardianCore), args.Error(1)
}

func (m *mockGuardianRepo) SelectByID(id string) (*guardians.GuardianCore, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*guardians.GuardianCore), args.Error(1)
}

func (m *mockGuardianRepo) SelectByUserID(userID string) (*guardians.GuardianCore, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*guardians.GuardianCore), args.Error(1)
}

func (m *mockGuardianRepo) Insert(input *guardians.GuardianCore) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *mockGuardianRepo) Update(id string, input *guardians.GuardianCore) error {
	args := m.Called(id, input)
	return args.Error(0)
}

func (m *mockGuardianRepo) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockGuardianRepo) LinkStudent(input *guardians.StudentGuardianCore) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *mockGuardianRepo) UnlinkStudent(guardianID string, siswaID string) error {
	args := m.Called(guardianID, siswaID)
	return args.Error(0)
}

func (m *mockGuardianRepo) SelectStudentsByGuardianID(guardianID string) ([]guardians.StudentGuardianCore, error) {
	args := m.Called(guardianID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]guardians.StudentGuardianCore), args.Error(1)
}

func (m *mockGuardianRepo) SelectGuardiansBySiswaID(siswaID string) ([]guardians.StudentGuardianCore, error) {
	args := m.Called(siswaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]guardians.StudentGuardianCore), args.Error(1)
}

func TestGuardianService_CRUD(t *testing.T) {
	mockRepo := new(mockGuardianRepo)
	svc := NewGuardianService(mockRepo)

	t.Run("Create validation errors", func(t *testing.T) {
		_, err := svc.Create(nil)
		assert.Error(t, err)

		_, err = svc.Create(&guardians.GuardianCore{})
		assert.ErrorContains(t, err, "nama wali murid")

		_, err = svc.Create(&guardians.GuardianCore{Nama: "Budi"})
		assert.ErrorContains(t, err, "nomor telepon")
	})

	t.Run("Create success with normalized hubungan", func(t *testing.T) {
		input := &guardians.GuardianCore{
			Nama:     "Budi Hartono",
			Telepon:  "08123456789",
			Hubungan: "bapak",
		}

		mockRepo.On("Insert", mock.MatchedBy(func(item *guardians.GuardianCore) bool {
			return item.Nama == "Budi Hartono" && item.Hubungan == "Ayah"
		})).Return(nil).Once()

		res, err := svc.Create(input)
		assert.NoError(t, err)
		assert.Equal(t, "Ayah", res.Hubungan)
	})

	t.Run("GetByID success with students list", func(t *testing.T) {
		mockRepo.On("SelectByID", "g-1").Return(&guardians.GuardianCore{
			ID:      "g-1",
			Nama:    "Siti Rahma",
			Telepon: "08129999",
		}, nil).Once()

		mockRepo.On("SelectStudentsByGuardianID", "g-1").Return([]guardians.StudentGuardianCore{
			{
				ID:        "sg-1",
				SiswaID:   "s-1",
				NamaSiswa: "Ahmad",
			},
			{
				ID:        "sg-2",
				SiswaID:   "s-2",
				NamaSiswa: "Aisyah",
			},
		}, nil).Once()

		res, err := svc.GetByID("g-1")
		assert.NoError(t, err)
		assert.Equal(t, "Siti Rahma", res.Nama)
		assert.Len(t, res.SiswaList, 2)
	})

	t.Run("Delete success", func(t *testing.T) {
		mockRepo.On("Delete", "g-1").Return(nil).Once()
		err := svc.Delete("g-1")
		assert.NoError(t, err)
	})
}

func TestGuardianService_Relations(t *testing.T) {
	mockRepo := new(mockGuardianRepo)
	svc := NewGuardianService(mockRepo)

	t.Run("LinkStudent validation errors", func(t *testing.T) {
		_, err := svc.LinkStudent("", "s-1", false)
		assert.ErrorContains(t, err, "id wali")

		_, err = svc.LinkStudent("g-1", "", false)
		assert.ErrorContains(t, err, "id siswa")
	})

	t.Run("LinkStudent success", func(t *testing.T) {
		mockRepo.On("LinkStudent", mock.MatchedBy(func(item *guardians.StudentGuardianCore) bool {
			return item.GuardianID == "g-1" && item.SiswaID == "s-1" && item.IsPrimary == true
		})).Return(nil).Once()

		res, err := svc.LinkStudent("g-1", "s-1", true)
		assert.NoError(t, err)
		assert.True(t, res.IsPrimary)
	})

	t.Run("UnlinkStudent success", func(t *testing.T) {
		mockRepo.On("UnlinkStudent", "g-1", "s-1").Return(nil).Once()
		err := svc.UnlinkStudent("g-1", "s-1")
		assert.NoError(t, err)
	})

	t.Run("GetGuardiansBySiswaID error on empty", func(t *testing.T) {
		_, err := svc.GetGuardiansBySiswaID("")
		assert.Error(t, err)
	})

	t.Run("GetGuardiansBySiswaID success", func(t *testing.T) {
		mockRepo.On("SelectGuardiansBySiswaID", "s-1").Return([]guardians.StudentGuardianCore{
			{
				ID:       "sg-1",
				NamaWali: "Ayah Budi",
				Hubungan: "Ayah",
			},
		}, nil).Once()

		res, err := svc.GetGuardiansBySiswaID("s-1")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}
