package service

import (
	"context"

	"github.com/arosace/WellnessWaveApi/internal/account/model"
	"github.com/arosace/WellnessWaveApi/internal/account/repository"
)

// UserService defines the interface for user operations.
type AccountService interface {
	AddAccount(ctx context.Context, user model.Account) error
	GetAccounts(ctx context.Context) ([]model.Account, error)
	CheckAccountExists(ctx context.Context, email string) bool
}

type accountService struct {
	accountRepository repository.AccountRepository
}

// NewUserService creates a new instance of the user service.
func NewAccountService(accountRepo repository.AccountRepository) AccountService {
	return &accountService{accountRepository: accountRepo}
}

func (s *accountService) AddAccount(ctx context.Context, user model.Account) error {
	return s.accountRepository.Add(ctx, user)
}

func (s *accountService) GetAccounts(ctx context.Context) ([]model.Account, error) {
	return s.accountRepository.List(ctx)
}

func (s *accountService) CheckAccountExists(ctx context.Context, email string) bool {
	_, err := s.accountRepository.FindByEmail(ctx, email)
	if err != nil {
		return false
	}

	return true
}
