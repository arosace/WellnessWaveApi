package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arosace/WellnessWaveApi/internal/account/domain"
	"github.com/arosace/WellnessWaveApi/internal/account/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAccountService struct {
	mock.Mock
}

func (s *mockAccountService) AddAccount(ctx context.Context, account model.Account) (*model.Account, error) {
	args := s.Called(account)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), nil
}

func (s *mockAccountService) CheckAccountExists(ctx context.Context, email string) bool {
	args := s.Called(email)
	return args.Bool(0)
}

func (s *mockAccountService) GetAccounts(ctx context.Context) ([]*model.Account, error) {
	args := s.Called()
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Account), nil
}

func (s *mockAccountService) AttachAccount(ctx context.Context, accountToAttach model.AttachAccountBody) error {
	args := s.Called(accountToAttach)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (s *mockAccountService) GetAccountByEmail(ctx context.Context, email string) (*model.Account, error) {
	args := s.Called(email)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), nil
}

func TestHandleAddAccount(t *testing.T) {
	t.Run("when post method is used, return method not allowed", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/accounts/register", bytes.NewBuffer([]byte{}))
		if err != nil {
			t.Fatal(err)
		}
		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)
		handler.HandleAddAccount(rr, req)
		assert.Equal(t, rr.Code, 405)
	})
	t.Run("when wrong body is sent, return bad request", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/accounts/register", bytes.NewBuffer([]byte{}))
		if err != nil {
			t.Fatal(err)
		}
		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)
		handler.HandleAddAccount(rr, req)
		assert.Equal(t, rr.Code, 400)
		assert.Contains(t, rr.Body.String(), "wrong_data_type")
	})

	t.Run("when wrong body is not valid, return bad request", func(t *testing.T) {
		testUser := model.Account{
			ID:        "test_id",
			FirstName: "Name",
			LastName:  "Surname",
			Email:     "test@example.com",
		}
		userJSON, _ := json.Marshal(testUser)
		rr := httptest.NewRecorder()

		req, err := http.NewRequest("POST", "/accounts/register", bytes.NewBuffer(userJSON))
		if err != nil {
			t.Fatal(err)
		}

		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)
		handler.HandleAddAccount(rr, req)

		assert.Equal(t, rr.Code, 400)
		assert.Contains(t, rr.Body.String(), "missing_data: role, password")
	})

	t.Run("when account already exists, return conflict", func(t *testing.T) {
		testUser := model.Account{
			ID:        "test_id",
			FirstName: "Name",
			LastName:  "Surname",
			Email:     "test@example.com",
			Role:      "role",
			Password:  "[]",
		}
		userJSON, _ := json.Marshal(testUser)
		rr := httptest.NewRecorder()

		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)
		mockService.On("CheckAccountExists", testUser.Email).Return(true)

		req, err := http.NewRequest("POST", "/accounts/register", bytes.NewBuffer(userJSON))
		if err != nil {
			t.Fatal(err)
		}

		handler.HandleAddAccount(rr, req)

		assert.Equal(t, rr.Code, 409)
		assert.Contains(t, rr.Body.String(), "email_already_in_use")
	})

	t.Run("when adding account returns an error, return internal server error", func(t *testing.T) {
		testUser := model.Account{
			ID:        "test_id",
			FirstName: "Name",
			LastName:  "Surname",
			Email:     "test@example.com",
			Role:      "role",
			Password:  "[]",
		}
		userJSON, _ := json.Marshal(testUser)
		rr := httptest.NewRecorder()

		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)
		mockService.On("CheckAccountExists", testUser.Email).Return(false)
		mockService.On("AddAccount", testUser).Return(nil, errors.New("error"))

		req, err := http.NewRequest("POST", "/accounts/register", bytes.NewBuffer(userJSON))
		if err != nil {
			t.Fatal(err)
		}

		handler.HandleAddAccount(rr, req)

		assert.Equal(t, rr.Code, 500)
	})
}

func TestHandleGetAccounts(t *testing.T) {
	t.Run("when non-GET method is used, return method not allowed", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/accounts", nil) // Using POST instead of GET
		if err != nil {
			t.Fatal(err)
		}
		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)
		handler.HandleGetAccounts(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("when service successfully returns accounts, return OK with accounts", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/accounts", nil)
		if err != nil {
			t.Fatal(err)
		}
		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)

		testAccounts := []*model.Account{
			{ID: "1", FirstName: "Test", LastName: "User", Email: "test@example.com"},
			// Add more test accounts as needed
		}

		mockService.On("GetAccounts", mock.Anything).Return(testAccounts, nil)

		handler.HandleGetAccounts(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var accounts []*model.Account
		err = json.NewDecoder(rr.Body).Decode(&accounts)
		assert.NoError(t, err)
		assert.Equal(t, testAccounts, accounts)
	})

	t.Run("when service fails to get accounts, return internal server error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/accounts", nil)
		if err != nil {
			t.Fatal(err)
		}
		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)

		mockService.On("GetAccounts", mock.Anything).Return(nil, errors.New("internal error"))

		handler.HandleGetAccounts(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestHandleAttachAccount(t *testing.T) {
	t.Run("when non-POST method is used, return method not allowed", func(t *testing.T) {
		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)

		req, err := http.NewRequest(http.MethodGet, "/accounts/attach", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.HandleAttachAccount(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("when wrong body is sent, return bad request", func(t *testing.T) {
		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)
		req, err := http.NewRequest(http.MethodPost, "/accounts/attach", bytes.NewBufferString("invalid body"))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.HandleAttachAccount(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "wrong_data_type")
	})

	t.Run("when data is invalid, return bad request", func(t *testing.T) {
		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)
		account := model.AttachAccountBody{}
		body, _ := json.Marshal(account)

		req, err := http.NewRequest(http.MethodPost, "/accounts/attach", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.HandleAttachAccount(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "missing_data: email")
	})

	t.Run("when role is invalid, return bad request", func(t *testing.T) {
		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)
		account := model.AttachAccountBody{
			Role:      "invalid",
			FirstName: "first",
			LastName:  "last",
			Email:     "email",
			ParentID:  "id",
		}
		body, _ := json.Marshal(account)

		req, err := http.NewRequest(http.MethodPost, "/accounts/attach", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.HandleAttachAccount(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid_role")
	})
	t.Run("when account service returns an error, return internal server error", func(t *testing.T) {
		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)
		account := model.AttachAccountBody{
			Role:      domain.PatientRole,
			FirstName: "first",
			LastName:  "last",
			Email:     "email",
			ParentID:  "id",
		}
		body, _ := json.Marshal(account)

		req, err := http.NewRequest(http.MethodPost, "/attachAccount", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		mockService.On("AttachAccount", account).Return(errors.New("service error"))

		rr := httptest.NewRecorder()
		handler.HandleAttachAccount(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Failed to attach account")
	})

	t.Run("successful account attachment", func(t *testing.T) {
		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)
		account := model.AttachAccountBody{
			Role:      domain.PatientRole,
			FirstName: "first",
			LastName:  "last",
			Email:     "email",
			ParentID:  "id",
		}
		body, _ := json.Marshal(account)

		req, err := http.NewRequest(http.MethodPost, "/attachAccount", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		mockService.On("AttachAccount", account).Return(nil)

		rr := httptest.NewRecorder()
		handler.HandleAttachAccount(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
