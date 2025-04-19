package service

import (
	"context"

	"sensor-data-service.backend/internal/dashboard/domain"
)

type DashboardService interface {
	GetDashboardByID(ctx context.Context, id int32) (*domain.Dashboard, error)
	ListDashboards(ctx context.Context) ([]*domain.Dashboard, error)
	SaveDashboard(ctx context.Context, d *domain.Dashboard) error
	DeleteDashboard(ctx context.Context, id int32) error
	UpdateDashboard(ctx context.Context, d *domain.Dashboard) error
	FindDashboardByName(ctx context.Context, name string) (*domain.Dashboard, error)
}
