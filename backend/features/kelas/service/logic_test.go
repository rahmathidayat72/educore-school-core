package service

import (
	"errors"
	"go_rest_native_sekolah/features/kelas"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock untuk DataKelasInterface
type mockDataKelas struct {
	mock.Mock
}

func (m *mockDataKelas) SelectAll(schoolUnitID ...string) ([]kelas.KelasCore, error) {
	var args mock.Arguments
	if len(schoolUnitID) > 0 && schoolUnitID[0] != "" {
		args = m.Called(schoolUnitID[0])
	} else {
		args = m.Called()
	}
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]kelas.KelasCore), args.Error(1)
}

func (m *mockDataKelas) SelectById(id string, schoolUnitID ...string) (*kelas.KelasCore, error) {
	var args mock.Arguments
	if len(schoolUnitID) > 0 && schoolUnitID[0] != "" {
		args = m.Called(id, schoolUnitID[0])
	} else {
		args = m.Called(id)
	}
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*kelas.KelasCore), args.Error(1)
}

func (m *mockDataKelas) Insert(insert *kelas.KelasCore) error {
	args := m.Called(insert)
	return args.Error(0)
}

func (m *mockDataKelas) Update(insert *kelas.KelasCore, id string, schoolUnitID ...string) error {
	var args mock.Arguments
	if len(schoolUnitID) > 0 && schoolUnitID[0] != "" {
		args = m.Called(insert, id, schoolUnitID[0])
	} else {
		args = m.Called(insert, id)
	}
	return args.Error(0)
}

func (m *mockDataKelas) DeleteById(id string, schoolUnitID ...string) error {
	var args mock.Arguments
	if len(schoolUnitID) > 0 && schoolUnitID[0] != "" {
		args = m.Called(id, schoolUnitID[0])
	} else {
		args = m.Called(id)
	}
	return args.Error(0)
}

// Test SelectAll
func TestSelectAllKelas(t *testing.T) {
	mockRepo := new(mockDataKelas)

	t.Run("success get all kelas", func(t *testing.T) {
		expectedKelas := []kelas.KelasCore{
			{
				ID:        "kelas-001",
				Kelas:     "10A",
				ID_Guru:   "guru-001",
				Nama_Guru: "Budi Santoso",
				Update_At: "2024-06-21",
			},
			{
				ID:        "kelas-002",
				Kelas:     "10B",
				ID_Guru:   "guru-002",
				Nama_Guru: "Siti Nurhaliza",
				Update_At: "2024-06-21",
			},
		}

		mockRepo.On("SelectAll").Return(expectedKelas, nil).Once()

		svc := &kelasService{kelasData: mockRepo}
		result, err := svc.SelectAll()

		assert.NoError(t, err)
		assert.Equal(t, expectedKelas, result)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed get all kelas - repository error", func(t *testing.T) {
		mockRepo.On("SelectAll").Return(nil, errors.New("database error")).Once()

		svc := &kelasService{kelasData: mockRepo}
		result, err := svc.SelectAll()

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed - nil repository", func(t *testing.T) {
		svc := &kelasService{kelasData: nil}
		result, err := svc.SelectAll()

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// Test SelectById
func TestSelectKelasByID(t *testing.T) {
	mockRepo := new(mockDataKelas)

	t.Run("success get kelas by id", func(t *testing.T) {
		expectedKelas := &kelas.KelasCore{
			ID:        "kelas-001",
			Kelas:     "10A",
			ID_Guru:   "guru-001",
			Nama_Guru: "Budi Santoso",
			Update_At: "2024-06-21",
		}

		mockRepo.On("SelectById", "kelas-001").Return(expectedKelas, nil).Once()

		svc := &kelasService{kelasData: mockRepo}
		result, err := svc.SelectById("kelas-001")

		assert.NoError(t, err)
		assert.Equal(t, expectedKelas, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed get kelas by id - not found", func(t *testing.T) {
		mockRepo.On("SelectById", "999").Return(nil, pgx.ErrNoRows).Once()

		svc := &kelasService{kelasData: mockRepo}
		result, err := svc.SelectById("999")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

// Test Insert
func TestInsertKelas(t *testing.T) {
	mockRepo := new(mockDataKelas)

	t.Run("success insert kelas", func(t *testing.T) {
		newKelas := &kelas.KelasCore{
			ID:      "kelas-001",
			Kelas:   "10A",
			ID_Guru: "guru-001",
		}

		mockRepo.On("Insert", newKelas).Return(nil).Once()

		svc := &kelasService{kelasData: mockRepo}
		err := svc.Insert(newKelas)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed insert kelas - nil repository", func(t *testing.T) {
		newKelas := &kelas.KelasCore{
			Kelas:   "10A",
			ID_Guru: "guru-001",
		}

		svc := &kelasService{kelasData: nil}
		err := svc.Insert(newKelas)

		assert.Error(t, err)
	})

	t.Run("failed insert kelas - repository error", func(t *testing.T) {
		newKelas := &kelas.KelasCore{
			Kelas:   "10A",
			ID_Guru: "guru-001",
		}

		mockRepo.On("Insert", newKelas).Return(errors.New("insert failed")).Once()

		svc := &kelasService{kelasData: mockRepo}
		err := svc.Insert(newKelas)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// Test Update
func TestUpdateKelas(t *testing.T) {
	mockRepo := new(mockDataKelas)

	t.Run("success update kelas", func(t *testing.T) {
		existingKelas := &kelas.KelasCore{
			ID:      "kelas-001",
			Kelas:   "10A",
			ID_Guru: "guru-001",
		}

		updatedKelas := &kelas.KelasCore{
			Kelas:   "10A-Updated",
			ID_Guru: "guru-002",
		}

		mockRepo.On("SelectById", "kelas-001").Return(existingKelas, nil).Once()
		mockRepo.On("Update", updatedKelas, "kelas-001").Return(nil).Once()

		svc := &kelasService{kelasData: mockRepo}
		err := svc.Update(updatedKelas, "kelas-001")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed update kelas - not found", func(t *testing.T) {
		mockRepo.On("SelectById", "999").Return(nil, pgx.ErrNoRows).Once()

		svc := &kelasService{kelasData: mockRepo}
		err := svc.Update(&kelas.KelasCore{}, "999")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// Test Delete
func TestDeleteKelas(t *testing.T) {
	mockRepo := new(mockDataKelas)

	t.Run("success delete kelas", func(t *testing.T) {
		mockRepo.On("DeleteById", "kelas-001").Return(nil).Once()

		svc := &kelasService{kelasData: mockRepo}
		err := svc.DeleteById("kelas-001")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed delete kelas - not found", func(t *testing.T) {
		mockRepo.On("DeleteById", "999").Return(errors.New("data not found")).Once()

		svc := &kelasService{kelasData: mockRepo}
		err := svc.DeleteById("999")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// Test Kurikulum Merdeka Multi-Jenjang Validation & Auto-assignment
func TestResolveFaseDanJurusan(t *testing.T) {
	t.Run("success auto-assign fase SD", func(t *testing.T) {
		k := &kelas.KelasCore{
			Kelas:   "Kelas 1A",
			Tingkat: "1",
		}
		err := ResolveFaseDanJurusan(k)
		assert.NoError(t, err)
		assert.Equal(t, "A", k.Fase)
	})

	t.Run("success auto-assign fase SMP", func(t *testing.T) {
		k := &kelas.KelasCore{
			Kelas:   "Kelas 7B",
			Tingkat: "7",
		}
		err := ResolveFaseDanJurusan(k)
		assert.NoError(t, err)
		assert.Equal(t, "D", k.Fase)
	})

	t.Run("success auto-assign fase SMA", func(t *testing.T) {
		k := &kelas.KelasCore{
			Kelas:   "Kelas 10-1",
			Tingkat: "10",
		}
		err := ResolveFaseDanJurusan(k)
		assert.NoError(t, err)
		assert.Equal(t, "E", k.Fase)
	})

	t.Run("failed invalid fase conflict", func(t *testing.T) {
		k := &kelas.KelasCore{
			Kelas:   "Kelas 1A",
			Tingkat: "1",
			Fase:    "F",
		}
		err := ResolveFaseDanJurusan(k)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tidak sesuai dengan tingkat")
	})

	t.Run("failed SMK without jurusan", func(t *testing.T) {
		smkUnitID := "unit-smk-001"
		k := &kelas.KelasCore{
			SchoolUnitID: &smkUnitID,
			Kelas:        "Kelas X TKJ",
			Tingkat:      "10",
			Jurusan:      "",
		}
		err := ResolveFaseDanJurusan(k)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "wajib mencantumkan jurusan")
	})

	t.Run("success SMK with jurusan", func(t *testing.T) {
		smkUnitID := "unit-smk-001"
		k := &kelas.KelasCore{
			SchoolUnitID: &smkUnitID,
			Kelas:        "Kelas X TKJ 1",
			Tingkat:      "10",
			Jurusan:      "Teknik Komputer dan Jaringan",
		}
		err := ResolveFaseDanJurusan(k)
		assert.NoError(t, err)
		assert.Equal(t, "E", k.Fase)
		assert.Equal(t, "Teknik Komputer dan Jaringan", k.Jurusan)
	})

	t.Run("success clear jurusan on SD/SMP", func(t *testing.T) {
		k := &kelas.KelasCore{
			Kelas:   "Kelas 5A",
			Tingkat: "5",
			Jurusan: "IPA",
		}
		err := ResolveFaseDanJurusan(k)
		assert.NoError(t, err)
		assert.Equal(t, "C", k.Fase)
		assert.Empty(t, k.Jurusan)
	})
}
