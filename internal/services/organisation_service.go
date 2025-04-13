package services

import (
	"context"
	"kowtha_be/internal/models"
	"kowtha_be/internal/repositories"
)

type OrganisationService struct {
	repo *repositories.OrganisationRepository
}

func NewOrganisationService(repo *repositories.OrganisationRepository) *OrganisationService {
	return &OrganisationService{repo: repo}
}

func (s *OrganisationService) CreateOrganisation(ctx context.Context, org *models.Organisation) (*models.Organisation, error) {
	return s.repo.Create(ctx, org)
}

func (s *OrganisationService) UpdateOrganisation(ctx context.Context, orgId string, org *models.Organisation) error {
	return s.repo.Update(ctx, orgId, org)
}

func (s *OrganisationService) DeleteOrganisation(ctx context.Context, orgId string) error {
	return s.repo.Delete(ctx, orgId)
}
func (s *OrganisationService) GetAllOrganisations(ctx context.Context) ([]*models.Organisation, error) {
	return s.repo.GetAllOrganisations(ctx)
}
