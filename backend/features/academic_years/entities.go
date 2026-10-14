package academicyears

import "time"

// AcademicYearCore merepresentasikan entitas tahun ajaran sekolah
type AcademicYearCore struct {
	ID             string     `json:"id"`
	SchoolUnitID   string     `json:"school_unit_id"`
	TahunAjaran    string     `json:"tahun_ajaran"` // Contoh: '2025/2026'
	Semester       string     `json:"semester"`     // 'Ganjil', 'Genap'
	IsActive       bool       `json:"is_active"`
	TanggalMulai   *string    `json:"tanggal_mulai,omitempty"`   // YYYY-MM-DD
	TanggalSelesai *string    `json:"tanggal_selesai,omitempty"` // YYYY-MM-DD
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdateAt       *time.Time `json:"update_at,omitempty"`
	DeleteAt       *time.Time `json:"delete_at,omitempty"`
}

// ClassEnrollmentCore merepresentasikan riwayat penempatan siswa di rombel/kelas
type ClassEnrollmentCore struct {
	ID             string     `json:"id"`
	SchoolUnitID   string     `json:"school_unit_id"`
	AcademicYearID string     `json:"academic_year_id"`
	KelasID        string     `json:"kelas_id"`
	SiswaID        string     `json:"siswa_id"`
	Status         string     `json:"status"` // 'Aktif', 'Lulus', 'Pindah', 'Tinggal Kelas'
	Keterangan     string     `json:"keterangan,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdateAt       *time.Time `json:"update_at,omitempty"`
	DeleteAt       *time.Time `json:"delete_at,omitempty"`

	// Enriched fields for display
	NamaSiswa   string `json:"nama_siswa,omitempty"`
	NIS         string `json:"nis,omitempty"`
	NISN        string `json:"nisn,omitempty"`
	NamaKelas   string `json:"nama_kelas,omitempty"`
	Tingkat     string `json:"tingkat,omitempty"`
	TahunAjaran string `json:"tahun_ajaran,omitempty"`
	Semester    string `json:"semester,omitempty"`
}

// DataAcademicYearInterface mendefinisikan kontrak akses data repository
type DataAcademicYearInterface interface {
	SelectAllYears(schoolUnitID string) ([]AcademicYearCore, error)
	SelectYearByID(id string, schoolUnitID ...string) (*AcademicYearCore, error)
	GetActiveYear(schoolUnitID string) (*AcademicYearCore, error)
	InsertYear(input *AcademicYearCore) error
	UpdateYear(id string, input *AcademicYearCore, schoolUnitID ...string) error
	SetActiveYear(id string, schoolUnitID string) error
	DeleteYear(id string, schoolUnitID ...string) error

	// Enrollments
	SelectEnrollments(schoolUnitID string, academicYearID string, kelasID string) ([]ClassEnrollmentCore, error)
	SelectEnrollmentByID(id string, schoolUnitID ...string) (*ClassEnrollmentCore, error)
	InsertEnrollment(input *ClassEnrollmentCore) error
	UpdateEnrollment(id string, input *ClassEnrollmentCore, schoolUnitID ...string) error
	DeleteEnrollment(id string, schoolUnitID ...string) error
	GetStudentEnrollmentHistory(siswaID string) ([]ClassEnrollmentCore, error)
}

// ServiceAcademicYearInterface mendefinisikan kontrak business logic service
type ServiceAcademicYearInterface interface {
	GetAllYears(schoolUnitID string) ([]AcademicYearCore, error)
	GetYearByID(id string, schoolUnitID ...string) (*AcademicYearCore, error)
	GetActiveYear(schoolUnitID string) (*AcademicYearCore, error)
	CreateYear(input *AcademicYearCore) (*AcademicYearCore, error)
	UpdateYear(id string, input *AcademicYearCore, schoolUnitID ...string) (*AcademicYearCore, error)
	SetActiveYear(id string, schoolUnitID string) error
	DeleteYear(id string, schoolUnitID ...string) error

	// Enrollments
	GetEnrollments(schoolUnitID string, academicYearID string, kelasID string) ([]ClassEnrollmentCore, error)
	GetEnrollmentByID(id string, schoolUnitID ...string) (*ClassEnrollmentCore, error)
	EnrollStudent(input *ClassEnrollmentCore) (*ClassEnrollmentCore, error)
	UpdateEnrollment(id string, input *ClassEnrollmentCore, schoolUnitID ...string) (*ClassEnrollmentCore, error)
	DeleteEnrollment(id string, schoolUnitID ...string) error
	GetStudentEnrollmentHistory(siswaID string) ([]ClassEnrollmentCore, error)
}
