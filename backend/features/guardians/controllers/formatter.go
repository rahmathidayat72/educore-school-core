package controllers

import (
	"go_rest_native_sekolah/features/guardians"
	"time"
)

type GuardianRequest struct {
	UserID      *string `json:"user_id"`
	Nama        string  `json:"nama"`
	NIK         string  `json:"nik"`
	Hubungan    string  `json:"hubungan"`
	Pekerjaan   string  `json:"pekerjaan"`
	Penghasilan string  `json:"penghasilan"`
	Telepon     string  `json:"telepon"`
	Email       string  `json:"email"`
	Alamat      string  `json:"alamat"`
}

type GuardianResponse struct {
	ID          string                    `json:"id"`
	UserID      *string                   `json:"user_id,omitempty"`
	Nama        string                    `json:"nama"`
	NIK         string                    `json:"nik,omitempty"`
	Hubungan    string                    `json:"hubungan"`
	Pekerjaan   string                    `json:"pekerjaan,omitempty"`
	Penghasilan string                    `json:"penghasilan,omitempty"`
	Telepon     string                    `json:"telepon"`
	Email       string                    `json:"email,omitempty"`
	Alamat      string                    `json:"alamat,omitempty"`
	CreatedAt   *time.Time                `json:"created_at,omitempty"`
	UpdateAt    *time.Time                `json:"update_at,omitempty"`
	SiswaList   []StudentGuardianResponse `json:"siswa_list,omitempty"`
}

type LinkStudentRequest struct {
	GuardianID string `json:"guardian_id"`
	SiswaID    string `json:"siswa_id"`
	IsPrimary  bool   `json:"is_primary"`
}

type StudentGuardianResponse struct {
	ID          string     `json:"id"`
	GuardianID  string     `json:"guardian_id"`
	SiswaID     string     `json:"siswa_id"`
	IsPrimary   bool       `json:"is_primary"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdateAt    *time.Time `json:"update_at,omitempty"`
	NamaSiswa   string     `json:"nama_siswa,omitempty"`
	NISN        string     `json:"nisn,omitempty"`
	NamaKelas   string     `json:"nama_kelas,omitempty"`
	NamaSekolah string     `json:"nama_sekolah,omitempty"`
	NamaWali    string     `json:"nama_wali,omitempty"`
	Hubungan    string     `json:"hubungan,omitempty"`
	TeleponWali string     `json:"telepon_wali,omitempty"`
}

func ToGuardianCore(req GuardianRequest) guardians.GuardianCore {
	return guardians.GuardianCore{
		UserID:      req.UserID,
		Nama:        req.Nama,
		NIK:         req.NIK,
		Hubungan:    req.Hubungan,
		Pekerjaan:   req.Pekerjaan,
		Penghasilan: req.Penghasilan,
		Telepon:     req.Telepon,
		Email:       req.Email,
		Alamat:      req.Alamat,
	}
}

func ToStudentGuardianResponse(core guardians.StudentGuardianCore) StudentGuardianResponse {
	return StudentGuardianResponse{
		ID:          core.ID,
		GuardianID:  core.GuardianID,
		SiswaID:     core.SiswaID,
		IsPrimary:   core.IsPrimary,
		CreatedAt:   core.CreatedAt,
		UpdateAt:    core.UpdateAt,
		NamaSiswa:   core.NamaSiswa,
		NISN:        core.NISN,
		NamaKelas:   core.NamaKelas,
		NamaSekolah: core.NamaSekolah,
		NamaWali:    core.NamaWali,
		Hubungan:    core.Hubungan,
		TeleponWali: core.TeleponWali,
	}
}

func ToStudentGuardianResponseList(cores []guardians.StudentGuardianCore) []StudentGuardianResponse {
	res := make([]StudentGuardianResponse, 0, len(cores))
	for _, c := range cores {
		res = append(res, ToStudentGuardianResponse(c))
	}
	return res
}

func ToGuardianResponse(core guardians.GuardianCore) GuardianResponse {
	resp := GuardianResponse{
		ID:          core.ID,
		UserID:      core.UserID,
		Nama:        core.Nama,
		NIK:         core.NIK,
		Hubungan:    core.Hubungan,
		Pekerjaan:   core.Pekerjaan,
		Penghasilan: core.Penghasilan,
		Telepon:     core.Telepon,
		Email:       core.Email,
		Alamat:      core.Alamat,
		CreatedAt:   core.CreatedAt,
		UpdateAt:    core.UpdateAt,
	}
	if len(core.SiswaList) > 0 {
		resp.SiswaList = ToStudentGuardianResponseList(core.SiswaList)
	}
	return resp
}

func ToGuardianResponseList(cores []guardians.GuardianCore) []GuardianResponse {
	res := make([]GuardianResponse, 0, len(cores))
	for _, c := range cores {
		res = append(res, ToGuardianResponse(c))
	}
	return res
}
