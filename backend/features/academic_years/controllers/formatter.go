package controllers

import (
	academicyears "go_rest_native_sekolah/features/academic_years"
	"time"
)

type AcademicYearRequest struct {
	TahunAjaran    string `json:"tahun_ajaran"`
	Semester       string `json:"semester"`
	IsActive       bool   `json:"is_active"`
	TanggalMulai   string `json:"tanggal_mulai"`
	TanggalSelesai string `json:"tanggal_selesai"`
}

type AcademicYearResponse struct {
	ID             string     `json:"id"`
	SchoolUnitID   string     `json:"school_unit_id"`
	TahunAjaran    string     `json:"tahun_ajaran"`
	Semester       string     `json:"semester"`
	IsActive       bool       `json:"is_active"`
	TanggalMulai   *string    `json:"tanggal_mulai,omitempty"`
	TanggalSelesai *string    `json:"tanggal_selesai,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdateAt       *time.Time `json:"update_at,omitempty"`
}

func ToYearCore(req AcademicYearRequest, schoolUnitID string) academicyears.AcademicYearCore {
	core := academicyears.AcademicYearCore{
		SchoolUnitID: schoolUnitID,
		TahunAjaran:  req.TahunAjaran,
		Semester:     req.Semester,
		IsActive:     req.IsActive,
	}
	if req.TanggalMulai != "" {
		core.TanggalMulai = &req.TanggalMulai
	}
	if req.TanggalSelesai != "" {
		core.TanggalSelesai = &req.TanggalSelesai
	}
	return core
}

func ToYearResponse(core academicyears.AcademicYearCore) AcademicYearResponse {
	return AcademicYearResponse{
		ID:             core.ID,
		SchoolUnitID:   core.SchoolUnitID,
		TahunAjaran:    core.TahunAjaran,
		Semester:       core.Semester,
		IsActive:       core.IsActive,
		TanggalMulai:   core.TanggalMulai,
		TanggalSelesai: core.TanggalSelesai,
		CreatedAt:      core.CreatedAt,
		UpdateAt:       core.UpdateAt,
	}
}

func ToYearResponseList(cores []academicyears.AcademicYearCore) []AcademicYearResponse {
	res := make([]AcademicYearResponse, 0, len(cores))
	for _, c := range cores {
		res = append(res, ToYearResponse(c))
	}
	return res
}

// Enrollments DTOs

type ClassEnrollmentRequest struct {
	AcademicYearID string `json:"academic_year_id"`
	KelasID        string `json:"kelas_id"`
	SiswaID        string `json:"siswa_id"`
	Status         string `json:"status"`
	Keterangan     string `json:"keterangan"`
}

type ClassEnrollmentResponse struct {
	ID             string     `json:"id"`
	SchoolUnitID   string     `json:"school_unit_id"`
	AcademicYearID string     `json:"academic_year_id"`
	KelasID        string     `json:"kelas_id"`
	SiswaID        string     `json:"siswa_id"`
	Status         string     `json:"status"`
	Keterangan     string     `json:"keterangan,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdateAt       *time.Time `json:"update_at,omitempty"`
	NamaSiswa      string     `json:"nama_siswa,omitempty"`
	NIS            string     `json:"nis,omitempty"`
	NISN           string     `json:"nisn,omitempty"`
	NamaKelas      string     `json:"nama_kelas,omitempty"`
	Tingkat        string     `json:"tingkat,omitempty"`
	TahunAjaran    string     `json:"tahun_ajaran,omitempty"`
	Semester       string     `json:"semester,omitempty"`
}

func ToEnrollmentCore(req ClassEnrollmentRequest, schoolUnitID string) academicyears.ClassEnrollmentCore {
	return academicyears.ClassEnrollmentCore{
		SchoolUnitID:   schoolUnitID,
		AcademicYearID: req.AcademicYearID,
		KelasID:        req.KelasID,
		SiswaID:        req.SiswaID,
		Status:         req.Status,
		Keterangan:     req.Keterangan,
	}
}

func ToEnrollmentResponse(core academicyears.ClassEnrollmentCore) ClassEnrollmentResponse {
	return ClassEnrollmentResponse{
		ID:             core.ID,
		SchoolUnitID:   core.SchoolUnitID,
		AcademicYearID: core.AcademicYearID,
		KelasID:        core.KelasID,
		SiswaID:        core.SiswaID,
		Status:         core.Status,
		Keterangan:     core.Keterangan,
		CreatedAt:      core.CreatedAt,
		UpdateAt:       core.UpdateAt,
		NamaSiswa:      core.NamaSiswa,
		NIS:            core.NIS,
		NISN:           core.NISN,
		NamaKelas:      core.NamaKelas,
		Tingkat:        core.Tingkat,
		TahunAjaran:    core.TahunAjaran,
		Semester:       core.Semester,
	}
}

func ToEnrollmentResponseList(cores []academicyears.ClassEnrollmentCore) []ClassEnrollmentResponse {
	res := make([]ClassEnrollmentResponse, 0, len(cores))
	for _, c := range cores {
		res = append(res, ToEnrollmentResponse(c))
	}
	return res
}
