package controllers

import (
	"encoding/json"
	"errors"
	"go_rest_native_sekolah/features/jadwal"
	"go_rest_native_sekolah/helper"
	"net/http"
	"strconv"
	"strings"
)

type JadwalController struct {
	service jadwal.ServiceJadwalInterface
}

func NewJadwalController(service jadwal.ServiceJadwalInterface) *JadwalController {
	if service == nil {
		panic("jadwal controller: nil service")
	}
	return &JadwalController{service: service}
}

func getAuthMeta(r *http.Request) (role, userID string) {
	authHeader := r.Header.Get("Authorization")
	tokenStr := helper.GetTokenFromAuthorizationHeader(authHeader)
	role = "guru"
	userID = ""
	if tokenStr != "" {
		if meta, err := helper.VerifyTokenHeader(tokenStr); err == nil {
			if meta.Role != "" {
				role = strings.ToLower(strings.TrimSpace(meta.Role))
			}
			userID = meta.ID
		}
	}
	return
}

func (c *JadwalController) GetAll(w http.ResponseWriter, r *http.Request) error {
	role, userID := getAuthMeta(r)
	_, _, unitID, _ := helper.ExtractUserContext(r)
	hari := r.URL.Query().Get("hari")
	kelasID := r.URL.Query().Get("kelas_id")
	guruID := r.URL.Query().Get("guru_id")
	tahun := r.URL.Query().Get("tahun_ajaran")
	semester := r.URL.Query().Get("semester")
	status := r.URL.Query().Get("status")

	var data []jadwal.JadwalCore
	var err error
	if unitID != "" {
		data, err = c.service.GetAll(role, userID, hari, kelasID, guruID, tahun, semester, status, unitID)
	} else {
		data, err = c.service.GetAll(role, userID, hari, kelasID, guruID, tahun, semester, status)
	}
	if err != nil {
		helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Berhasil mengambil data jadwal pelajaran", data)
	return json.NewEncoder(w).Encode(res)
}

func (c *JadwalController) GetByID(w http.ResponseWriter, r *http.Request) error {
	role, userID := getAuthMeta(r)
	_, _, unitID, _ := helper.ExtractUserContext(r)
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter id wajib diisi")
		return errors.New("id empty")
	}

	var data *jadwal.JadwalCore
	var err error
	if unitID != "" {
		data, err = c.service.GetByID(id, role, userID, unitID)
	} else {
		data, err = c.service.GetByID(id, role, userID)
	}
	if err != nil {
		helper.JSONResponse(w, http.StatusNotFound, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Berhasil mengambil rincian jadwal", data)
	return json.NewEncoder(w).Encode(res)
}

func (c *JadwalController) Create(w http.ResponseWriter, r *http.Request) error {
	role, userID := getAuthMeta(r)
	_, _, unitID, _ := helper.ExtractUserContext(r)

	var payload struct {
		SchoolUnitID    *string `json:"school_unit_id,omitempty"`
		KodeJadwal      string  `json:"kode_jadwal"`
		MataPelajaranID string  `json:"mata_pelajaran_id"`
		IDGuru          string  `json:"id_guru"`
		KelasID         string  `json:"kelas_id"`
		Hari            int     `json:"hari"`
		JamKeMulai      int     `json:"jam_ke_mulai"`
		JamKeSelesai    int     `json:"jam_ke_selesai"`
		JamMulai        string  `json:"jam_mulai"`
		JamSelesai      string  `json:"jam_selesai"`
		Ruang           string  `json:"ruang"`
		JenisKegiatan   string  `json:"jenis_kegiatan"`
		TahunAjaran     string  `json:"tahun_ajaran"`
		Semester        string  `json:"semester"`
		Keterangan      string  `json:"keterangan"`
		Aktif           *bool   `json:"aktif"`
	}

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			helper.JSONResponse(w, http.StatusBadRequest, "Format JSON tidak valid: "+err.Error())
			return err
		}
	} else {
		_ = r.ParseMultipartForm(10 << 20)
		if r.FormValue("school_unit_id") != "" {
			su := r.FormValue("school_unit_id")
			payload.SchoolUnitID = &su
		}
		payload.KodeJadwal = r.FormValue("kode_jadwal")
		payload.MataPelajaranID = r.FormValue("mata_pelajaran_id")
		payload.IDGuru = r.FormValue("id_guru")
		payload.KelasID = r.FormValue("kelas_id")
		h, _ := strconv.Atoi(r.FormValue("hari"))
		payload.Hari = h
		jkm, _ := strconv.Atoi(r.FormValue("jam_ke_mulai"))
		payload.JamKeMulai = jkm
		jks, _ := strconv.Atoi(r.FormValue("jam_ke_selesai"))
		payload.JamKeSelesai = jks
		payload.JamMulai = r.FormValue("jam_mulai")
		payload.JamSelesai = r.FormValue("jam_selesai")
		payload.Ruang = r.FormValue("ruang")
		payload.JenisKegiatan = r.FormValue("jenis_kegiatan")
		payload.TahunAjaran = r.FormValue("tahun_ajaran")
		payload.Semester = r.FormValue("semester")
		payload.Keterangan = r.FormValue("keterangan")
		if r.FormValue("aktif") != "" {
			val := r.FormValue("aktif") == "true" || r.FormValue("aktif") == "1"
			payload.Aktif = &val
		}
	}

	isAktif := true
	if payload.Aktif != nil {
		isAktif = *payload.Aktif
	}

	core := jadwal.JadwalCore{
		SchoolUnitID:    payload.SchoolUnitID,
		KodeJadwal:      payload.KodeJadwal,
		MataPelajaranID: payload.MataPelajaranID,
		IDGuru:          payload.IDGuru,
		KelasID:         payload.KelasID,
		Hari:            payload.Hari,
		JamKeMulai:      payload.JamKeMulai,
		JamKeSelesai:    payload.JamKeSelesai,
		JamMulai:        payload.JamMulai,
		JamSelesai:      payload.JamSelesai,
		Ruang:           payload.Ruang,
		JenisKegiatan:   payload.JenisKegiatan,
		TahunAjaran:     payload.TahunAjaran,
		Semester:        payload.Semester,
		Keterangan:      payload.Keterangan,
		Aktif:           isAktif,
	}

	if core.SchoolUnitID == nil && unitID != "" {
		core.SchoolUnitID = &unitID
	}

	created, err := c.service.Create(&core, role, userID)
	if err != nil {
		helper.JSONResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := helper.APIResponse(http.StatusCreated, "Jadwal pelajaran berhasil ditambahkan", created)
	return json.NewEncoder(w).Encode(res)
}

func (c *JadwalController) Update(w http.ResponseWriter, r *http.Request) error {
	role, userID := getAuthMeta(r)
	_, _, unitID, _ := helper.ExtractUserContext(r)
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter id wajib diisi")
		return errors.New("id empty")
	}

	var payload struct {
		KodeJadwal      string `json:"kode_jadwal"`
		MataPelajaranID string `json:"mata_pelajaran_id"`
		IDGuru          string `json:"id_guru"`
		KelasID         string `json:"kelas_id"`
		Hari            int    `json:"hari"`
		JamKeMulai      int    `json:"jam_ke_mulai"`
		JamKeSelesai    int    `json:"jam_ke_selesai"`
		JamMulai        string `json:"jam_mulai"`
		JamSelesai      string `json:"jam_selesai"`
		Ruang           string `json:"ruang"`
		JenisKegiatan   string `json:"jenis_kegiatan"`
		TahunAjaran     string `json:"tahun_ajaran"`
		Semester        string `json:"semester"`
		Keterangan      string `json:"keterangan"`
		Aktif           *bool  `json:"aktif"`
	}

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			helper.JSONResponse(w, http.StatusBadRequest, "Format JSON tidak valid: "+err.Error())
			return err
		}
	} else {
		_ = r.ParseMultipartForm(10 << 20)
		payload.KodeJadwal = r.FormValue("kode_jadwal")
		payload.MataPelajaranID = r.FormValue("mata_pelajaran_id")
		payload.IDGuru = r.FormValue("id_guru")
		payload.KelasID = r.FormValue("kelas_id")
		h, _ := strconv.Atoi(r.FormValue("hari"))
		payload.Hari = h
		jkm, _ := strconv.Atoi(r.FormValue("jam_ke_mulai"))
		payload.JamKeMulai = jkm
		jks, _ := strconv.Atoi(r.FormValue("jam_ke_selesai"))
		payload.JamKeSelesai = jks
		payload.JamMulai = r.FormValue("jam_mulai")
		payload.JamSelesai = r.FormValue("jam_selesai")
		payload.Ruang = r.FormValue("ruang")
		payload.JenisKegiatan = r.FormValue("jenis_kegiatan")
		payload.TahunAjaran = r.FormValue("tahun_ajaran")
		payload.Semester = r.FormValue("semester")
		payload.Keterangan = r.FormValue("keterangan")
		if r.FormValue("aktif") != "" {
			val := r.FormValue("aktif") == "true" || r.FormValue("aktif") == "1"
			payload.Aktif = &val
		}
	}

	isAktif := true
	if payload.Aktif != nil {
		isAktif = *payload.Aktif
	}

	core := jadwal.JadwalCore{
		KodeJadwal:      payload.KodeJadwal,
		MataPelajaranID: payload.MataPelajaranID,
		IDGuru:          payload.IDGuru,
		KelasID:         payload.KelasID,
		Hari:            payload.Hari,
		JamKeMulai:      payload.JamKeMulai,
		JamKeSelesai:    payload.JamKeSelesai,
		JamMulai:        payload.JamMulai,
		JamSelesai:      payload.JamSelesai,
		Ruang:           payload.Ruang,
		JenisKegiatan:   payload.JenisKegiatan,
		TahunAjaran:     payload.TahunAjaran,
		Semester:        payload.Semester,
		Keterangan:      payload.Keterangan,
		Aktif:           isAktif,
	}

	var updated *jadwal.JadwalCore
	var err error
	if unitID != "" {
		updated, err = c.service.Update(id, &core, role, userID, unitID)
	} else {
		updated, err = c.service.Update(id, &core, role, userID)
	}
	if err != nil {
		helper.JSONResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Jadwal pelajaran berhasil diperbarui", updated)
	return json.NewEncoder(w).Encode(res)
}

func (c *JadwalController) Delete(w http.ResponseWriter, r *http.Request) error {
	role, userID := getAuthMeta(r)
	_, _, unitID, _ := helper.ExtractUserContext(r)
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter id wajib diisi")
		return errors.New("id empty")
	}

	var err error
	if unitID != "" {
		err = c.service.Delete(id, role, userID, unitID)
	} else {
		err = c.service.Delete(id, role, userID)
	}
	if err != nil {
		helper.JSONResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	helper.JSONResponse(w, http.StatusOK, "Jadwal pelajaran berhasil dihapus")
	return nil
}
