package services

import (
	"context"
	"time"

	"fverify_be/internal/models"

	"fverify_be/internal/repositories"
)

type UserService struct {
	repo *repositories.UserRepositoryImpl
}

func NewUserService(repo *repositories.UserRepositoryImpl) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.UserModel) (*models.UserModel, error) {
	user.CreatedTime = time.Now()
	user.UpdatedTime = time.Now()
	user.Status = models.Created
	return s.repo.Create(ctx, user)
}

func (s *UserService) GetByUserUID(ctx context.Context, uid string) (*models.UserModel, error) {
	return s.repo.GetByUserUID(ctx, uid)
}

func (s *UserService) GetByUserID(ctx context.Context, userId string) (*models.UserModel, error) {
	return s.repo.GetByUserID(ctx, userId)
}
func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.UserModel, error) {
	return s.repo.GetAllUsers(ctx)
}

func (s *UserService) DeleteByUId(ctx context.Context, uId string) error {
	return s.repo.DeleteByUId(ctx, uId)
}

func (s *UserService) DeleteByUserId(ctx context.Context, userId string) error {
	return s.repo.DeleteByUserId(ctx, userId)
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.UserModel) error {
	return s.repo.Update(ctx, user)
}
func (s *UserService) LoginUser(ctx context.Context, username, password string, orgId string) (*models.UserModel, error) {
	return s.repo.ValidateUser(ctx, username, password, orgId)
}
func (s *UserService) SetPassword(ctx context.Context, uId string, newPassword string) error {
	return s.repo.SetPassword(ctx, uId, newPassword)
}
