package service

import (
	"context"

	"example.com/my-api/internal/model"
	"example.com/my-api/internal/repository"
)

type UserService interface {
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, id int64) (model.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type userService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{repository: repository}
}

func (s *userService) GetUsers(ctx context.Context) ([]model.User, error) {
	return s.repository.FindAll(ctx)
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (model.User, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}
