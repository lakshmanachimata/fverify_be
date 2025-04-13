package services

import (
	"context"
	"time"

	"kowtha_be/internal/models"

	"kowtha_be/internal/repositories"
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

func (s *UserService) GetByUserID(ctx context.Context, userId string) (*models.UserModel, error) {
	return s.repo.GetByUserID(ctx, userId)
}
func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.UserModel, error) {
	return s.repo.GetAllUsers(ctx)
}
