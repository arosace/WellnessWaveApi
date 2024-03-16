package service

import (
	"context"

	"../model"
	"../repository"
)

// UserService defines the interface for user operations.
type AccountService interface {
	AddAccount(ctx context.Context, user model.Account) error
}

type accountService struct {
	accountRepository repository.AccountRepository
}

// NewUserService creates a new instance of the user service.
func NewUserService(userRepo repository.AccountRepository) AccountService {
	return &accountService{accountRepository: userRepo}
}

func (s *accountService) AddAccount(ctx context.Context, user model.Account) error {
	return s.accountRepository.Add(ctx, user)
}
