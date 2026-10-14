package controllers

import (
	"encoding/json"
	"errors"
	"go_rest_native_sekolah/features/dashboard"
	"go_rest_native_sekolah/helper"
	"net/http"
	"strings"
)

type DashboardController struct {
	dashboardService dashboard.ServiceDashboardInterface
}

func NewDashboardController(service dashboard.ServiceDashboardInterface) *DashboardController {
	return &DashboardController{dashboardService: service}
}

// GetSummary menangani request GET /dashboard/summary
func (dc *DashboardController) GetSummary(w http.ResponseWriter, r *http.Request) error {
	if dc == nil || dc.dashboardService == nil {
		http.Error(w, "Dashboard service not available", http.StatusInternalServerError)
		return errors.New("dashboard controller: service is nil")
	}

	authHeader := r.Header.Get("Authorization")
	tokenStr := helper.GetTokenFromAuthorizationHeader(authHeader)
	userRole := "admin"
	userID := ""

	if tokenStr != "" {
		if meta, err := helper.VerifyTokenHeader(tokenStr); err == nil {
			if meta.Role != "" {
				userRole = strings.ToLower(strings.TrimSpace(meta.Role))
			}
			userID = meta.ID
		}
	}

	summary, err := dc.dashboardService.GetSummary(userRole, userID)
	if err != nil {
		http.Error(w, "Gagal mengambil data ringkasan dasbor: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	apiResponse := helper.APIResponse(http.StatusOK, "success get dashboard summary", summary)
	return json.NewEncoder(w).Encode(apiResponse)
}
