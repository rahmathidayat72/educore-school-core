package model

import (
	"database/sql"
	"go_rest_native_sekolah/features/siswa"
	"time"
)

// Siswa merepresentasikan data siswa di database
type Siswa struct {
	ID            string
	SchoolUnitID  sql.NullString
	Kelas_ID      sql.NullString
	Nama_Kelas    sql.NullString
	Nama          string
	NIS           sql.NullString
	NISN          sql.NullString
	JenisKelamin  sql.NullString
	TempatLahir   sql.NullString
	TanggalLahir  sql.NullTime
	Agama         sql.NullString
	NoTelepon     sql.NullString
	NamaWali      sql.NullString
	NoTeleponWali sql.NullString
	TahunMasuk    sql.NullString
	Status        sql.NullString
	Email         sql.NullString
	Alamat        sql.NullString
	Update_At     time.Time
	Delete_At     *time.Time
}

func (u *Siswa) TableName() string {
	return "siswa"
}

// FormatterResponse mengubah Siswa (db) menjadi siswa.SiswaCore
func FormatterResponse(res Siswa) siswa.SiswaCore {
	tglLahirStr := ""
	if res.TanggalLahir.Valid {
		tglLahirStr = res.TanggalLahir.Time.Format("2006-01-02")
	}

	status := "Aktif"
	if res.Status.Valid && res.Status.String != "" {
		status = res.Status.String
	}

	var schoolUnitID *string
	if res.SchoolUnitID.Valid && res.SchoolUnitID.String != "" {
		s := res.SchoolUnitID.String
		schoolUnitID = &s
	}

	return siswa.SiswaCore{
		ID:            res.ID,
		SchoolUnitID:  schoolUnitID,
		Nama:          res.Nama,
		Kelas_ID:      res.Kelas_ID.String,
		Nama_Kelas:    res.Nama_Kelas.String,
		NIS:           res.NIS.String,
		NISN:          res.NISN.String,
		JenisKelamin:  res.JenisKelamin.String,
		TempatLahir:   res.TempatLahir.String,
		TanggalLahir:  tglLahirStr,
		Agama:         res.Agama.String,
		NoTelepon:     res.NoTelepon.String,
		NamaWali:      res.NamaWali.String,
		NoTeleponWali: res.NoTeleponWali.String,
		TahunMasuk:    res.TahunMasuk.String,
		Status:        status,
		Email:         res.Email.String,
		Alamat:        res.Alamat.String,
		Update_At:     res.Update_At,
		Delete_At:     res.Delete_At,
	}
}
