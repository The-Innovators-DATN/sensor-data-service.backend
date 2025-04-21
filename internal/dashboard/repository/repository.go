package repository

import (
	"context"

	"sensor-data-service.backend/internal/dashboard/domain"
)

type DashboardRepository interface {
	FindByID(ctx context.Context, uid string, userID int32) (*domain.Dashboard, error)
	FindByIDAndUser(ctx context.Context, uid string, userID int32) (*domain.Dashboard, error)
	FindByNameAndUser(ctx context.Context, name string, userID int32) (*domain.Dashboard, error)
	ListByUser(ctx context.Context, userID int32) ([]*domain.Dashboard, error)
	ListAll(ctx context.Context) ([]*domain.Dashboard, error)
	Create(ctx context.Context, d *domain.Dashboard) error
	Update(ctx context.Context, d *domain.Dashboard) error
	Patch(ctx context.Context, d *domain.Dashboard) error
	Delete(ctx context.Context, uid string) error
}
