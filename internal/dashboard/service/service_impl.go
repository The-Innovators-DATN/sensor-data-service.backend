package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"sensor-data-service.backend/internal/dashboard/domain"
	"sensor-data-service.backend/internal/dashboard/repository"
)

type dashboardServiceImpl struct {
	repo repository.DashboardRepository
}

func NewDashboardService(repo repository.DashboardRepository) *dashboardServiceImpl {
	return &dashboardServiceImpl{repo: repo}
}

func (s *dashboardServiceImpl) GetDashboardByID(ctx context.Context, uid string, userID int32) (*domain.Dashboard, error) {
	d, err := s.repo.FindByID(ctx, uid, userID)
	if err != nil {
		return nil, fmt.Errorf("GetDashboardByID: %w", err)
	}
	if d.CreatedBy != userID {
		return nil, fmt.Errorf("unauthorized access to dashboard %s", uid)
	}
	return d, nil
}

func (s *dashboardServiceImpl) ListDashboardsByUser(ctx context.Context, userID int32) ([]*domain.Dashboard, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *dashboardServiceImpl) ListDashboards(ctx context.Context) ([]*domain.Dashboard, error) {
	return s.repo.ListAll(ctx)
}

func (s *dashboardServiceImpl) CreateDashboard(ctx context.Context, d *domain.Dashboard) error {
	// Sinh UID nếu chưa có
	if d.UID == uuid.Nil {
		log.Printf("[info] CreateDashboard: generating new UID")
		d.UID = uuid.New()
	}
	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now
	if d.Version == 0 {
		d.Version = 1
	}
	if d.Status == "" {
		d.Status = "active"
	}
	log.Printf("[info] CreateDashboard: uid=%s by user_id=%d", d.UID, d.CreatedBy)
	return s.repo.Create(ctx, d)
}

func (s *dashboardServiceImpl) UpdateDashboard(ctx context.Context, d *domain.Dashboard, userID int32) error {

	origin, err := s.repo.FindByID(ctx, d.UID.String(), userID)
	if err != nil {
		return fmt.Errorf("UpdateDashboard: %w", err)
	}
	if origin.CreatedBy != userID {
		return fmt.Errorf("unauthorized update for dashboard %s", d.UID)
	}
	log.Printf("[info] UpdateDashboard: uid=%s user_id=%d", d.UID, userID)
	return s.repo.Update(ctx, d)
}

func (s *dashboardServiceImpl) PatchDashboard(ctx context.Context, d *domain.Dashboard, userID int32) error {
	origin, err := s.repo.FindByID(ctx, d.UID.String(), userID)
	if err != nil {
		return fmt.Errorf("PatchDashboard: %w", err)
	}
	if origin.CreatedBy != userID {
		return fmt.Errorf("unauthorized patch for dashboard %s", d.UID)
	}
	log.Printf("[info] PatchDashboard: uid=%s user_id=%d", d.UID, userID)
	return s.repo.Update(ctx, d)
}

func (s *dashboardServiceImpl) DeleteDashboard(ctx context.Context, uid string, userID int32) error {
	origin, err := s.repo.FindByID(ctx, uid, userID)
	if err != nil {
		return fmt.Errorf("DeleteDashboard: %w", err)
	}
	if origin.CreatedBy != userID {
		return fmt.Errorf("unauthorized delete for dashboard %s", uid)
	}
	log.Printf("[info] DeleteDashboard: uid=%s user_id=%d", uid, userID)
	return s.repo.Delete(ctx, uid)
}

func (s *dashboardServiceImpl) FindDashboardByName(ctx context.Context, name string, userID int32) (*domain.Dashboard, error) {
	d, err := s.repo.FindByNameAndUser(ctx, name, userID)
	if err != nil {
		return nil, fmt.Errorf("FindDashboardByName: %w", err)
	}
	if d.CreatedBy != userID {
		return nil, fmt.Errorf("unauthorized access to dashboard name=%s", name)
	}
	return d, nil
}
