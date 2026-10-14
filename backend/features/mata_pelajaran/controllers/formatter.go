package controllers

import matapelajaran "go_rest_native_sekolah/features/mata_pelajaran"

// FormatterMataPelajaran digunakan untuk memformat data mata pelajaran lengkap
type FormatterMataPelajaran struct {
	ID             string  `json:"id"`
	SchoolUnitID   *string `json:"school_unit_id,omitempty"`
	Nama_Pelajaran string  `json:"mata_pelajaran"`
	KodeMapel      string  `json:"kode_mapel"`
	Kelompok       string  `json:"kelompok"`
	Kurikulum      string  `json:"kurikulum"`
	Tingkat        string  `json:"tingkat"`
	Fase           string  `json:"fase"`
	JamPerMinggu   int     `json:"jam_per_minggu"`
	KKM            float64 `json:"kkm"`
	Status         string  `json:"status"`
	ID_Guru        string  `json:"id_guru"`
	Guru           string  `json:"guru"`
	Kelas_ID       string  `json:"kelas_id"`
	Nama_Kelas     string  `json:"nama_kelas"`
	Deskripsi      string  `json:"deskripsi"`
}

func FormatterMapelList(cores []matapelajaran.MataPelajaranCore) []FormatterMataPelajaran {
	formatted := make([]FormatterMataPelajaran, 0, len(cores))
	for _, core := range cores {
		formatted = append(formatted, FormatterMataPelajaran{
			ID:             core.ID,
			SchoolUnitID:   core.SchoolUnitID,
			Nama_Pelajaran: core.Nama_Pelajaran,
			KodeMapel:      core.KodeMapel,
			Kelompok:       core.Kelompok,
			Kurikulum:      core.Kurikulum,
			Tingkat:        core.Tingkat,
			Fase:           core.Fase,
			JamPerMinggu:   core.JamPerMinggu,
			KKM:            core.KKM,
			Status:         core.Status,
			ID_Guru:        core.ID_Guru,
			Guru:           core.Guru,
			Kelas_ID:       core.Kelas_ID,
			Nama_Kelas:     core.Nama_Kelas,
			Deskripsi:      core.Deskripsi,
		})
	}
	return formatted
}

func FormatterMapelRequestToCore(req FormatterMataPelajaran) matapelajaran.MataPelajaranCore {
	return matapelajaran.MataPelajaranCore{
		ID:             req.ID,
		SchoolUnitID:   req.SchoolUnitID,
		Nama_Pelajaran: req.Nama_Pelajaran,
		KodeMapel:      req.KodeMapel,
		Kelompok:       req.Kelompok,
		Kurikulum:      req.Kurikulum,
		Tingkat:        req.Tingkat,
		Fase:           req.Fase,
		JamPerMinggu:   req.JamPerMinggu,
		KKM:            req.KKM,
		Status:         req.Status,
		ID_Guru:        req.ID_Guru,
		Guru:           req.Guru,
		Kelas_ID:       req.Kelas_ID,
		Nama_Kelas:     req.Nama_Kelas,
		Deskripsi:      req.Deskripsi,
	}
}
