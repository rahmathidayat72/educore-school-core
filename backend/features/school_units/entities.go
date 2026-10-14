package schoolunits

import "time"

type SchoolUnitCore struct {
	ID          string     `json:"id"`
	NamaSekolah string     `json:"nama_sekolah"`
	Jenjang     string     `json:"jenjang"` // SD, SMP, SMA, SMK
	NPSN        string     `json:"npsn"`
	Alamat      string     `json:"alamat"`
	Telepon     string     `json:"telepon"`
	Email       string     `json:"email"`
	Website     string     `json:"website"`
	Status      string     `json:"status"` // aktif, nonaktif, dalam_pembangunan
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdateAt    *time.Time `json:"update_at,omitempty"`
	DeleteAt    *time.Time `json:"delete_at,omitempty"`
}

type DataSchoolUnitInterface interface {
	SelectAll(jenjang, status string) ([]SchoolUnitCore, error)
	SelectByID(id string) (*SchoolUnitCore, error)
	Insert(input *SchoolUnitCore) error
	Update(id string, input *SchoolUnitCore) error
	Delete(id string) error
	CheckNPSNExists(npsn string, excludeID string) (bool, error)
}

type ServiceSchoolUnitInterface interface {
	GetAll(jenjang, status string) ([]SchoolUnitCore, error)
	GetByID(id string) (*SchoolUnitCore, error)
	Create(input *SchoolUnitCore) (*SchoolUnitCore, error)
	Update(id string, input *SchoolUnitCore) (*SchoolUnitCore, error)
	Delete(id string) error
}
