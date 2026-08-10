package service

import (
	"context"

	"example.com/my-api/internal/dto"
	"example.com/my-api/internal/model"
	"example.com/my-api/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error)
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, id int64) (model.User, error)
	UpdateUser(ctx context.Context, id int64, req dto.UpdateUserRequest) (dto.UserResponse, error)
	DeleteUser(ctx context.Context, id int64) error
}

type userService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{repository: repository}
}

func (s *userService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error) {
	user, err := s.repository.Create(ctx, model.User{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
	}, nil
}

func (s *userService) GetUsers(ctx context.Context) ([]model.User, error) {
	return s.repository.FindAll(ctx)
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (model.User, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, id int64, req dto.UpdateUserRequest) (dto.UserResponse, error) {
	user, err := s.repository.Update(ctx, model.User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}
