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
	GetAttachedAccounts(ctx context.Context, parentId string) ([]*model.Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*model.Account, error)
	CheckAccountExists(ctx context.Context, email string) bool
	AttachAccount(ctx context.Context, accountToAttach model.AttachAccountBody) (*model.Account, error)
	UpdateAccount(ctx context.Context, accountToUpdate model.Account, infoType string) (*model.Account, error)
	Authorize(ctx context.Context, credentials model.LogInCredentials) (*model.Account, error)
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

func (s *accountService) AttachAccount(ctx context.Context, accountToAttach model.AttachAccountBody) (*model.Account, error) {
	var account *model.Account

	//check if parentAccount exists
	parent, err := s.GetAccountByID(ctx, accountToAttach.ParentID)
	if err != nil {
		if err.Error() == "not_found" {
			return nil, errors.New("parent account not found")
		}
		return nil, err
	}

	//one should not be able to attach to themselves
	if parent.Email == accountToAttach.Email {
		return nil, errors.New("cannot attach account to itself")
	}

	//check if accountToAttach exists
	account, err = s.GetAccountByEmail(ctx, accountToAttach.Email)
	if err != nil && err.Error() != "not_found" {
		return nil, err
	}
	//if it does not exist create account and attach
	//account will be created without password, for now it's just for record
	if account == nil {
		newAccount, err := s.accountRepository.Add(ctx, model.Account{
			FirstName: accountToAttach.FirstName,
			LastName:  accountToAttach.LastName,
			Email:     accountToAttach.Email,
			Role:      domain.PatientRole,
			ParentID:  accountToAttach.ParentID,
		})
		if err != nil {
			return nil, err
		}
		return newAccount, nil
		// send patient account created event
	} else { //if it exists
		if account.ParentID != "" && account.ParentID != accountToAttach.ParentID { //if it is already attached return error
			return nil, errors.New("cannot attach an account that is already attached to another")
		} else if account.ParentID == accountToAttach.ParentID { //is account is already attached to right health specialist
			return nil, nil
		} else { //if it is not already attached update it
			account.ParentID = accountToAttach.ParentID
			_, err := s.accountRepository.Update(ctx, account)
			if err != nil {
				return nil, err
			}
			// send patient account updated event
		}
	}
	return account, nil
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

func (s *accountService) GetAttachedAccounts(ctx context.Context, parentId string) ([]*model.Account, error) {
	return s.accountRepository.FindByParentID(ctx, parentId)
}

func (s *accountService) UpdateAccount(ctx context.Context, account model.Account, infoType string) (*model.Account, error) {
	oldAccount, err := s.accountRepository.FindByID(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	isToUpdate := false
	if infoType == "personal" {
		if account.FirstName != oldAccount.FirstName {
			isToUpdate = true
			oldAccount.FirstName = account.FirstName
		}
		if account.LastName != oldAccount.LastName {
			isToUpdate = true
			oldAccount.LastName = account.LastName
		}
		if account.ParentID != oldAccount.ParentID {
			isToUpdate = true
			oldAccount.ParentID = account.ParentID
		}

		if !isToUpdate {
			return oldAccount, nil
		}

		return s.accountRepository.Update(ctx, oldAccount)
	}

	if account.Password != oldAccount.Password {
		isToUpdate = true
		encryptedPassword, err := s.encryptor.Encrypt(account.Password)
		if err != nil {
			return nil, errors.New("error when encrypting password")
		}
		oldAccount.Password = encryptedPassword
	}
	if account.Email != oldAccount.Email {
		isToUpdate = true
		oldAccount.Email = account.Email
	}

	if !isToUpdate {
		return oldAccount, nil
	}

	return s.accountRepository.UpdateAuth(ctx, oldAccount)
}

func (s *accountService) Authorize(ctx context.Context, credentials model.LogInCredentials) (*model.Account, error) {
	account, err := s.accountRepository.FindByEmail(ctx, credentials.Email)
	if err != nil {
		return nil, err
	}
	decryptedPassword, err := s.encryptor.Decrypt(account.Password)
	if err != nil {
		return nil, err
	}
	if decryptedPassword != credentials.Password {
		return nil, errors.New("not_authorized")
	}
	return account, nil
}
