package service

import (
	"context"
	"errors"

	"sensor-data-service.backend/internal/dashboard/domain"
	"sensor-data-service.backend/internal/dashboard/repository"
)

type dashboardService struct {
	repo repository.DashboardRepository
}

func NewDashboardService(repo repository.DashboardRepository) DashboardService {
	return &dashboardService{
		repo: repo,
	}
}

func (s *dashboardService) GetDashboardByID(ctx context.Context, id int32) (*domain.Dashboard, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *dashboardService) ListDashboards(ctx context.Context) ([]*domain.Dashboard, error) {
	return s.repo.List(ctx)
}

func (s *dashboardService) SaveDashboard(ctx context.Context, d *domain.Dashboard) error {
	if d == nil || d.Name == "" || d.LayoutJSON == "" {
		return errors.New("missing required fields")
	}
	return s.repo.Save(ctx, d)
}

func (s *dashboardService) DeleteDashboard(ctx context.Context, id int32) error {
	return s.repo.Delete(ctx, id)
}

func (s *dashboardService) UpdateDashboard(ctx context.Context, d *domain.Dashboard) error {
	if d == nil || d.ID == 0 {
		return errors.New("invalid dashboard ID")
	}
	return s.repo.Update(ctx, d)
}

func (s *dashboardService) FindDashboardByName(ctx context.Context, name string) (*domain.Dashboard, error) {
	if name == "" {
		return nil, errors.New("dashboard name is required")
	}
	return s.repo.FindByName(ctx, name)
}
