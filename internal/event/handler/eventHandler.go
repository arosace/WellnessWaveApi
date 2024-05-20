package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/arosace/WellnessWaveApi/internal/event/model"
	"github.com/arosace/WellnessWaveApi/internal/event/service"
	"github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
)

type EventHandler struct {
	eventService service.EventService
}

func NewEventHandler(eventService service.EventService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}

func (h *EventHandler) HandleScheduleEvent(ctx echo.Context) error {
	res := utils.GenericHttpResponse{}
	var event model.Event
	if err := ctx.Bind(&event); err != nil {
		return apis.NewBadRequestError("wrong_data_type", nil)
	}

	if err := event.ValidateModel(); err != nil {
		return apis.NewBadRequestError(err.Error(), nil)
	}

	scheduledEvent, err := h.eventService.ScheduleEvent(ctx, event)
	if err != nil {
		return apis.NewBadRequestError(fmt.Sprintf("Failed to add event due to: %v", err), nil)
	}

	res.Data = scheduledEvent
	return ctx.JSON(http.StatusCreated, res)
}

func (h *EventHandler) HandleGetEvents(ctx echo.Context) error {
	res := model.EventResponse{}
	var events []*models.Record
	var err error

	healthSpecialistId := ctx.QueryParam("healthSpecialistId")
	patientId := ctx.QueryParam("patientId")
	after := ctx.QueryParam("after")

	if healthSpecialistId != "" && patientId != "" {
		return apis.NewBadRequestError("Both healthSpecialistId and patientId were specified but only one of the two is expected", nil)
	}

	if healthSpecialistId != "" {
		events, err = h.getEventsByHealthSpecialistId(ctx, healthSpecialistId, after)
		if err != nil {
			return apis.NewBadRequestError(fmt.Sprintf("There was an error retrieving events by specialist id: %s", err.Error()), nil)
		}
	}

	if patientId != "" {
		events, err = h.getEventsByPatientId(ctx, patientId, after)
		if err != nil {
			return apis.NewBadRequestError(fmt.Sprintf("There was an error retrieving events by patient id: %s", err.Error()), nil)
		}
	}

	if events == nil {
		res.Data = []*model.Event{}
	} else {
		res.Data = events
	}
	return ctx.JSON(http.StatusOK, res)
}

func (h *EventHandler) HandleRescheduleEvent(ctx echo.Context) error {
	res := utils.GenericHttpResponse{}
	var rescheduleRequest model.RescheduleRequest
	if err := ctx.Bind(&rescheduleRequest); err != nil {
		return apis.NewBadRequestError("Invalid request body", nil)
	}

	if err := rescheduleRequest.ValidateModel(); err != nil {
		return apis.NewBadRequestError(err.Error(), nil)
	}

	rescheduledEvent, err := h.eventService.RescheduleEvent(ctx, rescheduleRequest)
	if err != nil {
		if utils.IsErrorNotFound(err) {
			return apis.NewApiError(http.StatusNotFound, "Failed to reschedule event: no event with that id was found", nil)
		}
		return apis.NewBadRequestError(fmt.Sprintf("Failed to reschedule event: %s", err.Error()), nil)
	}

	res.Data = rescheduledEvent
	return ctx.JSON(http.StatusOK, res)
}

func (h *EventHandler) getEventsByHealthSpecialistId(ctx echo.Context, id string, after string) ([]*models.Record, error) {
	events, err := h.eventService.GetEventsByHealthSpecialistId(ctx, id, after)
	if err != nil {
		return nil, errors.New("Failed to retrieve events")
	}
	return events, nil
}

func (h *EventHandler) getEventsByPatientId(ctx echo.Context, id string, after string) ([]*models.Record, error) {
	events, err := h.eventService.GetEventsByPatientId(ctx, id, after)
	if err != nil {
		return nil, errors.New("Failed to retrieve events")
	}
	return events, nil
}
