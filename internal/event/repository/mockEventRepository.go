package repository

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/arosace/WellnessWaveApi/internal/event/model"
)

const layout = "2006-01-02--15:04"

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

func (r *MockEventRepository) GetByHealthSpecialistId(ctx context.Context, healthSpecialistId string, after string) ([]*model.Event, error) {
	var events []*model.Event

	var afterDate time.Time
	var err error
	if after != "" {
		afterDate, err = time.Parse(layout, after)
		if err != nil {
			fmt.Printf("Error parsing date 1: %v\n", err)
			return nil, err
		}
	}

	for _, event := range r.events {
		if event.HealthSpecialistID == healthSpecialistId {
			if after != "" {
				eventDate, err := time.Parse(layout, event.EventDate)
				if err != nil {
					fmt.Printf("Error parsing date 1: %v\n", err)
					return nil, err
				}
				if eventDate.After(afterDate) {
					events = append(events, event)
				}
			} else {
				events = append(events, event)
			}

		}
	}
	return events, nil
}

func (r *MockEventRepository) GetByPatientId(ctx context.Context, patientId string, after string) ([]*model.Event, error) {
	var events []*model.Event

	var afterDate time.Time
	var err error
	if after != "" {
		afterDate, err = time.Parse(layout, after)
		if err != nil {
			fmt.Printf("Error parsing date 1: %v\n", err)
			return nil, err
		}
	}

	for _, event := range r.events {
		if event.PatientID == patientId {
			if after != "" {
				eventDate, err := time.Parse(layout, event.EventDate)
				if err != nil {
					fmt.Printf("Error parsing date 1: %v\n", err)
					return nil, err
				}
				if eventDate.After(afterDate) {
					events = append(events, event)
				}
			} else {
				events = append(events, event)
			}

		}
	}
	return events, nil
}

func (r *MockEventRepository) Update(ctx context.Context, rescheduleRequest model.RescheduleRequest) (*model.Event, error) {
	var event *model.Event
	for _, event := range r.events {
		if event.ID == rescheduleRequest.EventID {
			event.EventDate = rescheduleRequest.NewDate
		}
	}
	return event, nil
}
