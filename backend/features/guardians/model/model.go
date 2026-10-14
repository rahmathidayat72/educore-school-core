package model

import (
	"database/sql"
	"go_rest_native_sekolah/features/guardians"
	"time"
)

type Guardian struct {
	ID          string
	UserID      sql.NullString
	Nama        string
	NIK         sql.NullString
	Hubungan    string
	Pekerjaan   sql.NullString
	Penghasilan sql.NullString
	Telepon     string
	Email       sql.NullString
	Alamat      sql.NullString
	CreatedAt   *time.Time
	UpdateAt    *time.Time
	DeleteAt    *time.Time
}

func (m *Guardian) toCore() guardians.GuardianCore {
	core := guardians.GuardianCore{
		ID:        m.ID,
		Nama:      m.Nama,
		Hubungan:  m.Hubungan,
		Telepon:   m.Telepon,
		CreatedAt: m.CreatedAt,
		UpdateAt:  m.UpdateAt,
		DeleteAt:  m.DeleteAt,
	}
	if m.UserID.Valid {
		val := m.UserID.String
		core.UserID = &val
	}
	if m.NIK.Valid {
		core.NIK = m.NIK.String
	}
	if m.Pekerjaan.Valid {
		core.Pekerjaan = m.Pekerjaan.String
	}
	if m.Penghasilan.Valid {
		core.Penghasilan = m.Penghasilan.String
	}
	if m.Email.Valid {
		core.Email = m.Email.String
	}
	if m.Alamat.Valid {
		core.Alamat = m.Alamat.String
	}
	return core
}

func toGuardianCoreList(list []Guardian) []guardians.GuardianCore {
	cores := make([]guardians.GuardianCore, 0, len(list))
	for _, item := range list {
		cores = append(cores, item.toCore())
	}
	return cores
}

type StudentGuardian struct {
	ID          string
	GuardianID  string
	SiswaID     string
	IsPrimary   bool
	CreatedAt   *time.Time
	UpdateAt    *time.Time
	DeleteAt    *time.Time

	// Join fields
	NamaSiswa   sql.NullString
	NISN        sql.NullString
	NamaKelas   sql.NullString
	NamaSekolah sql.NullString
	NamaWali    sql.NullString
	Hubungan    sql.NullString
	TeleponWali sql.NullString
}

func (m *StudentGuardian) toCore() guardians.StudentGuardianCore {
	core := guardians.StudentGuardianCore{
		ID:         m.ID,
		GuardianID: m.GuardianID,
		SiswaID:    m.SiswaID,
		IsPrimary:  m.IsPrimary,
		CreatedAt:  m.CreatedAt,
		UpdateAt:   m.UpdateAt,
		DeleteAt:   m.DeleteAt,
	}
	if m.NamaSiswa.Valid {
		core.NamaSiswa = m.NamaSiswa.String
	}
	if m.NISN.Valid {
		core.NISN = m.NISN.String
	}
	if m.NamaKelas.Valid {
		core.NamaKelas = m.NamaKelas.String
	}
	if m.NamaSekolah.Valid {
		core.NamaSekolah = m.NamaSekolah.String
	}
	if m.NamaWali.Valid {
		core.NamaWali = m.NamaWali.String
	}
	if m.Hubungan.Valid {
		core.Hubungan = m.Hubungan.String
	}
	if m.TeleponWali.Valid {
		core.TeleponWali = m.TeleponWali.String
	}
	return core
}

func toStudentGuardianCoreList(list []StudentGuardian) []guardians.StudentGuardianCore {
	cores := make([]guardians.StudentGuardianCore, 0, len(list))
	for _, item := range list {
		cores = append(cores, item.toCore())
	}
	return cores
}
