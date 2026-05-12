package service

import (
	"context"

	"github.com/example/microservices/user-service/internal/model"
	"github.com/example/microservices/user-service/internal/repository"
	appErrors "github.com/example/microservices/shared/pkg/errors"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*model.User, error) {
	if userID == "" {
		return nil, appErrors.NewValidationError("User ID is required", nil)
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, appErrors.NewNotFoundError("User not found")
	}

	return user, nil
}

func (s *UserService) ListUsers(ctx context.Context, page int, pageSize int) ([]model.User, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	users, total, err := s.repo.ListUsers(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, appErrors.NewInternalError("Failed to list users", err)
	}

	if users == nil {
		users = []model.User{}
	}

	return users, total, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID string, req *model.UpdateUserRequest) (*model.User, error) {
	if userID == "" {
		return nil, appErrors.NewValidationError("User ID is required", nil)
	}

	user, err := s.repo.UpdateUser(ctx, userID, req)
	if err != nil {
		return nil, appErrors.NewNotFoundError("User not found")
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	if userID == "" {
		return appErrors.NewValidationError("User ID is required", nil)
	}

	err := s.repo.DeleteUser(ctx, userID)
	if err != nil {
		return appErrors.NewNotFoundError("User not found")
	}

	return nil
}
