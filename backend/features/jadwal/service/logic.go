package service

import (
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/jadwal"
	"strconv"
	"strings"
	"time"
)

type JadwalService struct {
	repo jadwal.DataJadwalInterface
}

func NewJadwalService(repo jadwal.DataJadwalInterface) jadwal.ServiceJadwalInterface {
	if repo == nil {
		panic("jadwal service: nil repo")
	}
	return &JadwalService{repo: repo}
}

func (s *JadwalService) GetAll(role, userID, hariFilter, kelasIDFilter, guruIDFilter, tahunFilter, semesterFilter, statusFilter string) ([]jadwal.JadwalCore, error) {
	return s.repo.SelectAll(role, userID, hariFilter, kelasIDFilter, guruIDFilter, tahunFilter, semesterFilter, statusFilter)
}

func (s *JadwalService) GetByID(id string, role, userID string) (*jadwal.JadwalCore, error) {
	data, err := s.repo.SelectByID(id)
	if err != nil {
		return nil, err
	}

	// Jika role guru, pastikan guru hanya bisa akses jadwal miliknya
	if strings.ToLower(role) == "guru" {
		guruID, err := s.repo.GetGuruIDByUserID(userID)
		if err != nil || guruID != data.IDGuru {
			return nil, errors.New("akses ditolak: Anda tidak memiliki hak akses ke jadwal ini")
		}
	}

	return data, nil
}

func (s *JadwalService) validateInput(input *jadwal.JadwalCore) error {
	if strings.TrimSpace(input.MataPelajaranID) == "" {
		return errors.New("mata pelajaran wajib dipilih")
	}
	if input.Hari < 1 || input.Hari > 7 {
		return errors.New("hari tidak valid (pilih rentang 1-7)")
	}

	if input.JamKeMulai > 0 && input.JamKeSelesai > 0 {
		if input.JamKeSelesai < input.JamKeMulai {
			return errors.New("nomor jam ke selesai tidak boleh lebih awal dari jam ke mulai")
		}
	}

	jamMulai := strings.TrimSpace(input.JamMulai)
	jamSelesai := strings.TrimSpace(input.JamSelesai)
	if jamMulai == "" || jamSelesai == "" {
		return errors.New("jam mulai dan jam selesai wajib diisi")
	}

	// Normalisasi format time jika hanya HH:MM
	if len(jamMulai) == 5 {
		jamMulai += ":00"
	}
	if len(jamSelesai) == 5 {
		jamSelesai += ":00"
	}

	tMulai, err := time.Parse("15:04:05", jamMulai)
	if err != nil {
		return fmt.Errorf("format jam mulai tidak valid (gunakan HH:MM): %w", err)
	}
	tSelesai, err := time.Parse("15:04:05", jamSelesai)
	if err != nil {
		return fmt.Errorf("format jam selesai tidak valid (gunakan HH:MM): %w", err)
	}

	if !tSelesai.After(tMulai) {
		return errors.New("jam selesai harus lebih besar dari jam mulai")
	}

	input.JamMulai = jamMulai
	input.JamSelesai = jamSelesai

	if strings.TrimSpace(input.TahunAjaran) == "" {
		input.TahunAjaran = "2026/2027"
	}
	if strings.TrimSpace(input.Semester) == "" {
		input.Semester = "Ganjil"
	}
	if strings.TrimSpace(input.JenisKegiatan) == "" {
		input.JenisKegiatan = "Tatap Muka"
	}

	return nil
}

func (s *JadwalService) Create(input *jadwal.JadwalCore, role, userID string) (*jadwal.JadwalCore, error) {
	if err := s.validateInput(input); err != nil {
		return nil, err
	}

	// Jika role guru, pastikan IDGuru diisi dengan guruID milik user login
	if strings.ToLower(role) == "guru" {
		guruID, err := s.repo.GetGuruIDByUserID(userID)
		if err != nil {
			return nil, errors.New("akun Anda belum terhubung dengan data profil guru")
		}
		if input.IDGuru != "" && input.IDGuru != guruID {
			return nil, errors.New("akses ditolak: guru hanya dapat membuat jadwal atas namanya sendiri")
		}
		input.IDGuru = guruID
	}

	// Cek bentrok jadwal
	conflict, err := s.repo.CheckConflict("", input.MataPelajaranID, input.IDGuru, input.KelasID,
		strconv.Itoa(input.Hari), input.JamMulai, input.JamSelesai, input.Ruang, input.TahunAjaran, input.Semester, input.Aktif)
	if err != nil {
		return nil, err
	}
	if conflict.HasConflict {
		return nil, fmt.Errorf("konflik jadwal: %s", conflict.Message)
	}

	if err := s.repo.Insert(input); err != nil {
		return nil, err
	}

	return s.repo.SelectByID(input.ID)
}

func (s *JadwalService) Update(id string, input *jadwal.JadwalCore, role, userID string) (*jadwal.JadwalCore, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("parameter id wajib diisi")
	}

	// Ambil data eksisting untuk merge dan cek otorisasi guru
	existing, err := s.repo.SelectByID(id)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(role) == "guru" {
		guruID, err := s.repo.GetGuruIDByUserID(userID)
		if err != nil || guruID != existing.IDGuru {
			return nil, errors.New("akses ditolak: Anda hanya dapat memperbarui jadwal mengajar Anda sendiri")
		}
		input.IDGuru = guruID
	}

	// Partial fallback nilai lama jika field kosong
	if strings.TrimSpace(input.MataPelajaranID) == "" {
		input.MataPelajaranID = existing.MataPelajaranID
	}
	if input.Hari <= 0 {
		input.Hari = existing.Hari
	}
	if strings.TrimSpace(input.JamMulai) == "" {
		input.JamMulai = existing.JamMulai
	}
	if strings.TrimSpace(input.JamSelesai) == "" {
		input.JamSelesai = existing.JamSelesai
	}
	if strings.TrimSpace(input.IDGuru) == "" {
		input.IDGuru = existing.IDGuru
	}
	if strings.TrimSpace(input.KelasID) == "" {
		input.KelasID = existing.KelasID
	}
	if strings.TrimSpace(input.TahunAjaran) == "" {
		input.TahunAjaran = existing.TahunAjaran
	}
	if strings.TrimSpace(input.Semester) == "" {
		input.Semester = existing.Semester
	}
	if strings.TrimSpace(input.JenisKegiatan) == "" {
		input.JenisKegiatan = existing.JenisKegiatan
	}

	if err := s.validateInput(input); err != nil {
		return nil, err
	}

	// Cek bentrok jadwal (exclude ID saat ini)
	conflict, err := s.repo.CheckConflict(id, input.MataPelajaranID, input.IDGuru, input.KelasID,
		strconv.Itoa(input.Hari), input.JamMulai, input.JamSelesai, input.Ruang, input.TahunAjaran, input.Semester, input.Aktif)
	if err != nil {
		return nil, err
	}
	if conflict.HasConflict {
		return nil, fmt.Errorf("konflik jadwal: %s", conflict.Message)
	}

	if err := s.repo.Update(id, input); err != nil {
		return nil, err
	}

	return s.repo.SelectByID(id)
}

func (s *JadwalService) Delete(id string, role, userID string) error {
	existing, err := s.repo.SelectByID(id)
	if err != nil {
		return err
	}

	if strings.ToLower(role) == "guru" {
		guruID, err := s.repo.GetGuruIDByUserID(userID)
		if err != nil || guruID != existing.IDGuru {
			return errors.New("akses ditolak: Anda hanya dapat menghapus jadwal mengajar Anda sendiri")
		}
	}

	return s.repo.Delete(id)
}
