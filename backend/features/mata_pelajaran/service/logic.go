package service

import (
	"errors"
	"fmt"
	matapelajaran "go_rest_native_sekolah/features/mata_pelajaran"
	"strings"

	"github.com/jackc/pgx/v5"
)

type mataPelajaranServiceinterface struct {
	mataPelajaranData matapelajaran.DataMataPelajaranInterface
}

func NewMataPelajaranService(repo matapelajaran.DataMataPelajaranInterface) matapelajaran.ServiceMapelInterface {
	if repo == nil {
		errors.New("Nil repository")
	}
	return &mataPelajaranServiceinterface{mataPelajaranData: repo}
}

func (m *mataPelajaranServiceinterface) InsertMapel(insert *matapelajaran.MataPelajaranCore) error {
	if m.mataPelajaranData == nil {
		return errors.New("Nil repository")
	}

	if insert == nil {
		return errors.New("insert data is nil")
	}

	insert.Nama_Pelajaran = strings.TrimSpace(insert.Nama_Pelajaran)
	if insert.Nama_Pelajaran == "" {
		return errors.New("Validation error: insert mapel is nil")
	}

	if insert.KKM < 0 || insert.KKM > 100 {
		return errors.New("nilai KKM harus berada di rentang 0 sampai 100")
	}
	if insert.KKM == 0 {
		insert.KKM = 75.0
	}

	if insert.JamPerMinggu < 0 {
		return errors.New("jam per minggu tidak boleh negatif")
	}
	if insert.JamPerMinggu == 0 {
		insert.JamPerMinggu = 2
	}

	if strings.TrimSpace(insert.Kurikulum) == "" {
		insert.Kurikulum = "Kurikulum Merdeka"
	}
	if strings.TrimSpace(insert.Status) == "" {
		insert.Status = "Aktif"
	}

	return m.mataPelajaranData.InsertMapel(insert)
}

func (m *mataPelajaranServiceinterface) SelectAllMapel(schoolUnitID ...string) ([]matapelajaran.MataPelajaranCore, error) {
	if m.mataPelajaranData == nil {
		return nil, errors.New("Nil repository")
	}

	mapels, err := m.mataPelajaranData.SelectAllMapel(schoolUnitID...)
	if err != nil {
		return nil, errors.New("MataPelajaranService: gagal mengambil data")
	}

	return mapels, nil
}

func (m *mataPelajaranServiceinterface) SelectMapelById(id string, schoolUnitID ...string) (*matapelajaran.MataPelajaranCore, error) {
	if m.mataPelajaranData == nil {
		return nil, errors.New("Nil repository")
	}

	mapel, err := m.mataPelajaranData.SelectMapelById(id, schoolUnitID...)
	if err != nil {
		return nil, errors.New("MataPelajaranService: gagal mengambil data")
	}

	return mapel, nil
}

func (m *mataPelajaranServiceinterface) UpdateMapel(update *matapelajaran.MataPelajaranCore, id string, schoolUnitID ...string) error {
	if m == nil || m.mataPelajaranData == nil {
		return errors.New("Nil repository")
	}

	if id == "" {
		return errors.New("Validation error: id is nil")
	}

	existingData, err := m.mataPelajaranData.SelectMapelById(id, schoolUnitID...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("mata pelajaran service: Data tidak ditemukan")
		}
		return fmt.Errorf("gagal mengambil data lama: %w", err)
	}

	if update.Nama_Pelajaran == "" {
		update.Nama_Pelajaran = existingData.Nama_Pelajaran
	}
	if update.JamPerMinggu <= 0 && existingData.JamPerMinggu > 0 {
		update.JamPerMinggu = existingData.JamPerMinggu
	}
	if update.KKM <= 0 && existingData.KKM > 0 {
		update.KKM = existingData.KKM
	}
	if update.Kurikulum == "" {
		update.Kurikulum = existingData.Kurikulum
	}
	if update.Status == "" {
		update.Status = existingData.Status
	}

	if err := m.mataPelajaranData.UpdateMapel(update, id, schoolUnitID...); err != nil {
		return fmt.Errorf("gagal update data mata pelajaran: %w", err)
	}

	return nil
}

func (m *mataPelajaranServiceinterface) DeleteMapel(id string, schoolUnitID ...string) error {
	if m == nil || m.mataPelajaranData == nil {
		return errors.New("Nil repository")
	}
	if id == "" {
		return errors.New("Validation error: id is nil")
	}
	if err := m.mataPelajaranData.DeleteMapel(id, schoolUnitID...); err != nil {
		return fmt.Errorf("gagal menghapus data mata pelajaran: %w", err)
	}
	return nil
}
