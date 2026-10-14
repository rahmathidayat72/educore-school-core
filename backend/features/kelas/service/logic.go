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

	if err := ResolveFaseDanJurusan(insert); err != nil {
		return err
	}

	return k.kelasData.Insert(insert)
}

func (k *kelasService) SelectAll(schoolUnitID ...string) ([]kelas.KelasCore, error) {
	if k.kelasData == nil {
		return nil, errors.New("kelas service: Nil repository")
	}

	kelass, err := k.kelasData.SelectAll(schoolUnitID...)
	if err != nil {
		return nil, errors.New("kelas service: gagal mengambil data")
	}

	return kelass, nil
}

func (k *kelasService) SelectById(id string, schoolUnitID ...string) (*kelas.KelasCore, error) {
	if k.kelasData == nil {
		return nil, errors.New("kelas service: Nil repository")
	}

	kelasItem, err := k.kelasData.SelectById(id, schoolUnitID...)
	if err != nil {
		return nil, fmt.Errorf("kelas service: gagal mengambil data kelas: %w", err)
	}

	return kelasItem, nil
}

func (k *kelasService) Update(insert *kelas.KelasCore, id string, schoolUnitID ...string) error {
	if k == nil || k.kelasData == nil {
		return errors.New("Nil repository")
	}

	if id == "" {
		return errors.New("Validation error: id is nil")
	}

	exisData, err := k.kelasData.SelectById(id, schoolUnitID...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("Data tidak ditemukan")
		}
		return fmt.Errorf("Id salah atau gagal mengambil data lama: %w", err)
	}

	if insert.Kelas == "" {
		insert.Kelas = exisData.Kelas
	}
	if insert.Tingkat == "" {
		insert.Tingkat = exisData.Tingkat
	}
	if insert.Fase == "" {
		insert.Fase = exisData.Fase
	}
	if insert.Jurusan == "" {
		insert.Jurusan = exisData.Jurusan
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
	if insert.SchoolUnitID == nil && exisData.SchoolUnitID != nil {
		insert.SchoolUnitID = exisData.SchoolUnitID
	}

	if err := ResolveFaseDanJurusan(insert); err != nil {
		return err
	}

	if err := k.kelasData.Update(insert, id, schoolUnitID...); err != nil {
		return err
	}

	return nil
}

func (k *kelasService) DeleteById(id string, schoolUnitID ...string) error {
	if k == nil || k.kelasData == nil {
		return errors.New("Nil repository")
	}

	if id == "" {
		return errors.New("Validation error: id harus diisi")
	}

	if err := k.kelasData.DeleteById(id, schoolUnitID...); err != nil {
		return fmt.Errorf("gagal menghapus data kelas: %w", err)
	}

	return nil
}

func normalizeTingkat(tingkat string) string {
	t := strings.TrimSpace(strings.ToUpper(tingkat))
	switch t {
	case "1", "I":
		return "1"
	case "2", "II":
		return "2"
	case "3", "III":
		return "3"
	case "4", "IV":
		return "4"
	case "5", "V":
		return "5"
	case "6", "VI":
		return "6"
	case "7", "VII":
		return "7"
	case "8", "VIII":
		return "8"
	case "9", "IX":
		return "9"
	case "10", "X":
		return "10"
	case "11", "XI":
		return "11"
	case "12", "XII":
		return "12"
	case "13", "XIII":
		return "13"
	default:
		return t
	}
}

func getExpectedFase(tingkatNorm string) string {
	switch tingkatNorm {
	case "1", "2":
		return "A"
	case "3", "4":
		return "B"
	case "5", "6":
		return "C"
	case "7", "8", "9":
		return "D"
	case "10":
		return "E"
	case "11", "12", "13":
		return "F"
	default:
		return ""
	}
}

func isValidFase(fase string) bool {
	f := strings.ToUpper(strings.TrimSpace(fase))
	return f == "A" || f == "B" || f == "C" || f == "D" || f == "E" || f == "F"
}

func isSDOrSMP(tingkatNorm string) bool {
	switch tingkatNorm {
	case "1", "2", "3", "4", "5", "6", "7", "8", "9":
		return true
	default:
		return false
	}
}

func ResolveFaseDanJurusan(insert *kelas.KelasCore) error {
	tingkatRaw := strings.TrimSpace(insert.Tingkat)
	faseRaw := strings.ToUpper(strings.TrimSpace(insert.Fase))
	jurusanRaw := strings.TrimSpace(insert.Jurusan)

	tingkatNorm := normalizeTingkat(tingkatRaw)
	expectedFase := getExpectedFase(tingkatNorm)

	if expectedFase != "" {
		if faseRaw == "" {
			insert.Fase = expectedFase
		} else if faseRaw != expectedFase {
			return fmt.Errorf("fase '%s' tidak sesuai dengan tingkat '%s' (seharusnya Fase %s)", faseRaw, tingkatRaw, expectedFase)
		} else {
			insert.Fase = faseRaw
		}
	} else if faseRaw != "" {
		if !isValidFase(faseRaw) {
			return fmt.Errorf("fase '%s' tidak valid: harus A, B, C, D, E, atau F", faseRaw)
		}
		insert.Fase = faseRaw
	}

	// Aturan jurusan untuk SMK
	isSMK := false
	if insert.SchoolUnitID != nil && strings.Contains(strings.ToLower(*insert.SchoolUnitID), "smk") {
		isSMK = true
	}
	if strings.Contains(strings.ToLower(insert.Kelas), "smk") {
		isSMK = true
	}

	if isSMK && jurusanRaw == "" {
		return errors.New("rombel SMK wajib mencantumkan jurusan/program keahlian")
	}

	// Untuk SD dan SMP, jurusan tidak berlaku
	if isSDOrSMP(tingkatNorm) && jurusanRaw != "" {
		insert.Jurusan = ""
	}

	return nil
}
