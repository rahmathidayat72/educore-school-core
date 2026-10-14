package controllers

import (
	"go_rest_native_sekolah/features/kelas"
)

type KelasFormatter struct {
	ID          string `json:"id"`
	Kelas       string `json:"kelas"`
	KodeKelas   string `json:"kode_kelas"`
	Tingkat     string `json:"tingkat"`
	Fase        string `json:"fase"`
	Jurusan     string `json:"jurusan"`
	TahunAjaran string `json:"tahun_ajaran"`
	Ruangan     string `json:"ruangan"`
	Kapasitas   int    `json:"kapasitas"`
	Status      string `json:"status"`
	Deskripsi   string `json:"deskripsi"`
	ID_Guru     string `json:"id_guru"`
	Nama_Guru   string `json:"nama_guru"`
	TotalSiswa  int    `json:"total_siswa"`
}

func FormatKelasList(cores []kelas.KelasCore) []KelasFormatter {
	formatted := make([]KelasFormatter, 0, len(cores))
	for _, core := range cores {
		formatted = append(formatted, KelasFormatter{
			ID:          core.ID,
			Kelas:       core.Kelas,
			KodeKelas:   core.KodeKelas,
			Tingkat:     core.Tingkat,
			Fase:        core.Fase,
			Jurusan:     core.Jurusan,
			TahunAjaran: core.TahunAjaran,
			Ruangan:     core.Ruangan,
			Kapasitas:   core.Kapasitas,
			Status:      core.Status,
			Deskripsi:   core.Deskripsi,
			ID_Guru:     core.ID_Guru,
			Nama_Guru:   core.Nama_Guru,
			TotalSiswa:  core.TotalSiswa,
		})
	}
	return formatted
}

func FormatKelasRequestToCore(req KelasFormatter) kelas.KelasCore {
	return kelas.KelasCore{
		ID:          req.ID,
		Kelas:       req.Kelas,
		KodeKelas:   req.KodeKelas,
		Tingkat:     req.Tingkat,
		Fase:        req.Fase,
		Jurusan:     req.Jurusan,
		TahunAjaran: req.TahunAjaran,
		Ruangan:     req.Ruangan,
		Kapasitas:   req.Kapasitas,
		Status:      req.Status,
		Deskripsi:   req.Deskripsi,
		ID_Guru:     req.ID_Guru,
		Nama_Guru:   req.Nama_Guru,
	}
}
