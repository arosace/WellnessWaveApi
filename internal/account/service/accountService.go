package service

import (
	"errors"
	"fmt"

	"github.com/arosace/WellnessWaveApi/internal/account/domain"
	"github.com/arosace/WellnessWaveApi/internal/account/model"
	"github.com/arosace/WellnessWaveApi/internal/account/repository"
	"github.com/arosace/WellnessWaveApi/pkg/utils"
	encryption "github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/models"
)

// AccountService defines the interface for account operations.
type AccountService interface {
	AddAccount(ctx echo.Context, account model.Account) (*models.Record, error)
	GetAccounts(ctx echo.Context) ([]*models.Record, error)
	GetAccountById(ctx echo.Context, id string) (*models.Record, error)
	GetAttachedAccounts(ctx echo.Context, parentId string) ([]*models.Record, error)
	GetAccountByEmail(ctx echo.Context, email string) (*models.Record, error)
	CheckAccountExists(ctx echo.Context, email string) bool
	AttachAccount(ctx echo.Context, accountToAttach model.AttachAccountBody) (*models.Record, error)
	UpdateAccount(ctx echo.Context, accountToUpdate model.Account, infoType string) (*models.Record, error)
	Authorize(ctx echo.Context, credentials model.LogInCredentials) (*models.Record, error)
	VerifyAccount(echo.Context, string) (*models.Record, error)
}

type accountService struct {
	accountRepository repository.AccountRepository
	encryptor         encryption.Encryption
}

// NewAccountService creates a new instance of the account service.
func NewAccountService(accountRepo repository.AccountRepository, encryptor encryption.Encryption) AccountService {
	return &accountService{
		accountRepository: accountRepo,
		encryptor:         encryptor,
	}
}

func (s *accountService) AddAccount(ctx echo.Context, account model.Account) (*models.Record, error) {
	encryptedPassword, err := s.encryptor.Encrypt(account.Password)
	if err != nil {
		return nil, errors.New("error when encrypting password")
	}
	account.EncryptedPassword = encryptedPassword
	return s.accountRepository.Add(ctx, account)
}

func (s *accountService) GetAccounts(ctx echo.Context) ([]*models.Record, error) {
	return s.accountRepository.List(ctx)
}

func (s *accountService) VerifyAccount(ctx echo.Context, email string) (*models.Record, error) {
	record, err := s.GetAccountByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !record.Verified() {
		err = s.accountRepository.UpdateVerify(ctx, record)
		if err != nil {
			return nil, fmt.Errorf("there was an error verifying the account: %w", err)
		}
	}

	return record, nil
}

func (s *accountService) CheckAccountExists(ctx echo.Context, email string) bool {
	_, err := s.accountRepository.FindByEmail(ctx, email)
	if err != nil {
		return false
	}
	return true
}

func (s *accountService) AttachAccount(ctx echo.Context, accountToAttach model.AttachAccountBody) (*models.Record, error) {
	var account *models.Record

	//check if parentAccount exists
	parent, err := s.GetAccountById(ctx, accountToAttach.ParentID)
	if err != nil {
		if err.Error() == "not_found" {
			return nil, errors.New("parent account not found")
		}
		return nil, err
	}

	//one should not be able to attach to themselves
	if parent.Email() == accountToAttach.Email {
		return nil, errors.New("cannot attach account to itself")
	}

	//check if accountToAttach exists
	account, err = s.GetAccountByEmail(ctx, accountToAttach.Email)
	if err != nil && err.Error() != "not_found" {
		return nil, fmt.Errorf("there was an error retrieving account by email: %w", err)
	}

	//if it does not exist create account and attach
	//account will be created with random password, this will need to be changed by patient
	if account == nil {
		randPassword, err := utils.GenerateRandomPassword(8)
		if err != nil {
			return nil, errors.New("error generating eandom password for attached account")
		}
		encryptedRandPassword, err := s.encryptor.Encrypt(randPassword)
		if err != nil {
			return nil, errors.New("error encrypting random password for attached account")
		}

		newAccount, err := s.accountRepository.Add(ctx, model.Account{
			FirstName:         accountToAttach.FirstName,
			LastName:          accountToAttach.LastName,
			Email:             accountToAttach.Email,
			Role:              domain.PatientRole,
			ParentID:          accountToAttach.ParentID,
			Password:          randPassword,
			EncryptedPassword: encryptedRandPassword,
			Username:          fmt.Sprintf("%s %s", accountToAttach.FirstName, accountToAttach.LastName),
		})
		if err != nil {
			return nil, err
		}
		return newAccount, nil
	} else { //if it exists
		if account.GetString("parent_id") != "" && account.GetString("parent_id") != accountToAttach.ParentID { //if it is already attached return error
			return nil, errors.New("cannot attach an account that is already attached to another")
		} else if account.GetString("parent_id") == accountToAttach.ParentID { //is account is already attached to right health specialist
			return nil, nil
		} else { //if it is not already attached update it
			_, err := s.accountRepository.Attach(ctx, account, accountToAttach.ParentID)
			if err != nil {
				return nil, err
			}
			// send patient account updated event
		}
	}
	return account, nil
}

func (s *accountService) GetAccountByEmail(ctx echo.Context, email string) (*models.Record, error) {
	account, err := s.accountRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *accountService) GetAccountById(ctx echo.Context, id string) (*models.Record, error) {
	account, err := s.accountRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *accountService) GetAttachedAccounts(ctx echo.Context, parentId string) ([]*models.Record, error) {
	return s.accountRepository.FindByParentID(ctx, parentId)
}

func (s *accountService) UpdateAccount(ctx echo.Context, account model.Account, infoType string) (*models.Record, error) {
	oldAccount, err := s.accountRepository.FindByID(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	isToUpdate, nameChanged := false, false
	if infoType == "personal" {
		if account.FirstName != oldAccount.GetString("first_name") && account.FirstName != "" {
			isToUpdate = true
			oldAccount.Set("first_name", account.FirstName)
		}
		if account.LastName != oldAccount.GetString("last_name") && account.LastName != "" {
			isToUpdate = true
			oldAccount.Set("last_name", account.FirstName)
		}
		if account.ParentID != oldAccount.GetString("parent_id") && account.ParentID != "" {
			isToUpdate = true
			oldAccount.Set("parent_id", account.ParentID)
		}

		if !isToUpdate {
			return oldAccount, nil
		}

		if nameChanged {
			oldAccount.Set("username", fmt.Sprintf("%s %s", account.FirstName, account.LastName))
		}

		return s.accountRepository.Update(ctx, oldAccount)
	}

	decryptedOldPassword, err := s.encryptor.Decrypt(oldAccount.GetString("encrypted_password"))
	if err != nil {
		return nil, fmt.Errorf("error decrypting old account password: %w", err)
	}

	if account.Password != decryptedOldPassword {
		isToUpdate = true
		if utils.PasswordIsValid(account.Password) {
			encryptedPassword, err := s.encryptor.Encrypt(account.Password)
			if err != nil {
				return nil, errors.New("error when encrypting password")
			}
			oldAccount.Set("encrypted_password", encryptedPassword)
			oldAccount.SetPassword(encryptedPassword)
		} else {
			return nil, errors.New("password does not comply with authentication rules")
		}
	}
	if account.Email != oldAccount.Email() && utils.EmailIsValid(account.Email) {
		isToUpdate = true
		oldAccount.SetEmail(account.Email)
	}

	if isToUpdate {
		return s.accountRepository.Update(ctx, oldAccount)
	}

	return oldAccount, nil
}

func (s *accountService) Authorize(ctx echo.Context, credentials model.LogInCredentials) (*models.Record, error) {
	account, err := s.accountRepository.FindByEmail(ctx, credentials.Email)
	if err != nil {
		return nil, err
	}
	if !account.Verified() {
		return nil, errors.New("account_not_verified")
	}

	decryptedPassword, err := s.encryptor.Decrypt(account.GetString("encrypted_password"))
	if err != nil {
		return nil, err
	}

	if decryptedPassword != credentials.Password {
		return nil, errors.New("not_authorized")
	}
	return account, nil
}
