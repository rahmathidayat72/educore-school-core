package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_rest_native_sekolah/features/logs"
	"go_rest_native_sekolah/helper"
	"log"
	"net/http"
	"strconv"
)

type LogController struct {
	logService logs.ServiceLogInterface
}

func NewLogController(logService logs.ServiceLogInterface) *LogController {
	if logService == nil {
		panic("NewLogController: logService is nil")
	}
	return &LogController{
		logService: logService,
	}
}

// GetLogs mengambil daftar log sistem dengan filter dan pagination
func (lc *LogController) GetLogs(w http.ResponseWriter, r *http.Request) error {
	if lc == nil || lc.logService == nil {
		http.Error(w, "Service is not available", http.StatusInternalServerError)
		return errors.New("log controller: service is nil")
	}

	query := r.URL.Query()
	search := query.Get("search")
	resultFilter := query.Get("result")
	moduleFilter := query.Get("module")
	activityTypeFilter := query.Get("activity_type")

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	noCount := query.Get("no_count") == "true"

	filter := logs.FilterLog{
		Search:       search,
		Result:       resultFilter,
		Module:       moduleFilter,
		ActivityType: activityTypeFilter,
		Limit:        limit,
		Offset:       offset,
		NoCount:      noCount,
	}

	dataLogs, total, err := lc.logService.GetAllLogs(filter)
	if err != nil {
		log.Printf("Error retrieving logs: %v", err)
		http.Error(w, "Gagal mengambil data log: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	formattedLogs := ToLogResponses(dataLogs)
	paginated := LogPaginationResponse{
		Total: total,
		Limit: limit,
		Page:  page,
		Data:  formattedLogs,
	}

	w.Header().Set("Content-Type", "application/json")
	apiResponse := helper.APIResponse(http.StatusOK, "success get all logs", paginated)
	if err := json.NewEncoder(w).Encode(apiResponse); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return err
	}

	return nil
}

// GetLogByID mengambil detail 1 log
func (lc *LogController) GetLogByID(w http.ResponseWriter, r *http.Request, id int) error {
	if lc == nil || lc.logService == nil {
		http.Error(w, "Service is not available", http.StatusInternalServerError)
		return errors.New("log controller: service is nil")
	}

	item, err := lc.logService.GetLogByID(id)
	if err != nil {
		http.Error(w, "Log tidak ditemukan", http.StatusNotFound)
		return fmt.Errorf("error retrieving log %d: %v", id, err)
	}

	w.Header().Set("Content-Type", "application/json")
	apiResponse := helper.APIResponse(http.StatusOK, "success get log detail", ToLogResponse(*item))
	if err := json.NewEncoder(w).Encode(apiResponse); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return err
	}

	return nil
}
