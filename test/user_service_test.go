package test

import (
	"context"
	"testing"

	"fverify_be/internal/models"
	"fverify_be/internal/repositories"
	"fverify_be/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
	*repositories.UserRepositoryImpl
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) (*models.UserResp, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*models.UserResp), args.Error(1)
}

func (m *MockUserRepository) GetByUserUID(ctx context.Context, uid string) (*models.UserResp, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(*models.UserResp), args.Error(1)
}

func (m *MockUserRepository) GetByUserID(ctx context.Context, userId string) (*models.UserResp, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(*models.UserResp), args.Error(1)
}

func (m *MockUserRepository) GetAllUsers(ctx context.Context) ([]*models.UserResp, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.UserResp), args.Error(1)
}

func (m *MockUserRepository) DeleteByUId(ctx context.Context, uId string) error {
	args := m.Called(ctx, uId)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteByUserId(ctx context.Context, userId string) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User, authUserName string) (*models.UserResp, error) {
	args := m.Called(ctx, user, authUserName)
	return args.Get(0).(*models.UserResp), args.Error(1)
}

func (m *MockUserRepository) ValidateUser(ctx context.Context, username, password, org_id string) (*models.User, error) {
	args := m.Called(ctx, username, password, org_id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) SetPassword(ctx context.Context, uId string, newPassword string) error {
	args := m.Called(ctx, uId, newPassword)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserStatus(ctx context.Context, userId string, status string) error {
	args := m.Called(ctx, userId, status)
	return args.Error(0)
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := &services.UserService{Repo: mockRepo}
	ctx := context.Background()

	tests := []struct {
		name    string
		user    *models.User
		mock    func()
		want    *models.UserResp
		wantErr bool
	}{
		{
			name: "Success",
			user: &models.User{
				UserId:   "testuser",
				Username: "testuser",
				Password: "password",
				Role:     "user",
				Status:   "active",
			},
			mock: func() {
				mockRepo.On("Create", ctx, mock.AnythingOfType("*models.User")).Return(
					&models.UserResp{
						UserId:   "testuser",
						Username: "testuser",
						Role:     "user",
						Status:   "active",
					}, nil)
			},
			want: &models.UserResp{
				UserId:   "testuser",
				Username: "testuser",
				Role:     "user",
				Status:   "active",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := service.CreateUser(ctx, tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserService_GetByUserUID(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := &services.UserService{Repo: mockRepo}
	ctx := context.Background()

	tests := []struct {
		name    string
		uid     string
		mock    func()
		want    *models.UserResp
		wantErr bool
	}{
		{
			name: "Success",
			uid:  "testuid",
			mock: func() {
				mockRepo.On("GetByUserUID", ctx, "testuid").Return(
					&models.UserResp{
						UserId:   "testuser",
						Username: "testuser",
						Role:     "user",
						Status:   "active",
					}, nil)
			},
			want: &models.UserResp{
				UserId:   "testuser",
				Username: "testuser",
				Role:     "user",
				Status:   "active",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := service.GetByUserUID(ctx, tt.uid)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
