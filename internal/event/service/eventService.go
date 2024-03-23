package service

import (
	"context"

	"github.com/arosace/WellnessWaveApi/internal/event/model"
	"github.com/arosace/WellnessWaveApi/internal/event/repository"
)

type EventService interface {
	ScheduleEvent(context.Context, model.Event) (*model.Event, error)
	GetEventsByHealthSpecialistId(context.Context, string) ([]*model.Event, error)
}

type eventService struct {
	eventRepository repository.EventRepository
}

func NewEventService(eventRepo repository.EventRepository) EventService {
	return &eventService{
		eventRepository: eventRepo,
	}
}

func (e *eventService) ScheduleEvent(ctx context.Context, event model.Event) (*model.Event, error) {
	return e.eventRepository.Add(ctx, event)
}

func (e *eventService) GetEventsByHealthSpecialistId(ctx context.Context, healthSpecialistId string) ([]*model.Event, error) {
	return e.eventRepository.GetByHealthSpecialistId(ctx, healthSpecialistId)
}
