package siswa

import "time"

type (
	// SiswaCore merepresentasikan data lengkap siswa
	SiswaCore struct {
		ID            string     `json:"id"`
		SchoolUnitID  *string    `json:"school_unit_id,omitempty"`
		Nama          string     `json:"nama"`
		Kelas_ID      string     `json:"kelas_id"`
		Nama_Kelas    string     `json:"nama_kelas"`
		NIS           string     `json:"nis"`
		NISN          string     `json:"nisn"`
		JenisKelamin  string     `json:"jenis_kelamin"`
		TempatLahir   string     `json:"tempat_lahir"`
		TanggalLahir  string     `json:"tanggal_lahir"` // format YYYY-MM-DD
		Agama         string     `json:"agama"`
		NoTelepon     string     `json:"no_telepon"`
		NamaWali      string     `json:"nama_wali"`
		NoTeleponWali string     `json:"no_telepon_wali"`
		TahunMasuk    string     `json:"tahun_masuk"`
		Status        string     `json:"status"` // Aktif, Lulus, Pindah, Nonaktif
		Email         string     `json:"email"`
		Alamat        string     `json:"alamat"`
		Update_At     time.Time  `json:"update_at"`
		Delete_At     *time.Time `json:"delete_at"`
	}

	DataSiswaInterface interface {
		SelectAllSiswa(schoolUnitID ...string) ([]SiswaCore, error)
		InsertSiswa(insert *SiswaCore) error
		Update(insert *SiswaCore, id string, schoolUnitID ...string) error
		SelectById(id string, schoolUnitID ...string) (*SiswaCore, error)
		DeleteById(id string, schoolUnitID ...string) error
	}

	ServiceSiswaInterface interface {
		SelectAllSiswa(schoolUnitID ...string) ([]SiswaCore, error)
		InsertSiswa(insert *SiswaCore) error
		Update(insert *SiswaCore, id string, schoolUnitID ...string) error
		SelectById(id string, schoolUnitID ...string) (*SiswaCore, error)
		DeleteById(id string, schoolUnitID ...string) error
	}
)
