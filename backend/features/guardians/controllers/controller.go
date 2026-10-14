package controllers

import (
	"encoding/json"
	"errors"
	"go_rest_native_sekolah/features/guardians"
	"go_rest_native_sekolah/helper"
	"net/http"
	"strings"
)

type GuardianController struct {
	service guardians.ServiceGuardianInterface
}

func NewGuardianController(service guardians.ServiceGuardianInterface) *GuardianController {
	if service == nil {
		panic("nil guardian service")
	}
	return &GuardianController{service: service}
}

func (c *GuardianController) GetAll(w http.ResponseWriter, r *http.Request) error {
	search := r.URL.Query().Get("search")

	data, err := c.service.GetAll(search)
	if err != nil {
		helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Berhasil mengambil daftar wali murid", ToGuardianResponseList(data))
	return json.NewEncoder(w).Encode(res)
}

func (c *GuardianController) GetByID(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter id wajib diisi")
		return errors.New("id empty")
	}

	data, err := c.service.GetByID(id)
	if err != nil {
		helper.JSONResponse(w, http.StatusNotFound, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Berhasil mengambil rincian data wali murid", ToGuardianResponse(*data))
	return json.NewEncoder(w).Encode(res)
}

func (c *GuardianController) GetMyProfile(w http.ResponseWriter, r *http.Request) error {
	userID, role, _, _ := helper.ExtractUserContext(r)
	if strings.ToLower(role) != "wali" && strings.ToLower(role) != "super_admin" && strings.ToLower(role) != "admin" {
		helper.JSONResponse(w, http.StatusForbidden, "Hanya akun wali murid yang dapat mengakses endpoint profil ini")
		return errors.New("forbidden")
	}

	data, err := c.service.GetByUserID(userID)
	if err != nil {
		helper.JSONResponse(w, http.StatusNotFound, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Berhasil mengambil profil wali murid", ToGuardianResponse(*data))
	return json.NewEncoder(w).Encode(res)
}

func (c *GuardianController) Create(w http.ResponseWriter, r *http.Request) error {
	var req GuardianRequest
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.JSONResponse(w, http.StatusBadRequest, "Format JSON tidak valid: "+err.Error())
			return err
		}
	} else {
		_ = r.ParseMultipartForm(10 << 20)
		req.Nama = r.FormValue("nama")
		req.NIK = r.FormValue("nik")
		req.Hubungan = r.FormValue("hubungan")
		req.Pekerjaan = r.FormValue("pekerjaan")
		req.Penghasilan = r.FormValue("penghasilan")
		req.Telepon = r.FormValue("telepon")
		req.Email = r.FormValue("email")
		req.Alamat = r.FormValue("alamat")
		if uid := r.FormValue("user_id"); uid != "" {
			req.UserID = &uid
		}
	}

	inputCore := ToGuardianCore(req)
	created, err := c.service.Create(&inputCore)
	if err != nil {
		helper.JSONResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := helper.APIResponse(http.StatusCreated, "Berhasil mendaftarkan wali murid baru", ToGuardianResponse(*created))
	return json.NewEncoder(w).Encode(res)
}

func (c *GuardianController) Update(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter id wajib diisi")
		return errors.New("id empty")
	}

	var req GuardianRequest
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.JSONResponse(w, http.StatusBadRequest, "Format JSON tidak valid: "+err.Error())
			return err
		}
	} else {
		_ = r.ParseMultipartForm(10 << 20)
		req.Nama = r.FormValue("nama")
		req.NIK = r.FormValue("nik")
		req.Hubungan = r.FormValue("hubungan")
		req.Pekerjaan = r.FormValue("pekerjaan")
		req.Penghasilan = r.FormValue("penghasilan")
		req.Telepon = r.FormValue("telepon")
		req.Email = r.FormValue("email")
		req.Alamat = r.FormValue("alamat")
		if uid := r.FormValue("user_id"); uid != "" {
			req.UserID = &uid
		}
	}

	inputCore := ToGuardianCore(req)
	updated, err := c.service.Update(id, &inputCore)
	if err != nil {
		helper.JSONResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Berhasil memperbarui data wali murid", ToGuardianResponse(*updated))
	return json.NewEncoder(w).Encode(res)
}

func (c *GuardianController) Delete(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter id wajib diisi")
		return errors.New("id empty")
	}

	if err := c.service.Delete(id); err != nil {
		helper.JSONResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	helper.JSONResponse(w, http.StatusOK, "Berhasil menghapus data wali murid")
	return nil
}

// Student-Guardian Link Handlers

func (c *GuardianController) LinkStudent(w http.ResponseWriter, r *http.Request) error {
	var req LinkStudentRequest
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.JSONResponse(w, http.StatusBadRequest, "Format JSON tidak valid: "+err.Error())
			return err
		}
	} else {
		_ = r.ParseMultipartForm(10 << 20)
		req.GuardianID = r.FormValue("guardian_id")
		req.SiswaID = r.FormValue("siswa_id")
		req.IsPrimary = r.FormValue("is_primary") == "true"
	}

	linked, err := c.service.LinkStudent(req.GuardianID, req.SiswaID, req.IsPrimary)
	if err != nil {
		helper.JSONResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := helper.APIResponse(http.StatusCreated, "Berhasil menautkan siswa dengan wali murid", ToStudentGuardianResponse(*linked))
	return json.NewEncoder(w).Encode(res)
}

func (c *GuardianController) UnlinkStudent(w http.ResponseWriter, r *http.Request) error {
	guardianID := r.URL.Query().Get("guardian_id")
	siswaID := r.URL.Query().Get("siswa_id")

	if guardianID == "" || siswaID == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter guardian_id dan siswa_id wajib diisi")
		return errors.New("missing parameters")
	}

	if err := c.service.UnlinkStudent(guardianID, siswaID); err != nil {
		helper.JSONResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	helper.JSONResponse(w, http.StatusOK, "Berhasil memutuskan tautan siswa dan wali murid")
	return nil
}

func (c *GuardianController) GetStudentsByGuardianID(w http.ResponseWriter, r *http.Request) error {
	guardianID := r.URL.Query().Get("guardian_id")
	if strings.TrimSpace(guardianID) == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter guardian_id wajib diisi")
		return errors.New("guardian_id empty")
	}

	data, err := c.service.GetStudentsByGuardianID(guardianID)
	if err != nil {
		helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Berhasil mengambil daftar anak yang terhubung", ToStudentGuardianResponseList(data))
	return json.NewEncoder(w).Encode(res)
}

func (c *GuardianController) GetGuardiansBySiswaID(w http.ResponseWriter, r *http.Request) error {
	siswaID := r.URL.Query().Get("siswa_id")
	if strings.TrimSpace(siswaID) == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter siswa_id wajib diisi")
		return errors.New("siswa_id empty")
	}

	data, err := c.service.GetGuardiansBySiswaID(siswaID)
	if err != nil {
		helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Berhasil mengambil daftar wali siswa", ToStudentGuardianResponseList(data))
	return json.NewEncoder(w).Encode(res)
}
