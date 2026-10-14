package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/guru"
	"go_rest_native_sekolah/helper"
	"log"
	"net/http"
	"strings"
)

type Gurucontroller struct {
	guruService guru.ServiceGuruInterface
}

func NewGuruController(service guru.ServiceGuruInterface) *Gurucontroller {
	return &Gurucontroller{guruService: service}
}

func (gc *Gurucontroller) Guru(w http.ResponseWriter, r *http.Request) error {
	gurus, err := gc.guruService.GetAllGuru()
	if err != nil {
		return fmt.Errorf("guru controller: Error retrieving data: %v", err)
	}

	formattedGurus := FormatGuruList(gurus)
	response := helper.APIResponse(http.StatusOK, "Berhasil mengambil data guru dari database", formattedGurus)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func (gc *Gurucontroller) InsertGuru(w http.ResponseWriter, r *http.Request) error {
	if gc == nil || gc.guruService == nil {
		return fmt.Errorf("guru controller: service is nil")
	}

	var guruReq GuruFormatter
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") {
		err := json.NewDecoder(r.Body).Decode(&guruReq)
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
		guruReq = GuruFormatter{
			ID_User:            r.FormValue("id_user"),
			Nama:               r.FormValue("nama"),
			NIP:                r.FormValue("nip"),
			NUPTK:              r.FormValue("nuptk"),
			JenisKelamin:       r.FormValue("jenis_kelamin"),
			TempatLahir:        r.FormValue("tempat_lahir"),
			TanggalLahir:       r.FormValue("tanggal_lahir"),
			NoTelepon:          r.FormValue("no_telepon"),
			Agama:              r.FormValue("agama"),
			PendidikanTerakhir: r.FormValue("pendidikan_terakhir"),
			StatusKepegawaian:  r.FormValue("status_kepegawaian"),
			Email:              r.FormValue("email"),
			Alamat:             r.FormValue("alamat"),
		}
	}

	if strings.TrimSpace(guruReq.Nama) == "" || strings.TrimSpace(guruReq.Email) == "" || strings.TrimSpace(guruReq.Alamat) == "" {
		http.Error(w, "nama, email, dan alamat wajib diisi", http.StatusBadRequest)
		return fmt.Errorf("validasi gagal: nama, email, dan alamat harus diisi")
	}

	guruCore := FormatGuruRequestToCore(guruReq)

	err := gc.guruService.InsertGuru(&guruCore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return fmt.Errorf("gagal insert guru: %v", err)
	}

	formattedGurus := FormatGuruList([]guru.GuruCore{guruCore})
	response := helper.APIResponse(http.StatusCreated, "Berhasil menginsert data guru ke database", formattedGurus)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(response)
}

func (gc *Gurucontroller) UpdateGuru(w http.ResponseWriter, r *http.Request) error {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		return errors.New("missing 'id' query parameter")
	}
	log.Printf("Request update guru dengan ID: %s", idStr)

	var guruReq GuruFormatter
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") {
		err := json.NewDecoder(r.Body).Decode(&guruReq)
		if err != nil {
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, "Gagal memproses data input", http.StatusBadRequest)
			return err
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "gagal membaca form data", http.StatusBadRequest)
			return err
		}
		guruReq = GuruFormatter{
			Nama:               r.FormValue("nama"),
			NIP:                r.FormValue("nip"),
			NUPTK:              r.FormValue("nuptk"),
			JenisKelamin:       r.FormValue("jenis_kelamin"),
			TempatLahir:        r.FormValue("tempat_lahir"),
			TanggalLahir:       r.FormValue("tanggal_lahir"),
			NoTelepon:          r.FormValue("no_telepon"),
			Agama:              r.FormValue("agama"),
			PendidikanTerakhir: r.FormValue("pendidikan_terakhir"),
			StatusKepegawaian:  r.FormValue("status_kepegawaian"),
			Email:              r.FormValue("email"),
			Alamat:             r.FormValue("alamat"),
		}
	}

	updateGuru := FormatGuruRequestToCore(guruReq)

	err := gc.guruService.UpdateGuru(&updateGuru, idStr)
	if err != nil {
		if strings.Contains(err.Error(), "validation") || strings.Contains(err.Error(), "sudah terdaftar") || strings.Contains(err.Error(), "format") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
		if strings.Contains(err.Error(), "tidak ditemukan") {
			http.Error(w, "ID guru tidak ditemukan atau tidak ada perubahan data", http.StatusNotFound)
			return err
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	updatedGuru, err := gc.guruService.SelectById(idStr)
	if err != nil {
		http.Error(w, "Gagal mengambil data setelah update", http.StatusInternalServerError)
		return err
	}

	formattedGurus := FormatGuruList([]guru.GuruCore{*updatedGuru})
	response := helper.APIResponse(http.StatusOK, "Berhasil mengupdate data guru ke database", formattedGurus)

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func (gc *Gurucontroller) GetGuruById(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "parameter 'id' wajib diisi", http.StatusBadRequest)
		return fmt.Errorf("guru controller: ID guru tidak ditemukan dalam query parameter")
	}

	guruData, err := gc.guruService.SelectById(id)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			http.Error(w, "Data guru tidak ditemukan", http.StatusNotFound)
			return nil
		}
		http.Error(w, "Terjadi kesalahan saat mengambil data guru", http.StatusInternalServerError)
		return fmt.Errorf("guru controller: gagal mengambil data guru berdasarkan ID: %v", err)
	}

	formattedGurus := FormatGuruList([]guru.GuruCore{*guruData})
	response := helper.APIResponse(http.StatusOK, "Berhasil mengambil data guru berdasarkan ID", formattedGurus)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func (gc *Gurucontroller) DeleteGuru(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID guru tidak ditemukan dalam query parameter", http.StatusBadRequest)
		return fmt.Errorf("guru controller: ID guru tidak ditemukan dalam query parameter")
	}

	err := gc.guruService.DeleteById(id)
	if err != nil {
		return fmt.Errorf("guru controller: gagal menghapus data guru berdasarkan ID: %v", err)
	}

	response := helper.APIResponse(http.StatusOK, "Berhasil menghapus data guru berdasarkan ID", "Terhapus")
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
