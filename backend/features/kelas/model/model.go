package model

import (
	"database/sql"
	"go_rest_native_sekolah/features/kelas"
	"time"
)

type Kelas struct {
	ID          string
	Kelas       string
	KodeKelas   sql.NullString
	Tingkat     sql.NullString
	Fase        sql.NullString
	Jurusan     sql.NullString
	TahunAjaran sql.NullString
	Ruangan     sql.NullString
	Kapasitas   sql.NullInt32
	Status      sql.NullString
	Deskripsi   sql.NullString
	ID_Guru     sql.NullString
	Nama_Guru   sql.NullString
	TotalSiswa  int
	Update_At   time.Time
	Delete_At   *time.Time
}

func (k *Kelas) TableName() string {
	return "kelas"
}

func FormatterResponse(res Kelas) kelas.KelasCore {
	kapasitas := 36
	if res.Kapasitas.Valid && res.Kapasitas.Int32 > 0 {
		kapasitas = int(res.Kapasitas.Int32)
	}

	status := "Aktif"
	if res.Status.Valid && res.Status.String != "" {
		status = res.Status.String
	}

	tahunAjaran := "2026/2027"
	if res.TahunAjaran.Valid && res.TahunAjaran.String != "" {
		tahunAjaran = res.TahunAjaran.String
	}

	upAt := ""
	if !res.Update_At.IsZero() {
		upAt = res.Update_At.Format("2006-01-02 15:04:05")
	}

	delAt := ""
	if res.Delete_At != nil {
		delAt = res.Delete_At.Format("2006-01-02 15:04:05")
	}

	return kelas.KelasCore{
		ID:          res.ID,
		Kelas:       res.Kelas,
		KodeKelas:   res.KodeKelas.String,
		Tingkat:     res.Tingkat.String,
		Fase:        res.Fase.String,
		Jurusan:     res.Jurusan.String,
		TahunAjaran: tahunAjaran,
		Ruangan:     res.Ruangan.String,
		Kapasitas:   kapasitas,
		Status:      status,
		Deskripsi:   res.Deskripsi.String,
		ID_Guru:     res.ID_Guru.String,
		Nama_Guru:   res.Nama_Guru.String,
		TotalSiswa:  res.TotalSiswa,
		Update_At:   upAt,
		Delete_At:   delAt,
	}
}
