package repository

import (
	"context"
	"errors"
	"sync"

	"../model"
)

// MockUserRepository is a mock implementation of UserRepository that stores user data in memory.
type MockUserRepository struct {
	users map[string]model.User
	mux   sync.RWMutex // ensures thread-safe access
}

// NewMockUserRepository creates a new instance of MockUserRepository.
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]model.User),
	}
}

// AddUser adds a new user to the repository.
func (r *MockUserRepository) Add(ctx context.Context, user model.User) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return errors.New("user already exists")
	}

	r.users[user.ID] = user
	return nil
}

// FindByID returns a user by their ID.
func (r *MockUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
