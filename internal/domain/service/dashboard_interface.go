package service

import (
	"context"

	"sensor-data-service.backend/internal/domain/model"
)

type DashboardService interface {
	GetDashboardByID(ctx context.Context, uid string, userID int32) (*model.Dashboard, error)
	ListDashboardsByUser(ctx context.Context, userID int32, page, limit int32) ([]*model.Dashboard, error)
	CreateDashboard(ctx context.Context, d *model.Dashboard) (string, error)
	UpdateDashboard(ctx context.Context, d *model.Dashboard, userID int32) error
	DeleteDashboard(ctx context.Context, uid string, userID int32) error
	PatchDashboard(ctx context.Context, d *model.Dashboard, userID int32) error
	FindDashboardByName(ctx context.Context, name string, userID int32) (*model.Dashboard, error)
	ListDashboards(ctx context.Context) ([]*model.Dashboard, error)
}
