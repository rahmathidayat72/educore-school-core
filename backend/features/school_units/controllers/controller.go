package controllers

import (
	"encoding/json"
	"errors"
	schoolunits "go_rest_native_sekolah/features/school_units"
	"go_rest_native_sekolah/helper"
	"net/http"
	"strings"
)

type SchoolUnitController struct {
	service schoolunits.ServiceSchoolUnitInterface
}

func NewSchoolUnitController(service schoolunits.ServiceSchoolUnitInterface) *SchoolUnitController {
	if service == nil {
		panic("nil school unit service")
	}
	return &SchoolUnitController{service: service}
}

func (c *SchoolUnitController) GetAll(w http.ResponseWriter, r *http.Request) error {
	jenjang := r.URL.Query().Get("jenjang")
	status := r.URL.Query().Get("status")

	data, err := c.service.GetAll(jenjang, status)
	if err != nil {
		helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Berhasil mengambil data unit sekolah", ToResponseList(data))
	return json.NewEncoder(w).Encode(res)
}

func (c *SchoolUnitController) GetByID(w http.ResponseWriter, r *http.Request) error {
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
	res := helper.APIResponse(http.StatusOK, "Berhasil mengambil rincian unit sekolah", ToResponse(*data))
	return json.NewEncoder(w).Encode(res)
}

func (c *SchoolUnitController) Create(w http.ResponseWriter, r *http.Request) error {
	var req SchoolUnitRequest

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.JSONResponse(w, http.StatusBadRequest, "Format JSON tidak valid: "+err.Error())
			return err
		}
	} else {
		_ = r.ParseMultipartForm(10 << 20)
		req.NamaSekolah = r.FormValue("nama_sekolah")
		req.Jenjang = r.FormValue("jenjang")
		req.NPSN = r.FormValue("npsn")
		req.Alamat = r.FormValue("alamat")
		req.Telepon = r.FormValue("telepon")
		req.Email = r.FormValue("email")
		req.Website = r.FormValue("website")
		req.Status = r.FormValue("status")
	}

	core := ToCore(req)
	created, err := c.service.Create(&core)
	if err != nil {
		helper.JSONResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := helper.APIResponse(http.StatusCreated, "Unit sekolah berhasil ditambahkan", ToResponse(*created))
	return json.NewEncoder(w).Encode(res)
}

func (c *SchoolUnitController) Update(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter id wajib diisi")
		return errors.New("id empty")
	}

	var req SchoolUnitRequest
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.JSONResponse(w, http.StatusBadRequest, "Format JSON tidak valid: "+err.Error())
			return err
		}
	} else {
		_ = r.ParseMultipartForm(10 << 20)
		req.NamaSekolah = r.FormValue("nama_sekolah")
		req.Jenjang = r.FormValue("jenjang")
		req.NPSN = r.FormValue("npsn")
		req.Alamat = r.FormValue("alamat")
		req.Telepon = r.FormValue("telepon")
		req.Email = r.FormValue("email")
		req.Website = r.FormValue("website")
		req.Status = r.FormValue("status")
	}

	core := ToCore(req)
	updated, err := c.service.Update(id, &core)
	if err != nil {
		helper.JSONResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Unit sekolah berhasil diperbarui", ToResponse(*updated))
	return json.NewEncoder(w).Encode(res)
}

func (c *SchoolUnitController) Delete(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter id wajib diisi")
		return errors.New("id empty")
	}

	if err := c.service.Delete(id); err != nil {
		helper.JSONResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	helper.JSONResponse(w, http.StatusOK, "Unit sekolah berhasil dihapus")
	return nil
}
