package guru

import "time"

type (
	// GuruCore struct untuk merepresentasikan data lengkap guru
	GuruCore struct {
		ID                 string     `json:"id"`
		ID_User            string     `json:"id_user"`
		Nama               string     `json:"nama"`
		NIP                string     `json:"nip"`
		NUPTK              string     `json:"nuptk"`
		JenisKelamin       string     `json:"jenis_kelamin"`
		TempatLahir        string     `json:"tempat_lahir"`
		TanggalLahir       string     `json:"tanggal_lahir"` // format YYYY-MM-DD
		NoTelepon          string     `json:"no_telepon"`
		Agama              string     `json:"agama"`
		PendidikanTerakhir string     `json:"pendidikan_terakhir"`
		StatusKepegawaian  string     `json:"status_kepegawaian"`
		Email              string     `json:"email"`
		Alamat             string     `json:"alamat"`
		Update_At          time.Time  `json:"update_at"`
		Delete_At          *time.Time `json:"delete_at"`
	}

	DataGuruInterface interface {
		SelectAllGuru() ([]GuruCore, error)
		InsertGuru(insert *GuruCore) error
		Update(insert *GuruCore, id string) error
		SelectById(id string) (*GuruCore, error)
		DeleteById(id string) error
	}

	ServiceGuruInterface interface {
		GetAllGuru() ([]GuruCore, error)
		InsertGuru(insert *GuruCore) error
		UpdateGuru(insert *GuruCore, id string) error
		SelectById(id string) (*GuruCore, error)
		DeleteById(id string) error
	}
)
