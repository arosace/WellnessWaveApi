package service

import (
	"github.com/arosace/WellnessWaveApi/internal/event/model"
	"github.com/arosace/WellnessWaveApi/internal/event/repository"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/models"
)

type EventService interface {
	ScheduleEvent(echo.Context, model.Event) (*models.Record, error)
	GetEventsByHealthSpecialistId(echo.Context, string, string) ([]*models.Record, error)
	GetEventsByPatientId(echo.Context, string, string) ([]*models.Record, error)
	RescheduleEvent(echo.Context, model.RescheduleRequest) (*models.Record, error)
}

type eventService struct {
	eventRepository repository.EventRepository
}

func NewEventService(eventRepo repository.EventRepository) EventService {
	return &eventService{
		eventRepository: eventRepo,
	}
}

func (e *eventService) ScheduleEvent(ctx echo.Context, event model.Event) (*models.Record, error) {
	return e.eventRepository.Add(ctx, event)
}

func (e *eventService) GetEventsByHealthSpecialistId(ctx echo.Context, healthSpecialistId string, after string) ([]*models.Record, error) {
	return e.eventRepository.GetByHealthSpecialistId(ctx, healthSpecialistId, after)
}

func (e *eventService) GetEventsByPatientId(ctx echo.Context, patientId string, after string) ([]*models.Record, error) {
	return e.eventRepository.GetByPatientId(ctx, patientId, after)
}

func (e *eventService) RescheduleEvent(ctx echo.Context, rescheduleRequest model.RescheduleRequest) (*models.Record, error) {
	return e.eventRepository.Update(ctx, rescheduleRequest)
}
