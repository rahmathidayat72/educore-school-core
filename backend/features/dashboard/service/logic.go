package service

import (
	"errors"
	"go_rest_native_sekolah/features/dashboard"
)

type DashboardService struct {
	dashboardRepo dashboard.DataDashboardInterface
}

func NewDashboardService(repo dashboard.DataDashboardInterface) dashboard.ServiceDashboardInterface {
	if repo == nil {
		panic("dashboard service: nil repo")
	}
	return &DashboardService{dashboardRepo: repo}
}

func (s *DashboardService) GetSummary(role string, userID string) (*dashboard.DashboardSummary, error) {
	if s.dashboardRepo == nil {
		return nil, errors.New("dashboard service: repository not initialized")
	}
	return s.dashboardRepo.GetSummary(role, userID)
}
