package services

import (
	"context"
	"fverify_be/internal/models"
	"fverify_be/internal/repositories"
)

type OrganisationService struct {
	repo     *repositories.OrganisationRepository
	userRepo *repositories.UserRepositoryImpl
}

func NewOrganisationService(repo *repositories.OrganisationRepository, userRepo *repositories.UserRepositoryImpl) *OrganisationService {
	return &OrganisationService{repo: repo, userRepo: userRepo}
}

func (s *OrganisationService) CreateOrganisation(ctx context.Context, org *models.Organisation) (*models.Organisation, error) {
	return s.repo.Create(ctx, org)
}

func (s *OrganisationService) UpdateOrganisation(ctx context.Context, org_id string, org *models.Organisation) error {
	return s.repo.Update(ctx, org_id, org)
}

func (s *OrganisationService) DeleteOrganisation(ctx context.Context, org_id string) error {
	return s.repo.Delete(ctx, org_id)
}
func (s *OrganisationService) GetAllOrganisations(ctx context.Context) ([]*models.Organisation, error) {
	return s.repo.GetAllOrganisations(ctx)
}
func (s *OrganisationService) IsOrgActive(ctx context.Context, org_id string) (bool, error) {
	return s.repo.IsOrgActive(ctx, org_id)
}
func (s *OrganisationService) GetOrganisationByID(ctx context.Context, org_id string) (*models.Organisation, error) {
	return s.repo.GetOrganisationByID(ctx, org_id)
}
func (s *OrganisationService) UpdateUsersStatusByOrgUUID(ctx context.Context, orgUUID string, status models.UserStatus) error {
	return s.userRepo.UpdateUsersStatusByOrgUUID(ctx, orgUUID, status)
}
