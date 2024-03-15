package service

import (
	"context"

	"../model"
	"../repository"
)

// UserService defines the interface for user operations.
type UserService interface {
	AddUser(ctx context.Context, user model.User) error
}

type userService struct {
	userRepository repository.UserRepository
}

// NewUserService creates a new instance of the user service.
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepository: userRepo}
}

func (s *userService) AddUser(ctx context.Context, user model.User) error {
	return s.userRepository.Add(ctx, user)
}
