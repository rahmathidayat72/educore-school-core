package controllers

import (
	"encoding/json"
	"errors"
	academicyears "go_rest_native_sekolah/features/academic_years"
	"go_rest_native_sekolah/helper"
	"net/http"
	"strings"
)

type AcademicYearController struct {
	service academicyears.ServiceAcademicYearInterface
}

func NewAcademicYearController(service academicyears.ServiceAcademicYearInterface) *AcademicYearController {
	if service == nil {
		panic("nil academic year service")
	}
	return &AcademicYearController{service: service}
}

// Academic Years Handlers

func (c *AcademicYearController) GetAllYears(w http.ResponseWriter, r *http.Request) error {
	_, _, unitID, _ := helper.ExtractUserContext(r)

	data, err := c.service.GetAllYears(unitID)
	if err != nil {
		helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Berhasil mengambil daftar tahun ajaran", ToYearResponseList(data))
	return json.NewEncoder(w).Encode(res)
}

func (c *AcademicYearController) GetYearByID(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter id wajib diisi")
		return errors.New("id empty")
	}

	_, _, unitID, isSuperAdmin := helper.ExtractUserContext(r)
	var data *academicyears.AcademicYearCore
	var err error
	if isSuperAdmin || unitID == "" {
		data, err = c.service.GetYearByID(id)
	} else {
		data, err = c.service.GetYearByID(id, unitID)
	}

	if err != nil {
		helper.JSONResponse(w, http.StatusNotFound, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Berhasil mengambil rincian tahun ajaran", ToYearResponse(*data))
	return json.NewEncoder(w).Encode(res)
}

func (c *AcademicYearController) GetActiveYear(w http.ResponseWriter, r *http.Request) error {
	_, _, unitID, _ := helper.ExtractUserContext(r)

	data, err := c.service.GetActiveYear(unitID)
	if err != nil {
		helper.JSONResponse(w, http.StatusNotFound, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Berhasil mengambil tahun ajaran aktif", ToYearResponse(*data))
	return json.NewEncoder(w).Encode(res)
}

func (c *AcademicYearController) CreateYear(w http.ResponseWriter, r *http.Request) error {
	_, _, unitID, isSuperAdmin := helper.ExtractUserContext(r)

	var req AcademicYearRequest
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.JSONResponse(w, http.StatusBadRequest, "Format JSON tidak valid: "+err.Error())
			return err
		}
	} else {
		_ = r.ParseMultipartForm(10 << 20)
		req.TahunAjaran = r.FormValue("tahun_ajaran")
		req.Semester = r.FormValue("semester")
		req.IsActive = r.FormValue("is_active") == "true"
		req.TanggalMulai = r.FormValue("tanggal_mulai")
		req.TanggalSelesai = r.FormValue("tanggal_selesai")
	}

	targetUnitID := unitID
	if isSuperAdmin {
		paramUnit := r.URL.Query().Get("unit_id")
		if paramUnit != "" {
			targetUnitID = paramUnit
		}
	}

	if targetUnitID == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Unit sekolah (school_unit_id) tidak terdeteksi. Sertakan ?unit_id= untuk Super Admin")
		return errors.New("unit id empty")
	}

	inputCore := ToYearCore(req, targetUnitID)
	created, err := c.service.CreateYear(&inputCore)
	if err != nil {
		helper.JSONResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := helper.APIResponse(http.StatusCreated, "Berhasil membuat tahun ajaran baru", ToYearResponse(*created))
	return json.NewEncoder(w).Encode(res)
}

func (c *AcademicYearController) UpdateYear(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter id wajib diisi")
		return errors.New("id empty")
	}

	_, _, unitID, isSuperAdmin := helper.ExtractUserContext(r)

	var req AcademicYearRequest
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.JSONResponse(w, http.StatusBadRequest, "Format JSON tidak valid: "+err.Error())
			return err
		}
	} else {
		_ = r.ParseMultipartForm(10 << 20)
		req.TahunAjaran = r.FormValue("tahun_ajaran")
		req.Semester = r.FormValue("semester")
		req.TanggalMulai = r.FormValue("tanggal_mulai")
		req.TanggalSelesai = r.FormValue("tanggal_selesai")
	}

	inputCore := ToYearCore(req, unitID)
	var updated *academicyears.AcademicYearCore
	var err error
	if isSuperAdmin || unitID == "" {
		updated, err = c.service.UpdateYear(id, &inputCore)
	} else {
		updated, err = c.service.UpdateYear(id, &inputCore, unitID)
	}

	if err != nil {
		helper.JSONResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Berhasil memperbarui data tahun ajaran", ToYearResponse(*updated))
	return json.NewEncoder(w).Encode(res)
}

func (c *AcademicYearController) SetActiveYear(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter id wajib diisi")
		return errors.New("id empty")
	}

	_, _, unitID, isSuperAdmin := helper.ExtractUserContext(r)
	targetUnitID := unitID
	if isSuperAdmin {
		paramUnit := r.URL.Query().Get("unit_id")
		if paramUnit != "" {
			targetUnitID = paramUnit
		}
	}

	if targetUnitID == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Unit sekolah wajib ditentukan untuk aktivasi tahun ajaran")
		return errors.New("unit id empty")
	}

	if err := c.service.SetActiveYear(id, targetUnitID); err != nil {
		helper.JSONResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	helper.JSONResponse(w, http.StatusOK, "Berhasil mengaktifkan tahun ajaran untuk unit ini")
	return nil
}

func (c *AcademicYearController) DeleteYear(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter id wajib diisi")
		return errors.New("id empty")
	}

	_, _, unitID, isSuperAdmin := helper.ExtractUserContext(r)
	var err error
	if isSuperAdmin || unitID == "" {
		err = c.service.DeleteYear(id)
	} else {
		err = c.service.DeleteYear(id, unitID)
	}

	if err != nil {
		helper.JSONResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	helper.JSONResponse(w, http.StatusOK, "Berhasil menghapus tahun ajaran")
	return nil
}

// Enrollments Handlers

func (c *AcademicYearController) GetEnrollments(w http.ResponseWriter, r *http.Request) error {
	_, _, unitID, _ := helper.ExtractUserContext(r)
	academicYearID := r.URL.Query().Get("academic_year_id")
	kelasID := r.URL.Query().Get("kelas_id")

	data, err := c.service.GetEnrollments(unitID, academicYearID, kelasID)
	if err != nil {
		helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Berhasil mengambil data penempatan kelas", ToEnrollmentResponseList(data))
	return json.NewEncoder(w).Encode(res)
}

func (c *AcademicYearController) GetEnrollmentByID(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter id wajib diisi")
		return errors.New("id empty")
	}

	_, _, unitID, isSuperAdmin := helper.ExtractUserContext(r)
	var data *academicyears.ClassEnrollmentCore
	var err error
	if isSuperAdmin || unitID == "" {
		data, err = c.service.GetEnrollmentByID(id)
	} else {
		data, err = c.service.GetEnrollmentByID(id, unitID)
	}

	if err != nil {
		helper.JSONResponse(w, http.StatusNotFound, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Berhasil mengambil rincian penempatan kelas", ToEnrollmentResponse(*data))
	return json.NewEncoder(w).Encode(res)
}

func (c *AcademicYearController) EnrollStudent(w http.ResponseWriter, r *http.Request) error {
	_, _, unitID, isSuperAdmin := helper.ExtractUserContext(r)

	var req ClassEnrollmentRequest
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.JSONResponse(w, http.StatusBadRequest, "Format JSON tidak valid: "+err.Error())
			return err
		}
	} else {
		_ = r.ParseMultipartForm(10 << 20)
		req.AcademicYearID = r.FormValue("academic_year_id")
		req.KelasID = r.FormValue("kelas_id")
		req.SiswaID = r.FormValue("siswa_id")
		req.Status = r.FormValue("status")
		req.Keterangan = r.FormValue("keterangan")
	}

	targetUnitID := unitID
	if isSuperAdmin {
		paramUnit := r.URL.Query().Get("unit_id")
		if paramUnit != "" {
			targetUnitID = paramUnit
		}
	}

	if targetUnitID == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Unit sekolah tidak terdeteksi")
		return errors.New("unit id empty")
	}

	inputCore := ToEnrollmentCore(req, targetUnitID)
	created, err := c.service.EnrollStudent(&inputCore)
	if err != nil {
		helper.JSONResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := helper.APIResponse(http.StatusCreated, "Berhasil mendaftarkan siswa ke kelas", ToEnrollmentResponse(*created))
	return json.NewEncoder(w).Encode(res)
}

func (c *AcademicYearController) UpdateEnrollment(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter id wajib diisi")
		return errors.New("id empty")
	}

	_, _, unitID, isSuperAdmin := helper.ExtractUserContext(r)

	var req ClassEnrollmentRequest
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.JSONResponse(w, http.StatusBadRequest, "Format JSON tidak valid: "+err.Error())
			return err
		}
	} else {
		_ = r.ParseMultipartForm(10 << 20)
		req.KelasID = r.FormValue("kelas_id")
		req.Status = r.FormValue("status")
		req.Keterangan = r.FormValue("keterangan")
	}

	inputCore := ToEnrollmentCore(req, unitID)
	var updated *academicyears.ClassEnrollmentCore
	var err error
	if isSuperAdmin || unitID == "" {
		updated, err = c.service.UpdateEnrollment(id, &inputCore)
	} else {
		updated, err = c.service.UpdateEnrollment(id, &inputCore, unitID)
	}

	if err != nil {
		helper.JSONResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Berhasil memperbarui penempatan kelas siswa", ToEnrollmentResponse(*updated))
	return json.NewEncoder(w).Encode(res)
}

func (c *AcademicYearController) DeleteEnrollment(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter id wajib diisi")
		return errors.New("id empty")
	}

	_, _, unitID, isSuperAdmin := helper.ExtractUserContext(r)
	var err error
	if isSuperAdmin || unitID == "" {
		err = c.service.DeleteEnrollment(id)
	} else {
		err = c.service.DeleteEnrollment(id, unitID)
	}

	if err != nil {
		helper.JSONResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	helper.JSONResponse(w, http.StatusOK, "Berhasil menghapus penempatan kelas siswa")
	return nil
}

func (c *AcademicYearController) GetStudentEnrollmentHistory(w http.ResponseWriter, r *http.Request) error {
	siswaID := r.URL.Query().Get("siswa_id")
	if strings.TrimSpace(siswaID) == "" {
		helper.JSONResponse(w, http.StatusBadRequest, "Parameter siswa_id wajib diisi")
		return errors.New("siswa_id empty")
	}

	data, err := c.service.GetStudentEnrollmentHistory(siswaID)
	if err != nil {
		helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	res := helper.APIResponse(http.StatusOK, "Berhasil mengambil riwayat kelas siswa", ToEnrollmentResponseList(data))
	return json.NewEncoder(w).Encode(res)
}
