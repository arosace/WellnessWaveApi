package repository

import (
	"context"
	"math/rand"
	"strconv"
	"sync"

	"github.com/arosace/WellnessWaveApi/internal/event/model"
)

// MockUserRepository is a mock implementation of UserRepository that stores user data in memory.
type MockEventRepository struct {
	events map[string]*model.Event
	mux    sync.RWMutex // ensures thread-safe access
}

// NewMockUserRepository creates a new instance of MockUserRepository.
func NewMockEventRepository() *MockEventRepository {
	return &MockEventRepository{
		events: make(map[string]*model.Event),
	}
}

func (r *MockEventRepository) Add(ctx context.Context, event model.Event) (*model.Event, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	// Generate a random integer
	randomInteger := rand.Intn(1000000) // Generates a number in [0, 1000000)

	// Convert the integer to a string
	randomIntegerStr := strconv.Itoa(randomInteger)

	event.ID = randomIntegerStr
	r.events[event.ID] = &event
	return &event, nil
}
