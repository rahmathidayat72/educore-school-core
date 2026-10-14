package model

import (
	"database/sql"
	matapelajaran "go_rest_native_sekolah/features/mata_pelajaran"
	"time"
)

type MataPelajaran struct {
	ID             string
	Nama_Pelajaran string
	KodeMapel      sql.NullString
	Kelompok       sql.NullString
	Kurikulum      sql.NullString
	Tingkat        sql.NullString
	Fase           sql.NullString
	JamPerMinggu   sql.NullInt32
	KKM            sql.NullFloat64
	Status         sql.NullString
	ID_Guru        sql.NullString
	Guru           sql.NullString
	Kelas_ID       sql.NullString
	Nama_Kelas     sql.NullString
	Deskripsi      sql.NullString
	Update_At      time.Time
	Delete_At      *time.Time
}

func (m *MataPelajaran) TableName() string {
	return "mata_pelajaran"
}

func FormatterResponse(res MataPelajaran) matapelajaran.MataPelajaranCore {
	kurikulum := "Kurikulum Merdeka"
	if res.Kurikulum.Valid && res.Kurikulum.String != "" {
		kurikulum = res.Kurikulum.String
	}

	jam := 2
	if res.JamPerMinggu.Valid && res.JamPerMinggu.Int32 > 0 {
		jam = int(res.JamPerMinggu.Int32)
	}

	kkm := 75.0
	if res.KKM.Valid && res.KKM.Float64 > 0 {
		kkm = res.KKM.Float64
	}

	status := "Aktif"
	if res.Status.Valid && res.Status.String != "" {
		status = res.Status.String
	}

	upAt := ""
	if !res.Update_At.IsZero() {
		upAt = res.Update_At.Format("2006-01-02 15:04:05")
	}

	delAt := ""
	if res.Delete_At != nil {
		delAt = res.Delete_At.Format("2006-01-02 15:04:05")
	}

	return matapelajaran.MataPelajaranCore{
		ID:             res.ID,
		Nama_Pelajaran: res.Nama_Pelajaran,
		KodeMapel:      res.KodeMapel.String,
		Kelompok:       res.Kelompok.String,
		Kurikulum:      kurikulum,
		Tingkat:        res.Tingkat.String,
		Fase:           res.Fase.String,
		JamPerMinggu:   jam,
		KKM:            kkm,
		Status:         status,
		ID_Guru:        res.ID_Guru.String,
		Guru:           res.Guru.String,
		Kelas_ID:       res.Kelas_ID.String,
		Nama_Kelas:     res.Nama_Kelas.String,
		Deskripsi:      res.Deskripsi.String,
		Update_At:      upAt,
		Delete_At:      delAt,
	}
}
