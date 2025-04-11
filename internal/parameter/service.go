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
func (s *Service) GetParameter(ctx context.Context, id int) (Parameter, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) CreateParameter(ctx context.Context, p Parameter) error {
	return s.repo.Create(ctx, p)
}

func (s *Service) UpdateParameter(ctx context.Context, p Parameter) error {
	return s.repo.Update(ctx, p)
}

func (s *Service) DeleteParameter(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
