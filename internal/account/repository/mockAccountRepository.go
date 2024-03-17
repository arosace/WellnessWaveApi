package repository

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"sync"

	"github.com/arosace/WellnessWaveApi/internal/account/model"
)

// MockUserRepository is a mock implementation of UserRepository that stores user data in memory.
type MockAccountRepository struct {
	accounts map[string]model.Account
	mux      sync.RWMutex // ensures thread-safe access
}

// NewMockUserRepository creates a new instance of MockUserRepository.
func NewMockAccountRepository() *MockAccountRepository {
	return &MockAccountRepository{
		accounts: make(map[string]model.Account),
	}
}

// AddUser adds a new user to the repository.
func (r *MockAccountRepository) Add(ctx context.Context, user model.Account) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	// Generate a random integer
	randomInteger := rand.Intn(1000000) // Generates a number in [0, 1000000)

	// Convert the integer to a string
	randomIntegerStr := strconv.Itoa(randomInteger)

	user.ID = randomIntegerStr
	r.accounts[user.Email] = user
	return nil
}

func (r *MockAccountRepository) List(ctx context.Context) ([]model.Account, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	var accounts []model.Account
	for _, a := range r.accounts {
		accounts = append(accounts, a)
	}

	return accounts, nil
}

// FindByID returns a user by their ID.
func (r *MockAccountRepository) FindByID(ctx context.Context, id string) (*model.Account, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	user, exists := r.accounts[id]
	if !exists {
		return nil, errors.New("account not found")
	}

	return &user, nil
}

// FindByEmail returns a user by their email.
func (r *MockAccountRepository) FindByEmail(ctx context.Context, email string) (*model.Account, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	user, exists := r.accounts[email]
	if !exists {
		return nil, errors.New("account not found")
	}

	return &user, nil
}
