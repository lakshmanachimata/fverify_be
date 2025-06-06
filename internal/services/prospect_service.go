package services

import (
	"context"
	"fverify_be/internal/models"
	"fverify_be/internal/repositories"
)

type ProspectService struct {
	repo *repositories.ProspectRepositoryImpl
}

func NewProspectService(repo *repositories.ProspectRepositoryImpl) *ProspectService {
	return &ProspectService{repo: repo}
}

func (s *ProspectService) CreateProspect(ctx context.Context, prospect *models.Prospect) error {
	return s.repo.Create(ctx, prospect)
}

func (s *ProspectService) GetProspectByID(ctx context.Context, id string) (*models.Prospect, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProspectService) UpdateProspect(ctx context.Context, prospect *models.Prospect) error {
	return s.repo.Update(ctx, prospect)
}

func (s *ProspectService) DeleteProspect(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *ProspectService) ListProspects(ctx context.Context) ([]*models.Prospect, error) {
	return s.repo.FindAll(ctx)
}
func (s *ProspectService) GetProspects(ctx context.Context, skip int, limit int) ([]models.Prospect, error) {
	return s.repo.GetProspects(ctx, skip, limit)
}
func (s *ProspectService) GetProspectsCount(ctx context.Context) (int, error) {
	return s.repo.GetProspectsCount(ctx)
}
