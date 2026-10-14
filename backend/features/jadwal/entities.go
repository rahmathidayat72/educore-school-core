package jadwal

import "time"

type (
	JadwalCore struct {
		ID              string     `json:"id"`
		KodeJadwal      string     `json:"kode_jadwal"`
		MataPelajaranID string     `json:"mata_pelajaran_id"`
		KodeMapel       string     `json:"kode_mapel"`
		NamaPelajaran   string     `json:"nama_pelajaran"`
		IDGuru          string     `json:"id_guru"`
		NamaGuru        string     `json:"nama_guru"`
		NIPGuru         string     `json:"nip_guru"`
		KelasID         string     `json:"kelas_id"`
		NamaKelas       string     `json:"nama_kelas"`
		KodeKelas       string     `json:"kode_kelas"`
		TingkatKelas    string     `json:"tingkat_kelas"`
		JurusanKelas    string     `json:"jurusan_kelas"`
		Hari            int        `json:"hari"` // 1: Senin ... 7: Minggu
		NamaHari        string     `json:"nama_hari"`
		JamKeMulai      int        `json:"jam_ke_mulai"`
		JamKeSelesai    int        `json:"jam_ke_selesai"`
		JamMulai        string     `json:"jam_mulai"`   // HH:MM:SS atau HH:MM
		JamSelesai      string     `json:"jam_selesai"` // HH:MM:SS atau HH:MM
		Ruang           string     `json:"ruang"`
		JenisKegiatan   string     `json:"jenis_kegiatan"` // Tatap Muka, Praktikum, dll
		TahunAjaran     string     `json:"tahun_ajaran"`
		Semester        string     `json:"semester"`
		Keterangan      string     `json:"keterangan"`
		Aktif           bool       `json:"aktif"`
		CreatedAt       *time.Time `json:"created_at,omitempty"`
		UpdateAt        *time.Time `json:"update_at,omitempty"`
		DeleteAt        *time.Time `json:"delete_at,omitempty"`
	}

	ConflictInfo struct {
		HasConflict bool   `json:"has_conflict"`
		Message     string `json:"message"`
	}

	DataJadwalInterface interface {
		SelectAll(role, userID, hariFilter, kelasIDFilter, guruIDFilter, tahunFilter, semesterFilter, statusFilter string) ([]JadwalCore, error)
		SelectByID(id string) (*JadwalCore, error)
		Insert(input *JadwalCore) error
		Update(id string, input *JadwalCore) error
		Delete(id string) error
		CheckConflict(excludeID, mapelID, guruID, kelasID, hari, jamMulai, jamSelesai, ruang, tahunAjaran, semester string, aktif bool) (*ConflictInfo, error)
		GetGuruIDByUserID(userID string) (string, error)
	}

	ServiceJadwalInterface interface {
		GetAll(role, userID, hariFilter, kelasIDFilter, guruIDFilter, tahunFilter, semesterFilter, statusFilter string) ([]JadwalCore, error)
		GetByID(id string, role, userID string) (*JadwalCore, error)
		Create(input *JadwalCore, role, userID string) (*JadwalCore, error)
		Update(id string, input *JadwalCore, role, userID string) (*JadwalCore, error)
		Delete(id string, role, userID string) error
	}
)

func NamaHariIndo(h int) string {
	switch h {
	case 1:
		return "Senin"
	case 2:
		return "Selasa"
	case 3:
		return "Rabu"
	case 4:
		return "Kamis"
	case 5:
		return "Jumat"
	case 6:
		return "Sabtu"
	case 7:
		return "Minggu"
	default:
		return "-"
	}
}
