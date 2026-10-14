package controllers

import "go_rest_native_sekolah/features/siswa"

// SiswaFormatter memformat respons dan request data siswa
type SiswaFormatter struct {
	ID            string `json:"id"`
	Nama          string `json:"nama"`
	Kelas_ID      string `json:"kelas_id"`
	Nama_Kelas    string `json:"nama_kelas"`
	NIS           string `json:"nis"`
	NISN          string `json:"nisn"`
	JenisKelamin  string `json:"jenis_kelamin"`
	TempatLahir   string `json:"tempat_lahir"`
	TanggalLahir  string `json:"tanggal_lahir"`
	Agama         string `json:"agama"`
	NoTelepon     string `json:"no_telepon"`
	NamaWali      string `json:"nama_wali"`
	NoTeleponWali string `json:"no_telepon_wali"`
	TahunMasuk    string `json:"tahun_masuk"`
	Status        string `json:"status"`
	Email         string `json:"email"`
	Alamat        string `json:"alamat"`
}

func FormatterKelasList(cores []siswa.SiswaCore) []SiswaFormatter {
	formatted := make([]SiswaFormatter, 0, len(cores))
	for _, core := range cores {
		formatted = append(formatted, SiswaFormatter{
			ID:            core.ID,
			Nama:          core.Nama,
			Kelas_ID:      core.Kelas_ID,
			Nama_Kelas:    core.Nama_Kelas,
			NIS:           core.NIS,
			NISN:          core.NISN,
			JenisKelamin:  core.JenisKelamin,
			TempatLahir:   core.TempatLahir,
			TanggalLahir:  core.TanggalLahir,
			Agama:         core.Agama,
			NoTelepon:     core.NoTelepon,
			NamaWali:      core.NamaWali,
			NoTeleponWali: core.NoTeleponWali,
			TahunMasuk:    core.TahunMasuk,
			Status:        core.Status,
			Email:         core.Email,
			Alamat:        core.Alamat,
		})
	}
	return formatted
}

func FormatSiswaRequestToCore(req SiswaFormatter) siswa.SiswaCore {
	return siswa.SiswaCore{
		ID:            req.ID,
		Nama:          req.Nama,
		Kelas_ID:      req.Kelas_ID,
		Nama_Kelas:    req.Nama_Kelas,
		NIS:           req.NIS,
		NISN:          req.NISN,
		JenisKelamin:  req.JenisKelamin,
		TempatLahir:   req.TempatLahir,
		TanggalLahir:  req.TanggalLahir,
		Agama:         req.Agama,
		NoTelepon:     req.NoTelepon,
		NamaWali:      req.NamaWali,
		NoTeleponWali: req.NoTeleponWali,
		TahunMasuk:    req.TahunMasuk,
		Status:        req.Status,
		Email:         req.Email,
		Alamat:        req.Alamat,
	}
}
