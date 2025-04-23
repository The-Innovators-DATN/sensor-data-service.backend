package repository

import (
	"context"

	"sensor-data-service.backend/internal/domain/model"
)

type DashboardRepository interface {
	FindByID(ctx context.Context, uid string, userID int32) (*model.Dashboard, error)
	FindByIDAndUser(ctx context.Context, uid string, userID int32) (*model.Dashboard, error)
	FindByNameAndUser(ctx context.Context, name string, userID int32) (*model.Dashboard, error)
	ListByUser(ctx context.Context, userID int32) ([]*model.Dashboard, error)
	ListAll(ctx context.Context) ([]*model.Dashboard, error)
	Create(ctx context.Context, d *model.Dashboard) error
	Update(ctx context.Context, d *model.Dashboard) error
	Patch(ctx context.Context, d *model.Dashboard) error
	Delete(ctx context.Context, uid string) error
}
