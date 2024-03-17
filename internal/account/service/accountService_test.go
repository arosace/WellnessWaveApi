package service

import (
	"context"
	"errors"
	"testing"

	"github.com/arosace/WellnessWaveApi/internal/account/domain"
	"github.com/arosace/WellnessWaveApi/internal/account/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Add(ctx context.Context, account model.Account) (*model.Account, error) {
	args := m.Called(account)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), nil
}

func (m *MockAccountRepository) List(ctx context.Context) ([]*model.Account, error) {
	args := m.Called()
	return args.Get(0).([]*model.Account), args.Error(1)
}

func (m *MockAccountRepository) FindByEmail(ctx context.Context, email string) (*model.Account, error) {
	args := m.Called(email)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), nil
}

func (m *MockAccountRepository) FindByID(ctx context.Context, id string) (*model.Account, error) {
	args := m.Called(id)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), nil
}

func (m *MockAccountRepository) Update(ctx context.Context, account *model.Account) (*model.Account, error) {
	args := m.Called(account)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), nil
}

type MockEncryptor struct {
	mock.Mock
}

func (m *MockEncryptor) Encrypt(data string) (string, error) {
	args := m.Called(data)
	return args.String(0), args.Error(1)
}

func (m *MockEncryptor) Decrypt(data string) (string, error) {
	args := m.Called(data)
	return args.String(0), args.Error(1)
}

func TestAddAccount(t *testing.T) {
	testAccount := model.Account{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
		Password:  "password123",
	}

	encryptedPassword := "password123"

	t.Run("when error when encrypting password, return error", func(t *testing.T) {
		mockRepo := new(MockAccountRepository)
		mockEncryptor := new(MockEncryptor)
		accountService := NewAccountService(mockRepo, mockEncryptor)
		mockEncryptor.On("Encrypt", testAccount.Password).Return("", errors.New("encryption error"))

		acc, err := accountService.AddAccount(context.Background(), testAccount)

		assert.Nil(t, acc)
		assert.EqualError(t, err, "error when encrypting password")
		mockEncryptor.AssertExpectations(t)
	})

	t.Run("when successful account addition, return account", func(t *testing.T) {
		mockRepo := new(MockAccountRepository)
		mockEncryptor := new(MockEncryptor)
		accountService := NewAccountService(mockRepo, mockEncryptor)
		mockEncryptor.On("Encrypt", testAccount.Password).Return(encryptedPassword, nil)
		mockRepo.On("Add", mock.Anything).Return(&testAccount, nil)

		acc, err := accountService.AddAccount(context.Background(), testAccount)

		assert.Nil(t, err)
		assert.NotNil(t, acc)
		assert.Equal(t, encryptedPassword, acc.Password)
		mockRepo.AssertExpectations(t)
		mockEncryptor.AssertExpectations(t)
	})
}

func TestAttachAccount(t *testing.T) {
	t.Run("when attach is successfull and account does not exist, return nil", func(t *testing.T) {
		mockRepo := new(MockAccountRepository)
		mockEncryptor := new(MockEncryptor)
		accountService := NewAccountService(mockRepo, mockEncryptor)
		mockAttachBody := model.AttachAccountBody{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			ParentID:  "parent123",
		}
		mockParent := model.Account{Role: domain.PatientRole, ParentID: "parent123"}
		mockAccount := model.Account{Role: domain.PatientRole, ParentID: "parent123"}
		mockRepo.On("FindByID", mockAttachBody.ParentID).Return(&mockParent, nil)
		mockRepo.On("FindByEmail", mockAttachBody.Email).Return(nil, errors.New("not_found"))
		mockRepo.On("Add", mock.Anything).Return(&mockAccount, nil)

		err := accountService.AttachAccount(context.Background(), mockAttachBody)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("when attach is successfull and account exists and is already attached to another account, return error", func(t *testing.T) {
		mockRepo := new(MockAccountRepository)
		mockEncryptor := new(MockEncryptor)
		accountService := NewAccountService(mockRepo, mockEncryptor)
		mockAttachBody := model.AttachAccountBody{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			ParentID:  "parent567",
		}
		mockParent := model.Account{Role: domain.PatientRole, ParentID: "parent123"}
		mockAccount := model.Account{Role: domain.PatientRole, ParentID: "parent123"}
		mockRepo.On("FindByID", mockAttachBody.ParentID).Return(&mockParent, nil)
		mockRepo.On("FindByEmail", mockAttachBody.Email).Return(&mockAccount, nil)

		err := accountService.AttachAccount(context.Background(), mockAttachBody)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "another")
	})

	t.Run("when attach is successfull and account exists and is already attached to the account, return nil", func(t *testing.T) {
		mockRepo := new(MockAccountRepository)
		mockEncryptor := new(MockEncryptor)
		accountService := NewAccountService(mockRepo, mockEncryptor)
		mockAttachBody := model.AttachAccountBody{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			ParentID:  "parent123",
		}
		mockParent := model.Account{Role: domain.PatientRole, ParentID: "parent123"}
		mockAccount := model.Account{Role: domain.PatientRole, ParentID: "parent123"}
		mockRepo.On("FindByID", mockAttachBody.ParentID).Return(&mockParent, nil)
		mockRepo.On("FindByEmail", mockAttachBody.Email).Return(&mockAccount, nil)

		err := accountService.AttachAccount(context.Background(), mockAttachBody)

		assert.Nil(t, err)
	})

	t.Run("when attach is successfull and account exists and is not already attached to any account, return nil", func(t *testing.T) {
		mockRepo := new(MockAccountRepository)
		mockEncryptor := new(MockEncryptor)
		accountService := NewAccountService(mockRepo, mockEncryptor)
		mockAttachBody := model.AttachAccountBody{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			ParentID:  "parent123",
		}
		mockParent := model.Account{Role: domain.PatientRole, ParentID: "parent123"}
		mockAccount := model.Account{Role: domain.PatientRole, ParentID: "parent123"}
		mockRepo.On("FindByID", mockAttachBody.ParentID).Return(&mockParent, nil)
		mockRepo.On("FindByEmail", mockAttachBody.Email).Return(&mockAccount, nil)
		mockRepo.On("Update", mock.Anything).Return(&mockAccount, nil)

		err := accountService.AttachAccount(context.Background(), mockAttachBody)

		assert.Nil(t, err)
	})

	t.Run("when update is not successfull, return error", func(t *testing.T) {
		mockRepo := new(MockAccountRepository)
		mockEncryptor := new(MockEncryptor)
		accountService := NewAccountService(mockRepo, mockEncryptor)
		mockAttachBody := model.AttachAccountBody{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			ParentID:  "parent123",
		}
		mockParent := model.Account{Role: domain.PatientRole, ID: "parent123"}
		mockAccount := model.Account{Role: domain.PatientRole}
		mockRepo.On("FindByID", mockAttachBody.ParentID).Return(&mockParent, nil)
		mockRepo.On("FindByEmail", mockAttachBody.Email).Return(&mockAccount, nil)
		mockRepo.On("Update", mock.Anything).Return(nil, errors.New("err"))

		err := accountService.AttachAccount(context.Background(), mockAttachBody)

		assert.Error(t, err)
	})

	t.Run("when parent account does not exits, return error", func(t *testing.T) {
		mockRepo := new(MockAccountRepository)
		mockEncryptor := new(MockEncryptor)
		accountService := NewAccountService(mockRepo, mockEncryptor)
		mockAttachBody := model.AttachAccountBody{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			ParentID:  "parent123",
		}
		mockRepo.On("FindByID", mockAttachBody.ParentID).Return(nil, errors.New("not_found"))
		err := accountService.AttachAccount(context.Background(), mockAttachBody)

		assert.Error(t, err)
	})
}
