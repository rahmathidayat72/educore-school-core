package service

import (
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/kelas"
	"strings"

	"github.com/jackc/pgx/v5"
)

type kelasService struct {
	kelasData kelas.DataKelasInterface
}

func NewServiceKelas(repo kelas.DataKelasInterface) kelas.ServiceKelasInterface {
	if repo == nil {
		errors.New("Nil repository")
	}
	return &kelasService{kelasData: repo}
}

func (k *kelasService) Insert(insert *kelas.KelasCore) error {
	if k.kelasData == nil {
		return errors.New("kelas service: Nil repository")
	}

	insert.Kelas = strings.TrimSpace(insert.Kelas)
	if insert.Kelas == "" {
		return errors.New("Validation error: nama kelas tidak boleh kosong")
	}

	if insert.Kapasitas < 0 {
		return errors.New("kapasitas kelas tidak boleh bernilai negatif")
	}
	if insert.Kapasitas == 0 {
		insert.Kapasitas = 36
	}

	if strings.TrimSpace(insert.Status) == "" {
		insert.Status = "Aktif"
	}

	return k.kelasData.Insert(insert)
}

func (k *kelasService) SelectAll() ([]kelas.KelasCore, error) {
	if k.kelasData == nil {
		return nil, errors.New("kelas service: Nil repository")
	}

	kelass, err := k.kelasData.SelectAll()
	if err != nil {
		return nil, errors.New("kelas service: gagal mengambil data")
	}

	return kelass, nil
}

func (k *kelasService) SelectById(id string) (*kelas.KelasCore, error) {
	if k.kelasData == nil {
		return nil, errors.New("kelas service: Nil repository")
	}

	kelasItem, err := k.kelasData.SelectById(id)
	if err != nil {
		return nil, fmt.Errorf("kelas service: gagal mengambil data kelas: %w", err)
	}

	return kelasItem, nil
}

func (k *kelasService) Update(insert *kelas.KelasCore, id string) error {
	if k == nil || k.kelasData == nil {
		return errors.New("Nil repository")
	}

	if id == "" {
		return errors.New("Validation error: id is nil")
	}

	exisData, err := k.kelasData.SelectById(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("Data tidak ditemukan")
		}
		return fmt.Errorf("Id salah atau gagal mengambil data lama: %w", err)
	}

	if insert.Kelas == "" {
		insert.Kelas = exisData.Kelas
	}
	if insert.Kapasitas <= 0 && exisData.Kapasitas > 0 {
		insert.Kapasitas = exisData.Kapasitas
	}
	if insert.TahunAjaran == "" {
		insert.TahunAjaran = exisData.TahunAjaran
	}
	if insert.Status == "" {
		insert.Status = exisData.Status
	}

	if err := k.kelasData.Update(insert, id); err != nil {
		return err
	}

	return nil
}

func (k *kelasService) DeleteById(id string) error {
	if k == nil || k.kelasData == nil {
		return errors.New("Nil repository")
	}

	if id == "" {
		return errors.New("Validation error: id harus diisi")
	}

	if err := k.kelasData.DeleteById(id); err != nil {
		return fmt.Errorf("gagal menghapus data kelas: %w", err)
	}

	return nil
}
