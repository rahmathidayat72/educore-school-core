package gurumodels

import (
	"database/sql"
	"go_rest_native_sekolah/features/guru"
	"time"
)

// Guru merepresentasikan data guru di database
type Guru struct {
	ID                 string
	ID_User            string
	SchoolUnitID       sql.NullString
	Nama               string
	NIP                sql.NullString
	NUPTK              sql.NullString
	JenisKelamin       sql.NullString
	TempatLahir        sql.NullString
	TanggalLahir       sql.NullTime
	NoTelepon          sql.NullString
	Agama              sql.NullString
	PendidikanTerakhir sql.NullString
	StatusKepegawaian  sql.NullString
	Email              sql.NullString
	Alamat             sql.NullString
	Update_At          time.Time
	Delete_At          *time.Time
}

func (u *Guru) TableName() string {
	return "guru"
}

// FormatterResponse mengubah data Guru dari database menjadi GuruCore
func FormatterResponse(res Guru) guru.GuruCore {
	tglLahirStr := ""
	if res.TanggalLahir.Valid {
		tglLahirStr = res.TanggalLahir.Time.Format("2006-01-02")
	}

	var schoolUnitID *string
	if res.SchoolUnitID.Valid && res.SchoolUnitID.String != "" {
		s := res.SchoolUnitID.String
		schoolUnitID = &s
	}

	return guru.GuruCore{
		ID:                 res.ID,
		ID_User:            res.ID_User,
		SchoolUnitID:       schoolUnitID,
		Nama:               res.Nama,
		NIP:                res.NIP.String,
		NUPTK:              res.NUPTK.String,
		JenisKelamin:       res.JenisKelamin.String,
		TempatLahir:        res.TempatLahir.String,
		TanggalLahir:       tglLahirStr,
		NoTelepon:          res.NoTelepon.String,
		Agama:              res.Agama.String,
		PendidikanTerakhir: res.PendidikanTerakhir.String,
		StatusKepegawaian:  res.StatusKepegawaian.String,
		Email:              res.Email.String,
		Alamat:             res.Alamat.String,
		Update_At:          res.Update_At,
		Delete_At:          res.Delete_At,
	}
}
