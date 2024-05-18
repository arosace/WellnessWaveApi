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

// UserService defines the interface for user operations.
type AccountService interface {
	AddAccount(ctx echo.Context, account model.Account) (*models.Record, error)
	GetAccounts(ctx echo.Context) ([]*models.Record, error)
	GetAccountById(ctx echo.Context, id string) (*models.Record, error)
	GetAttachedAccounts(ctx echo.Context, parentId string) ([]*model.Account, error)
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

// NewUserService creates a new instance of the user service.
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
	account.Password = encryptedPassword
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

	err = s.accountRepository.UpdateVerify(ctx, record)

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

func (s *accountService) GetAttachedAccounts(ctx echo.Context, parentId string) ([]*model.Account, error) {
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

	if account.Password != oldAccount.GetString("password") {
		isToUpdate = true
		if utils.PasswordIsValid(account.Password) {
			encryptedPassword, err := s.encryptor.Encrypt(account.Password)
			if err != nil {
				return nil, errors.New("error when encrypting password")
			}
			oldAccount.SetPassword(encryptedPassword)
		} else {
			return nil, errors.New("password does not comply with authentication rules")
		}
	}
	if account.Email != oldAccount.Email() && utils.EmailIsValid(account.Email) {
		isToUpdate = true
		oldAccount.SetEmail(account.Email)
	}

	if !isToUpdate {
		return oldAccount, nil
	}

	return s.accountRepository.Update(ctx, oldAccount)
}

func (s *accountService) Authorize(ctx echo.Context, credentials model.LogInCredentials) (*models.Record, error) {
	account, err := s.accountRepository.FindByEmail(ctx, credentials.Email)
	if err != nil {
		return nil, err
	}
	if !account.GetBool("verified") {
		return nil, errors.New("account_not_verified")
	}
	decryptedPassword, err := s.encryptor.Decrypt(account.GetString("account_password"))
	if err != nil {
		return nil, err
	}

	if decryptedPassword != credentials.Password {
		return nil, errors.New("not_authorized")
	}
	return account, nil
}
