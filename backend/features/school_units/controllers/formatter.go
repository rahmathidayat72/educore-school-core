package controllers

import (
	schoolunits "go_rest_native_sekolah/features/school_units"
	"time"
)

type SchoolUnitRequest struct {
	NamaSekolah string `json:"nama_sekolah"`
	Jenjang     string `json:"jenjang"`
	NPSN        string `json:"npsn"`
	Alamat      string `json:"alamat"`
	Telepon     string `json:"telepon"`
	Email       string `json:"email"`
	Website     string `json:"website"`
	Status      string `json:"status"`
}

type SchoolUnitResponse struct {
	ID          string     `json:"id"`
	NamaSekolah string     `json:"nama_sekolah"`
	Jenjang     string     `json:"jenjang"`
	NPSN        string     `json:"npsn,omitempty"`
	Alamat      string     `json:"alamat,omitempty"`
	Telepon     string     `json:"telepon,omitempty"`
	Email       string     `json:"email,omitempty"`
	Website     string     `json:"website,omitempty"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdateAt    *time.Time `json:"update_at,omitempty"`
}

func ToCore(req SchoolUnitRequest) schoolunits.SchoolUnitCore {
	return schoolunits.SchoolUnitCore{
		NamaSekolah: req.NamaSekolah,
		Jenjang:     req.Jenjang,
		NPSN:        req.NPSN,
		Alamat:      req.Alamat,
		Telepon:     req.Telepon,
		Email:       req.Email,
		Website:     req.Website,
		Status:      req.Status,
	}
}

func ToResponse(core schoolunits.SchoolUnitCore) SchoolUnitResponse {
	return SchoolUnitResponse{
		ID:          core.ID,
		NamaSekolah: core.NamaSekolah,
		Jenjang:     core.Jenjang,
		NPSN:        core.NPSN,
		Alamat:      core.Alamat,
		Telepon:     core.Telepon,
		Email:       core.Email,
		Website:     core.Website,
		Status:      core.Status,
		CreatedAt:   core.CreatedAt,
		UpdateAt:    core.UpdateAt,
	}
}

func ToResponseList(cores []schoolunits.SchoolUnitCore) []SchoolUnitResponse {
	res := make([]SchoolUnitResponse, 0, len(cores))
	for _, c := range cores {
		res = append(res, ToResponse(c))
	}
	return res
}
