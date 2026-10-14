package controllers

import "go_rest_native_sekolah/features/guru"

// GuruFormatter digunakan untuk memformat data guru agar sesuai kebutuhan response API
type GuruFormatter struct {
	ID                 string `json:"id"`
	ID_User            string `json:"id_user"`
	Nama               string `json:"nama"`
	NIP                string `json:"nip"`
	NUPTK              string `json:"nuptk"`
	JenisKelamin       string `json:"jenis_kelamin"`
	TempatLahir        string `json:"tempat_lahir"`
	TanggalLahir       string `json:"tanggal_lahir"`
	NoTelepon          string `json:"no_telepon"`
	Agama              string `json:"agama"`
	PendidikanTerakhir string `json:"pendidikan_terakhir"`
	StatusKepegawaian  string `json:"status_kepegawaian"`
	Email              string `json:"email"`
	Alamat             string `json:"alamat"`
}

func FormatGuruList(cores []guru.GuruCore) []GuruFormatter {
	formatted := make([]GuruFormatter, 0, len(cores))
	for _, core := range cores {
		formatted = append(formatted, GuruFormatter{
			ID:                 core.ID,
			ID_User:            core.ID_User,
			Nama:               core.Nama,
			NIP:                core.NIP,
			NUPTK:              core.NUPTK,
			JenisKelamin:       core.JenisKelamin,
			TempatLahir:        core.TempatLahir,
			TanggalLahir:       core.TanggalLahir,
			NoTelepon:          core.NoTelepon,
			Agama:              core.Agama,
			PendidikanTerakhir: core.PendidikanTerakhir,
			StatusKepegawaian:  core.StatusKepegawaian,
			Email:              core.Email,
			Alamat:             core.Alamat,
		})
	}
	return formatted
}

func FormatGuruRequestToCore(req GuruFormatter) guru.GuruCore {
	return guru.GuruCore{
		ID:                 req.ID,
		ID_User:            req.ID_User,
		Nama:               req.Nama,
		NIP:                req.NIP,
		NUPTK:              req.NUPTK,
		JenisKelamin:       req.JenisKelamin,
		TempatLahir:        req.TempatLahir,
		TanggalLahir:       req.TanggalLahir,
		NoTelepon:          req.NoTelepon,
		Agama:              req.Agama,
		PendidikanTerakhir: req.PendidikanTerakhir,
		StatusKepegawaian:  req.StatusKepegawaian,
		Email:              req.Email,
		Alamat:             req.Alamat,
	}
}
