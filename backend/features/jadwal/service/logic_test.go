package service

import (
	"go_rest_native_sekolah/features/jadwal"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockJadwalRepo struct {
	mock.Mock
}

func (m *mockJadwalRepo) SelectAll(role, userID, hariFilter, kelasIDFilter, guruIDFilter, tahunFilter, semesterFilter, statusFilter string, schoolUnitID ...string) ([]jadwal.JadwalCore, error) {
	args := m.Called(role, userID, hariFilter, kelasIDFilter, guruIDFilter, tahunFilter, semesterFilter, statusFilter, schoolUnitID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]jadwal.JadwalCore), args.Error(1)
}

func (m *mockJadwalRepo) SelectByID(id string, schoolUnitID ...string) (*jadwal.JadwalCore, error) {
	args := m.Called(id, schoolUnitID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jadwal.JadwalCore), args.Error(1)
}

func (m *mockJadwalRepo) Insert(input *jadwal.JadwalCore) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *mockJadwalRepo) Update(id string, input *jadwal.JadwalCore, schoolUnitID ...string) error {
	args := m.Called(id, input, schoolUnitID)
	return args.Error(0)
}

func (m *mockJadwalRepo) Delete(id string, schoolUnitID ...string) error {
	args := m.Called(id, schoolUnitID)
	return args.Error(0)
}

func (m *mockJadwalRepo) CheckConflict(excludeID, mapelID, guruID, kelasID, hari, jamMulai, jamSelesai, ruang, tahunAjaran, semester string, aktif bool, schoolUnitID ...string) (*jadwal.ConflictInfo, error) {
	args := m.Called(excludeID, mapelID, guruID, kelasID, hari, jamMulai, jamSelesai, ruang, tahunAjaran, semester, aktif, schoolUnitID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jadwal.ConflictInfo), args.Error(1)
}

func (m *mockJadwalRepo) GetGuruIDByUserID(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func TestJadwalService_ConflictDetection(t *testing.T) {
	mockRepo := new(mockJadwalRepo)
	svc := NewJadwalService(mockRepo)

	t.Run("Create failed - validation jam selesai lebih awal dari jam mulai", func(t *testing.T) {
		input := &jadwal.JadwalCore{
			MataPelajaranID: "mp-1",
			Hari:            1,
			JamMulai:        "10:00",
			JamSelesai:      "08:00",
		}
		_, err := svc.Create(input, "admin", "user-1")
		assert.ErrorContains(t, err, "jam selesai harus lebih besar dari jam mulai")
	})

	t.Run("Create failed - guru bentrok mengajar", func(t *testing.T) {
		input := &jadwal.JadwalCore{
			MataPelajaranID: "mp-1",
			IDGuru:          "guru-1",
			KelasID:         "kls-1",
			Hari:            1,
			JamMulai:        "08:00",
			JamSelesai:      "09:30",
			Ruang:           "Lab 1",
			TahunAjaran:     "2026/2027",
			Semester:        "Ganjil",
			Aktif:           true,
		}

		mockRepo.On("CheckConflict", "", "mp-1", "guru-1", "kls-1", "1", "08:00:00", "09:30:00", "Lab 1", "2026/2027", "Ganjil", true, mock.Anything).
			Return(&jadwal.ConflictInfo{
				HasConflict: true,
				Message:     "Guru Pak Budi sudah mengajar Matematika di kelas 7A pukul 08:00-09:30.",
			}, nil).Once()

		_, err := svc.Create(input, "admin", "user-1")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "konflik jadwal: Guru Pak Budi sudah mengajar")
	})

	t.Run("Create failed - ruangan bentrok digunakan kelas lain", func(t *testing.T) {
		input := &jadwal.JadwalCore{
			MataPelajaranID: "mp-2",
			IDGuru:          "guru-2",
			KelasID:         "kls-2",
			Hari:            2,
			JamMulai:        "10:00",
			JamSelesai:      "11:30",
			Ruang:           "Lab Komputer",
			TahunAjaran:     "2026/2027",
			Semester:        "Ganjil",
			Aktif:           true,
		}

		mockRepo.On("CheckConflict", "", "mp-2", "guru-2", "kls-2", "2", "10:00:00", "11:30:00", "Lab Komputer", "2026/2027", "Ganjil", true, mock.Anything).
			Return(&jadwal.ConflictInfo{
				HasConflict: true,
				Message:     "Ruangan Lab Komputer sudah digunakan kelas 8B untuk pelajaran Informatika pukul 10:00-11:30.",
			}, nil).Once()

		_, err := svc.Create(input, "admin", "user-1")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "konflik jadwal: Ruangan Lab Komputer sudah digunakan")
	})

	t.Run("Create success - no conflict", func(t *testing.T) {
		input := &jadwal.JadwalCore{
			ID:              "jdw-new-1",
			MataPelajaranID: "mp-3",
			IDGuru:          "guru-3",
			KelasID:         "kls-3",
			Hari:            3,
			JamMulai:        "07:30",
			JamSelesai:      "09:00",
			Ruang:           "Ruang 101",
			TahunAjaran:     "2026/2027",
			Semester:        "Ganjil",
			Aktif:           true,
		}

		mockRepo.On("CheckConflict", "", "mp-3", "guru-3", "kls-3", "3", "07:30:00", "09:00:00", "Ruang 101", "2026/2027", "Ganjil", true, mock.Anything).
			Return(&jadwal.ConflictInfo{HasConflict: false}, nil).Once()

		mockRepo.On("Insert", mock.Anything).Return(nil).Once()
		mockRepo.On("SelectByID", mock.Anything, mock.Anything).Return(&jadwal.JadwalCore{
			ID:            "jdw-new-1",
			NamaPelajaran: "Bahasa Indonesia",
			NamaGuru:      "Ibu Ratna",
			NamaKelas:     "7A",
			Hari:          3,
			JamMulai:      "07:30:00",
			JamSelesai:    "09:00:00",
		}, nil).Once()

		res, err := svc.Create(input, "admin", "user-1")
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "Bahasa Indonesia", res.NamaPelajaran)
	})
}
