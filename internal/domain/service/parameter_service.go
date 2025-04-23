package service

import (
	"context"

	"sensor-data-service.backend/internal/domain/model"
	"sensor-data-service.backend/internal/domain/repository"
)

type ParameterService struct {
	repo repository.ParameterRepository
}

func NewParameterService(repo repository.ParameterRepository) *ParameterService {
	return &ParameterService{repo: repo}
}

func (s *ParameterService) ListParameters(ctx context.Context) ([]model.Parameter, error) {
	return s.repo.GetAll(ctx)
}
func (s *ParameterService) GetParameter(ctx context.Context, id int) (model.Parameter, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ParameterService) CreateParameter(ctx context.Context, p model.Parameter) error {
	return s.repo.Create(ctx, p)
}

func (s *ParameterService) UpdateParameter(ctx context.Context, p model.Parameter) error {
	return s.repo.Update(ctx, p)
}

func (s *ParameterService) DeleteParameter(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
