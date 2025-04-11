package parameter

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListParameters(ctx context.Context) ([]Parameter, error) {
	return s.repo.GetAll(ctx)
}
