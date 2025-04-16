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

func (s *OrganisationService) UpdateOrganisation(ctx context.Context, orgId string, org *models.Organisation) error {
	return s.repo.Update(ctx, orgId, org)
}

func (s *OrganisationService) DeleteOrganisation(ctx context.Context, orgId string) error {
	return s.repo.Delete(ctx, orgId)
}
func (s *OrganisationService) GetAllOrganisations(ctx context.Context) ([]*models.Organisation, error) {
	return s.repo.GetAllOrganisations(ctx)
}
func (s *OrganisationService) IsOrgActive(ctx context.Context, orgId string) (bool, error) {
	return s.repo.IsOrgActive(ctx, orgId)
}
func (s *OrganisationService) GetOrganisationByID(ctx context.Context, orgId string) (*models.Organisation, error) {
	return s.repo.GetOrganisationByID(ctx, orgId)
}
func (s *OrganisationService) UpdateUsersStatusByOrgUUID(ctx context.Context, orgUUID string, status models.UserStatus) error {
	return s.userRepo.UpdateUsersStatusByOrgUUID(ctx, orgUUID, status)
}
