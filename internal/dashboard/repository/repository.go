package repository

import (
	"context"

	"sensor-data-service.backend/internal/dashboard/domain"
)

type DashboardRepository interface {
	FindByID(ctx context.Context, id int32) (*domain.Dashboard, error)
	List(ctx context.Context) ([]*domain.Dashboard, error)
	Save(ctx context.Context, d *domain.Dashboard) error
	Delete(ctx context.Context, id int32) error
	Update(ctx context.Context, d *domain.Dashboard) error
	FindByName(ctx context.Context, name string) (*domain.Dashboard, error)
}
