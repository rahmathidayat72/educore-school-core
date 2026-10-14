package matapelajaran

// MataPelajaranCore merepresentasikan data lengkap mata pelajaran
type MataPelajaranCore struct {
	ID             string  `json:"id"`
	Nama_Pelajaran string  `json:"mata_pelajaran"`
	KodeMapel      string  `json:"kode_mapel"`
	Kelompok       string  `json:"kelompok"`
	Kurikulum      string  `json:"kurikulum"`
	Tingkat        string  `json:"tingkat"`
	Fase           string  `json:"fase"`
	JamPerMinggu   int     `json:"jam_per_minggu"`
	KKM            float64 `json:"kkm"`
	Status         string  `json:"status"`
	ID_Guru        string  `json:"id_guru"`
	Guru           string  `json:"guru"`
	Kelas_ID       string  `json:"kelas_id"`
	Nama_Kelas     string  `json:"nama_kelas"`
	Deskripsi      string  `json:"deskripsi"`
	Update_At      string  `json:"update_at"`
	Delete_At      string  `json:"delete_at"`
}

// DataMataPelajaranInterface adalah interface untuk repository mata pelajaran
type DataMataPelajaranInterface interface {
	SelectAllMapel() ([]MataPelajaranCore, error)
	SelectMapelById(id string) (*MataPelajaranCore, error)
	InsertMapel(insert *MataPelajaranCore) error
	UpdateMapel(insert *MataPelajaranCore, id string) error
	DeleteMapel(id string) error
}

// ServiceMapelInterface adalah interface untuk service mata pelajaran
type ServiceMapelInterface interface {
	SelectAllMapel() ([]MataPelajaranCore, error)
	SelectMapelById(id string) (*MataPelajaranCore, error)
	InsertMapel(insert *MataPelajaranCore) error
	UpdateMapel(insert *MataPelajaranCore, id string) error
	DeleteMapel(id string) error
}
