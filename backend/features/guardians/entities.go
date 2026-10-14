package guardians

import "time"

// GuardianCore merepresentasikan data orang tua / wali murid
type GuardianCore struct {
	ID          string                `json:"id"`
	UserID      *string               `json:"user_id,omitempty"`
	Nama        string                `json:"nama"`
	NIK         string                `json:"nik,omitempty"`
	Hubungan    string                `json:"hubungan"` // 'Ayah', 'Ibu', 'Wali', 'Kakak'
	Pekerjaan   string                `json:"pekerjaan,omitempty"`
	Penghasilan string                `json:"penghasilan,omitempty"`
	Telepon     string                `json:"telepon"`
	Email       string                `json:"email,omitempty"`
	Alamat      string                `json:"alamat,omitempty"`
	CreatedAt   *time.Time            `json:"created_at,omitempty"`
	UpdateAt    *time.Time            `json:"update_at,omitempty"`
	DeleteAt    *time.Time            `json:"delete_at,omitempty"`
	SiswaList   []StudentGuardianCore `json:"siswa_list,omitempty"`
}

// StudentGuardianCore merepresentasikan relasi many-to-many siswa dan wali murid
type StudentGuardianCore struct {
	ID          string     `json:"id"`
	GuardianID  string     `json:"guardian_id"`
	SiswaID     string     `json:"siswa_id"`
	IsPrimary   bool       `json:"is_primary"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdateAt    *time.Time `json:"update_at,omitempty"`
	DeleteAt    *time.Time `json:"delete_at,omitempty"`

	// Relational data
	NamaSiswa   string     `json:"nama_siswa,omitempty"`
	NISN        string     `json:"nisn,omitempty"`
	NamaKelas   string     `json:"nama_kelas,omitempty"`
	NamaSekolah string     `json:"nama_sekolah,omitempty"`
	NamaWali    string     `json:"nama_wali,omitempty"`
	Hubungan    string     `json:"hubungan,omitempty"`
	TeleponWali string     `json:"telepon_wali,omitempty"`
}

type DataGuardianInterface interface {
	SelectAll(search string) ([]GuardianCore, error)
	SelectByID(id string) (*GuardianCore, error)
	SelectByUserID(userID string) (*GuardianCore, error)
	Insert(input *GuardianCore) error
	Update(id string, input *GuardianCore) error
	Delete(id string) error

	// Relasi Student-Guardian
	LinkStudent(input *StudentGuardianCore) error
	UnlinkStudent(guardianID string, siswaID string) error
	SelectStudentsByGuardianID(guardianID string) ([]StudentGuardianCore, error)
	SelectGuardiansBySiswaID(siswaID string) ([]StudentGuardianCore, error)
}

type ServiceGuardianInterface interface {
	GetAll(search string) ([]GuardianCore, error)
	GetByID(id string) (*GuardianCore, error)
	GetByUserID(userID string) (*GuardianCore, error)
	Create(input *GuardianCore) (*GuardianCore, error)
	Update(id string, input *GuardianCore) (*GuardianCore, error)
	Delete(id string) error

	// Relasi Student-Guardian
	LinkStudent(guardianID string, siswaID string, isPrimary bool) (*StudentGuardianCore, error)
	UnlinkStudent(guardianID string, siswaID string) error
	GetStudentsByGuardianID(guardianID string) ([]StudentGuardianCore, error)
	GetGuardiansBySiswaID(siswaID string) ([]StudentGuardianCore, error)
}
