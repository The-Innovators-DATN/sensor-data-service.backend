package service

import (
	"context"

	"sensor-data-service.backend/internal/dashboard/domain"
)

type DashboardService interface {
	GetDashboardByID(ctx context.Context, uid string, userID int32) (*domain.Dashboard, error)
	ListDashboardsByUser(ctx context.Context, userID int32) ([]*domain.Dashboard, error)
	CreateDashboard(ctx context.Context, d *domain.Dashboard) error
	UpdateDashboard(ctx context.Context, d *domain.Dashboard, userID int32) error
	DeleteDashboard(ctx context.Context, uid string, userID int32) error
	PatchDashboard(ctx context.Context, d *domain.Dashboard, userID int32) error
	FindDashboardByName(ctx context.Context, name string, userID int32) (*domain.Dashboard, error)
	ListDashboards(ctx context.Context) ([]*domain.Dashboard, error)
}
