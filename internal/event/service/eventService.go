package service

import (
	"fmt"
	"time"

	"github.com/arosace/WellnessWaveApi/internal/event/domain"
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
	GetEventById(echo.Context, string) (*models.Record, error)
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

func (e *eventService) GetEventById(ctx echo.Context, eventId string) (*models.Record, error) {
	return e.eventRepository.GetById(ctx, eventId)
}

func (e *eventService) RescheduleEvent(ctx echo.Context, rescheduleRequest model.RescheduleRequest) (*models.Record, error) {
	record, err := e.eventRepository.GetById(ctx, rescheduleRequest.EventID)
	if err != nil {
		return nil, err
	}
	parsedTime, err := time.Parse(domain.Layout, rescheduleRequest.NewDate)
	if err != nil {
		return nil, fmt.Errorf("there was an error parsing the new date: %w", err)
	}
	isSame := record.GetDateTime("event_date").Time().Compare(parsedTime) == 0
	if isSame {
		return record, nil
	}

	record.Set("event_date", rescheduleRequest.NewDate)
	return e.eventRepository.Update(ctx, record)
}
