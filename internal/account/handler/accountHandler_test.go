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
	"github.com/gorilla/mux"

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

func (s *mockAccountService) UpdateAccount(ctx context.Context, account model.Account, infoType string) (*model.Account, error) {
	args := s.Called(account, infoType)
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

func (s *mockAccountService) GetAccountById(ctx context.Context, id string) (*model.Account, error) {
	args := s.Called(id)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), nil
}

func (s *mockAccountService) AttachAccount(ctx context.Context, accountToAttach model.AttachAccountBody) (*model.Account, error) {
	args := s.Called(accountToAttach)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), nil
}

func (s *mockAccountService) GetAccountByEmail(ctx context.Context, email string) (*model.Account, error) {
	args := s.Called(email)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), nil
}

func (s *mockAccountService) GetAttachedAccounts(ctx context.Context, parentId string) ([]*model.Account, error) {
	args := s.Called(parentId)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Account), nil
}

func (s *mockAccountService) Authorize(ctx context.Context, credentials model.LogInCredentials) (*model.Account, error) {
	args := s.Called(credentials)
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

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(userJSON))
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
			Role:      domain.HhealthSpecialistRole,
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
			Role:      domain.HhealthSpecialistRole,
			Password:  "[]",
		}
		userJSON, _ := json.Marshal(testUser)
		rr := httptest.NewRecorder()

		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)
		mockService.On("CheckAccountExists", testUser.Email).Return(false)
		mockService.On("AddAccount", testUser).Return(nil, errors.New("error"))

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(userJSON))
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

		mockService.On("AttachAccount", mock.Anything).Return(nil, errors.New("service error"))

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

		mockService.On("AttachAccount", mock.Anything).Return(&model.Account{}, nil)

		rr := httptest.NewRecorder()
		handler.HandleAttachAccount(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
	})
}

func TestHandleGetAttachedAccounts(t *testing.T) {
	t.Run("when non-GET method is used, return method not allowed", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/accounts/parent_id", nil) // Using POST instead of GET
		if err != nil {
			t.Fatal(err)
		}
		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)
		handler.HandleGetAccounts(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("when parentID is missing, return bad request code", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/accounts/", nil)
		if err != nil {
			t.Fatal(err)
		}
		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)
		handler.HandleGetAttachedAccounts(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run("when unexpected parameter is passed, return bad request code", func(t *testing.T) {
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		req, err := http.NewRequest("GET", "/accounts/parent_id", nil)
		if err != nil {
			t.Fatal(err)
		}
		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)
		router.HandleFunc("/accounts/{parent_id}", handler.HandleGetAttachedAccounts)
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("when service fails to get accounts, return internal server error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		req, err := http.NewRequest("GET", "/accounts/123", nil)
		if err != nil {
			t.Fatal(err)
		}
		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)

		mockService.On("GetAttachedAccounts", mock.Anything).Return(nil, errors.New("internal error"))

		router.HandleFunc("/accounts/{parent_id}", handler.HandleGetAttachedAccounts)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("when get request is successfull, return list of attached accounts", func(t *testing.T) {
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		req, err := http.NewRequest("GET", "/accounts/123", nil)
		if err != nil {
			t.Fatal(err)
		}
		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)

		mockService.On("GetAttachedAccounts", mock.Anything).Return([]*model.Account{
			{ID: "1"}, {ID: "1"},
		}, nil)

		router.HandleFunc("/accounts/{parent_id}", handler.HandleGetAttachedAccounts)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestHandleUpdateAccount(t *testing.T) {
	t.Run("when method is not PUT, return Method Not Allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/?type=personal", nil)
		rr := httptest.NewRecorder()
		mockService := &mockAccountService{}
		handler := NewAccountHandler(mockService)
		r := mux.NewRouter()
		r.HandleFunc("/", handler.HandleUpdateAccount)
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("when type url parameter is missing or invalid, return Bad Request", func(t *testing.T) {
		bodies := []string{"", "invalid"}
		for _, body := range bodies {
			req, _ := http.NewRequest(http.MethodPut, "/?type="+body, nil)
			rr := httptest.NewRecorder()
			mockService := &mockAccountService{}
			h := NewAccountHandler(mockService)
			r := mux.NewRouter()
			r.HandleFunc("/", h.HandleUpdateAccount)
			r.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("when body has invalid data format, return Bad Request", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/?type=personal", bytes.NewBufferString("invalid json"))
		rr := httptest.NewRecorder()
		mockService := &mockAccountService{}
		h := NewAccountHandler(mockService)
		r := mux.NewRouter()
		r.HandleFunc("/", h.HandleUpdateAccount)
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("when personal info update is successful, return OK", func(t *testing.T) {
		accountJSON := `{"email":"test@example.com","first_name":"Name","last_name":"Surname","id":"12"}`
		req, _ := http.NewRequest(http.MethodPut, "/?type=personal", bytes.NewBufferString(accountJSON))
		rr := httptest.NewRecorder()
		mockService := &mockAccountService{}
		h := NewAccountHandler(mockService)
		mockService.On("UpdateAccount", mock.AnythingOfType("model.Account"), "personal").Return(&model.Account{}, nil)

		r := mux.NewRouter()
		r.HandleFunc("/", h.HandleUpdateAccount)
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestHandleLogIn(t *testing.T) {
	t.Run("when method is not POST, return Method Not Allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		mockService := &mockAccountService{}
		h := NewAccountHandler(mockService)
		h.HandleLogIn(rr, req)
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("when body has invalid data format, return Bad Request", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString("invalid json"))
		rr := httptest.NewRecorder()
		mockService := &mockAccountService{}
		h := NewAccountHandler(mockService)
		h.HandleLogIn(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("when login credentials are missing, return Bad Request", func(t *testing.T) {
		credentials := model.LogInCredentials{Email: "", Password: ""} // Missing credentials
		body, _ := json.Marshal(credentials)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		mockService := &mockAccountService{}
		h := NewAccountHandler(mockService)
		h.HandleLogIn(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "missing_parameters:email,password")
	})

	t.Run("when authorization fails, return Unauthorized", func(t *testing.T) {
		credentials := model.LogInCredentials{Email: "user@example.com", Password: "wrongpassword"}
		body, _ := json.Marshal(credentials)
		req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		mockService := &mockAccountService{}
		h := NewAccountHandler(mockService)
		mockService.On("Authorize", credentials).Return(nil, errors.New("not_authorized"))

		h.HandleLogIn(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("when login is successful, return OK", func(t *testing.T) {
		credentials := model.LogInCredentials{Email: "user@example.com", Password: "correctpassword"}
		body, _ := json.Marshal(credentials)
		req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		mockService := &mockAccountService{}
		h := NewAccountHandler(mockService)
		mockService.On("Authorize", credentials).Return(nil)

		h.HandleLogIn(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockService.AssertExpectations(t)
	})
}
