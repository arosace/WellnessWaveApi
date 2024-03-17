package service

import (
	"context"
	"errors"

	"github.com/arosace/WellnessWaveApi/internal/account/domain"
	"github.com/arosace/WellnessWaveApi/internal/account/model"
	"github.com/arosace/WellnessWaveApi/internal/account/repository"
	encryption "github.com/arosace/WellnessWaveApi/pkg/utils"
)

// UserService defines the interface for user operations.
type AccountService interface {
	AddAccount(ctx context.Context, account model.Account) (*model.Account, error)
	GetAccounts(ctx context.Context) ([]*model.Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*model.Account, error)
	CheckAccountExists(ctx context.Context, email string) bool
	AttachAccount(ctx context.Context, accountToAttach model.AttachAccountBody) error
}

type accountService struct {
	accountRepository repository.AccountRepository
	encryptor         encryption.Encryption
}

// NewUserService creates a new instance of the user service.
func NewAccountService(accountRepo repository.AccountRepository, encryptor encryption.Encryption) AccountService {
	return &accountService{
		accountRepository: accountRepo,
		encryptor:         encryptor,
	}
}

func (s *accountService) AddAccount(ctx context.Context, account model.Account) (*model.Account, error) {
	encryptedPassword, err := s.encryptor.Encrypt(account.Password)
	if err != nil {
		return nil, errors.New("error when encrypting password")
	}
	account.Password = encryptedPassword
	return s.accountRepository.Add(ctx, account)
}

func (s *accountService) GetAccounts(ctx context.Context) ([]*model.Account, error) {
	return s.accountRepository.List(ctx)
}

func (s *accountService) CheckAccountExists(ctx context.Context, email string) bool {
	_, err := s.accountRepository.FindByEmail(ctx, email)
	if err != nil {
		return false
	}
	return true
}

func (s *accountService) AttachAccount(ctx context.Context, accountToAttach model.AttachAccountBody) error {
	var account *model.Account

	//check if parentAccount exists
	_, err := s.GetAccountByID(ctx, accountToAttach.ParentID)
	if err != nil {
		if err.Error() == "not_found" {
			return errors.New("parent account not found")
		}
		return err
	}

	//check if accountToAttach exists
	account, err = s.GetAccountByEmail(ctx, accountToAttach.Email)
	if err != nil && err.Error() != "not_found" {
		return err
	}
	//if it does not exist create account and attach
	//account will be created without password, for now it's just for record
	if account == nil {
		_, err := s.accountRepository.Add(ctx, model.Account{
			FirstName: accountToAttach.FirstName,
			LastName:  accountToAttach.LastName,
			Email:     accountToAttach.Email,
			Role:      domain.PatientRole,
			ParentID:  accountToAttach.ParentID,
		})
		if err != nil {
			return err
		}
		// send patient account created event
	} else { //if it exists
		if account.ParentID != "" && account.ParentID != accountToAttach.ParentID { //if it is already attached return error
			return errors.New("cannot attach an account that is already attached to another")
		} else if account.ParentID == accountToAttach.ParentID { //is account is already attached to right health specialist
			return nil
		} else { //if it is not already attached update it
			account.ParentID = accountToAttach.ParentID
			_, err := s.accountRepository.Update(ctx, account)
			if err != nil {
				return err
			}
			// send patient account updated event
		}
	}
	return nil
}

func (s *accountService) GetAccountByEmail(ctx context.Context, email string) (*model.Account, error) {
	account, err := s.accountRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *accountService) GetAccountByID(ctx context.Context, id string) (*model.Account, error) {
	account, err := s.accountRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return account, nil
}
