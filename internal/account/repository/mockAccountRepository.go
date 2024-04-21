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
	accounts map[string]*model.Account
	mux      sync.RWMutex // ensures thread-safe access
}

// NewMockUserRepository creates a new instance of MockUserRepository.
func NewMockAccountRepository() *MockAccountRepository {
	return &MockAccountRepository{
		accounts: make(map[string]*model.Account),
	}
}

func (r *MockAccountRepository) Add(ctx context.Context, user model.Account) (*model.Account, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	// Generate a random integer
	randomInteger := rand.Intn(1000000) // Generates a number in [0, 1000000)

	// Convert the integer to a string
	randomIntegerStr := strconv.Itoa(randomInteger)

	user.ID = randomIntegerStr
	r.accounts[user.Email] = &user
	return &user, nil
}

func (r *MockAccountRepository) List(ctx context.Context) ([]*model.Account, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	var accounts []*model.Account
	for _, a := range r.accounts {
		accounts = append(accounts, a)
	}

	return accounts, nil
}

// FindByID returns a user by their ID.
func (r *MockAccountRepository) FindByID(ctx context.Context, id string) (*model.Account, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	var user *model.Account
	for _, a := range r.accounts {
		if a.ID == id {
			user = a
			return user, nil
		}
	}

	return nil, errors.New("not_found")
}

// FindByEmail returns a user by their email.
func (r *MockAccountRepository) FindByEmail(ctx context.Context, email string) (*model.Account, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	user, exists := r.accounts[email]
	if !exists {
		return nil, errors.New("not_found")
	}

	return user, nil
}

func (r *MockAccountRepository) Update(ctx context.Context, user *model.Account) (*model.Account, error) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.accounts[user.Email] = user
	return user, nil
}

func (r *MockAccountRepository) UpdateAuth(ctx context.Context, user *model.Account) (*model.Account, error) {
	r.mux.Lock()
	defer r.mux.Unlock()
	if _, emailHasNotChanged := r.accounts[user.Email]; emailHasNotChanged {
		r.accounts[user.Email] = user
	} else {
		for _, a := range r.accounts {
			if a.ID == user.ID {
				delete(r.accounts, a.Email)
			}
		}

		r.accounts[user.Email] = user
	}

	return user, nil
}

func (r *MockAccountRepository) FindByParentID(ctx context.Context, parentId string) ([]*model.Account, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	accounts := make([]*model.Account, 0)
	for _, a := range r.accounts {
		if a.ParentID == parentId {
			accounts = append(accounts, a)
		}
	}

	return accounts, nil
}
