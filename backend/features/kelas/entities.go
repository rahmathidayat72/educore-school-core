package kelas

// KelasCore merepresentasikan data lengkap kelas / rombongan belajar
type KelasCore struct {
	ID           string  `json:"id"`
	SchoolUnitID *string `json:"school_unit_id,omitempty"`
	Kelas        string  `json:"kelas"`
	KodeKelas    string  `json:"kode_kelas"`
	Tingkat      string  `json:"tingkat"`
	Fase         string  `json:"fase"`
	Jurusan      string  `json:"jurusan"`
	TahunAjaran  string  `json:"tahun_ajaran"`
	Ruangan      string  `json:"ruangan"`
	Kapasitas    int     `json:"kapasitas"`
	Status       string  `json:"status"` // Aktif, Nonaktif, Lulus
	Deskripsi    string  `json:"deskripsi"`
	ID_Guru      string  `json:"id_guru"`
	Nama_Guru    string  `json:"nama_guru"`
	TotalSiswa   int     `json:"total_siswa"` // Statistik jumlah siswa terdaftar
	Update_At    string  `json:"update_at"`
	Delete_At    string  `json:"delete_at"`
}

// DataKelasInterface antarmuka repository kelas
type DataKelasInterface interface {
	SelectAll(schoolUnitID ...string) ([]KelasCore, error)
	SelectById(id string, schoolUnitID ...string) (*KelasCore, error)
	Insert(insert *KelasCore) error
	Update(insert *KelasCore, id string, schoolUnitID ...string) error
	DeleteById(id string, schoolUnitID ...string) error
}

// ServiceKelasInterface antarmuka logika bisnis kelas
type ServiceKelasInterface interface {
	SelectAll(schoolUnitID ...string) ([]KelasCore, error)
	SelectById(id string, schoolUnitID ...string) (*KelasCore, error)
	Insert(insert *KelasCore) error
	Update(insert *KelasCore, id string, schoolUnitID ...string) error
	DeleteById(id string, schoolUnitID ...string) error
}
