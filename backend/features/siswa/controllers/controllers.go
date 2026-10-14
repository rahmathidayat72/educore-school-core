package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/siswa"
	"go_rest_native_sekolah/helper"
	"log"
	"net/http"
	"strings"
)

type SiswaController struct {
	SiswaService siswa.ServiceSiswaInterface
}

func NewSiswaController(service siswa.ServiceSiswaInterface) *SiswaController {
	return &SiswaController{
		SiswaService: service,
	}
}

func (sc *SiswaController) InsertSiswa(w http.ResponseWriter, r *http.Request) error {
	if sc == nil || sc.SiswaService == nil {
		return errors.New("Nil controller")
	}

	var siswaReq SiswaFormatter
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") {
		err := json.NewDecoder(r.Body).Decode(&siswaReq)
		if err != nil {
			http.Error(w, "gagal membaca JSON", http.StatusBadRequest)
			return fmt.Errorf("error decoding JSON: %v", err)
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "gagal membaca form data", http.StatusBadRequest)
			return fmt.Errorf("error parsing form-data: %v", err)
		}
		siswaReq = SiswaFormatter{
			Nama:          r.FormValue("nama"),
			Kelas_ID:      r.FormValue("kelas_id"),
			Nama_Kelas:    r.FormValue("nama_kelas"),
			NIS:           r.FormValue("nis"),
			NISN:          r.FormValue("nisn"),
			JenisKelamin:  r.FormValue("jenis_kelamin"),
			TempatLahir:   r.FormValue("tempat_lahir"),
			TanggalLahir:  r.FormValue("tanggal_lahir"),
			Agama:         r.FormValue("agama"),
			NoTelepon:     r.FormValue("no_telepon"),
			NamaWali:      r.FormValue("nama_wali"),
			NoTeleponWali: r.FormValue("no_telepon_wali"),
			TahunMasuk:    r.FormValue("tahun_masuk"),
			Status:        r.FormValue("status"),
			Email:         r.FormValue("email"),
			Alamat:        r.FormValue("alamat"),
		}
	}

	siswaCore := FormatSiswaRequestToCore(siswaReq)

	_, _, unitID, _ := helper.ExtractUserContext(r)
	if siswaCore.SchoolUnitID == nil && unitID != "" {
		siswaCore.SchoolUnitID = &unitID
	}

	err := sc.SiswaService.InsertSiswa(&siswaCore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	formatKelas := FormatterKelasList([]siswa.SiswaCore{siswaCore})
	respon := helper.APIResponse(http.StatusCreated, "Success insert data siswa", formatKelas)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(respon)
}

func (sc *SiswaController) Siswa(w http.ResponseWriter, r *http.Request) error {
	if sc == nil || sc.SiswaService == nil {
		return errors.New("Nil controller")
	}

	_, _, unitID, _ := helper.ExtractUserContext(r)
	siswaList, err := sc.SiswaService.SelectAllSiswa(unitID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	formatKelas := FormatterKelasList(siswaList)
	respon := helper.APIResponse(http.StatusOK, "Success get data siswa", formatKelas)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(respon)
}

func (sc *SiswaController) GetSiswaById(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		return errors.New("missing 'id' query parameter")
	}

	_, _, unitID, _ := helper.ExtractUserContext(r)
	siswaData, err := sc.SiswaService.SelectById(id, unitID)
	if err != nil {
		http.Error(w, "Data siswa tidak ditemukan", http.StatusNotFound)
		return err
	}

	formatKelas := FormatterKelasList([]siswa.SiswaCore{*siswaData})
	respon := helper.APIResponse(http.StatusOK, "Success get data siswaById", formatKelas)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(respon)
}

func (sc *SiswaController) UpdateSiswa(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		return errors.New("missing 'id' query parameter")
	}

	var siswaReq SiswaFormatter
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") {
		err := json.NewDecoder(r.Body).Decode(&siswaReq)
		if err != nil {
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, "gagal memproses data input", http.StatusBadRequest)
			return err
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "gagal membaca form data", http.StatusBadRequest)
			return err
		}
		siswaReq = SiswaFormatter{
			Nama:          r.FormValue("nama"),
			Kelas_ID:      r.FormValue("kelas_id"),
			Nama_Kelas:    r.FormValue("nama_kelas"),
			NIS:           r.FormValue("nis"),
			NISN:          r.FormValue("nisn"),
			JenisKelamin:  r.FormValue("jenis_kelamin"),
			TempatLahir:   r.FormValue("tempat_lahir"),
			TanggalLahir:  r.FormValue("tanggal_lahir"),
			Agama:         r.FormValue("agama"),
			NoTelepon:     r.FormValue("no_telepon"),
			NamaWali:      r.FormValue("nama_wali"),
			NoTeleponWali: r.FormValue("no_telepon_wali"),
			TahunMasuk:    r.FormValue("tahun_masuk"),
			Status:        r.FormValue("status"),
			Email:         r.FormValue("email"),
			Alamat:        r.FormValue("alamat"),
		}
	}

	siswaUpdate := FormatSiswaRequestToCore(siswaReq)
	_, _, unitID, _ := helper.ExtractUserContext(r)
	err := sc.SiswaService.Update(&siswaUpdate, id, unitID)
	if err != nil {
		if strings.Contains(err.Error(), "validation") || strings.Contains(err.Error(), "sudah terdaftar") || strings.Contains(err.Error(), "format") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	siswaData, err := sc.SiswaService.SelectById(id, unitID)
	if err != nil {
		http.Error(w, "Gagal mengambil data setelah update", http.StatusInternalServerError)
		return err
	}

	formatKelas := FormatterKelasList([]siswa.SiswaCore{*siswaData})
	respon := helper.APIResponse(http.StatusOK, "Berhasil mengupdate data siswa", formatKelas)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(respon)
}

func (sc *SiswaController) DeleteSiswa(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		return errors.New("missing 'id' query parameter")
	}

	_, _, unitID, _ := helper.ExtractUserContext(r)
	err := sc.SiswaService.DeleteById(id, unitID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	respon := helper.APIResponse(http.StatusOK, "Berhasil menghapus data siswa", nil)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(respon)
}
