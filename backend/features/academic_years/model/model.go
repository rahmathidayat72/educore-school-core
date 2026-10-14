package model

import (
	"database/sql"
	academicyears "go_rest_native_sekolah/features/academic_years"
	"time"
)

type AcademicYear struct {
	ID             string
	SchoolUnitID   string
	TahunAjaran    string
	Semester       string
	IsActive       bool
	TanggalMulai   sql.NullString
	TanggalSelesai sql.NullString
	CreatedAt      *time.Time
	UpdateAt       *time.Time
	DeleteAt       *time.Time
}

func (m *AcademicYear) toCore() academicyears.AcademicYearCore {
	core := academicyears.AcademicYearCore{
		ID:           m.ID,
		SchoolUnitID: m.SchoolUnitID,
		TahunAjaran:  m.TahunAjaran,
		Semester:     m.Semester,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdateAt:     m.UpdateAt,
		DeleteAt:     m.DeleteAt,
	}
	if m.TanggalMulai.Valid {
		val := m.TanggalMulai.String
		core.TanggalMulai = &val
	}
	if m.TanggalSelesai.Valid {
		val := m.TanggalSelesai.String
		core.TanggalSelesai = &val
	}
	return core
}

func toAcademicYearCoreList(list []AcademicYear) []academicyears.AcademicYearCore {
	cores := make([]academicyears.AcademicYearCore, 0, len(list))
	for _, item := range list {
		cores = append(cores, item.toCore())
	}
	return cores
}

type ClassEnrollment struct {
	ID             string
	SchoolUnitID   string
	AcademicYearID string
	KelasID        string
	SiswaID        string
	Status         string
	Keterangan     sql.NullString
	CreatedAt      *time.Time
	UpdateAt       *time.Time
	DeleteAt       *time.Time

	// Join fields
	NamaSiswa   sql.NullString
	NIS         sql.NullString
	NISN        sql.NullString
	NamaKelas   sql.NullString
	Tingkat     sql.NullString
	TahunAjaran sql.NullString
	Semester    sql.NullString
}

func (m *ClassEnrollment) toCore() academicyears.ClassEnrollmentCore {
	core := academicyears.ClassEnrollmentCore{
		ID:             m.ID,
		SchoolUnitID:   m.SchoolUnitID,
		AcademicYearID: m.AcademicYearID,
		KelasID:        m.KelasID,
		SiswaID:        m.SiswaID,
		Status:         m.Status,
		CreatedAt:      m.CreatedAt,
		UpdateAt:       m.UpdateAt,
		DeleteAt:       m.DeleteAt,
	}
	if m.Keterangan.Valid {
		core.Keterangan = m.Keterangan.String
	}
	if m.NamaSiswa.Valid {
		core.NamaSiswa = m.NamaSiswa.String
	}
	if m.NIS.Valid {
		core.NIS = m.NIS.String
	}
	if m.NISN.Valid {
		core.NISN = m.NISN.String
	}
	if m.NamaKelas.Valid {
		core.NamaKelas = m.NamaKelas.String
	}
	if m.Tingkat.Valid {
		core.Tingkat = m.Tingkat.String
	}
	if m.TahunAjaran.Valid {
		core.TahunAjaran = m.TahunAjaran.String
	}
	if m.Semester.Valid {
		core.Semester = m.Semester.String
	}
	return core
}

func toEnrollmentCoreList(list []ClassEnrollment) []academicyears.ClassEnrollmentCore {
	cores := make([]academicyears.ClassEnrollmentCore, 0, len(list))
	for _, item := range list {
		cores = append(cores, item.toCore())
	}
	return cores
}
