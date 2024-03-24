package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arosace/WellnessWaveApi/internal/event/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockEventService struct {
	mock.Mock
}

func (m *mockEventService) ScheduleEvent(ctx context.Context, event model.Event) (*model.Event, error) {
	args := m.Called(event)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Event), nil
}

func (m *mockEventService) GetEventsByHealthSpecialistId(ctx context.Context, healthSpecialistId string) ([]*model.Event, error) {
	args := m.Called(healthSpecialistId)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Event), nil
}

func (m *mockEventService) GetEventsByPatientId(ctx context.Context, id string) ([]*model.Event, error) {
	args := m.Called(id)
	return args.Get(0).([]*model.Event), args.Error(1)
}

func (m *mockEventService) RescheduleEvent(ctx context.Context, rescheduleRequest model.RescheduleRequest) (*model.Event, error) {
	args := m.Called(rescheduleRequest)
	return args.Get(0).(*model.Event), args.Error(1)
}

func TestHandleScheduleEvent(t *testing.T) {
	mockEventService := &mockEventService{}
	handler := &EventHandler{eventService: mockEventService}

	t.Run("when method is not POST, return Method Not Allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		handler.HandleScheduleEvent(rr, req)
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("when Invalid JSONm retrun bad request", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString("invalid json"))
		rr := httptest.NewRecorder()
		handler.HandleScheduleEvent(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("when Model Validation Failure, return bad request", func(t *testing.T) {
		event := model.Event{}
		body, _ := json.Marshal(event)
		req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		handler.HandleScheduleEvent(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockEventService.AssertExpectations(t)
	})

	t.Run("when Scheduling Failure, return system error", func(t *testing.T) {
		event := model.Event{
			HealthSpecialistID: "1",
			PatientID:          "2",
			EventType:          "whatever",
			EventDate:          "12/12/2012",
		}
		body, _ := json.Marshal(event)
		req, _ := http.NewRequest(http.MethodPost, "/events/schedule", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		mockEventService.On("ScheduleEvent", mock.Anything).Return(nil, errors.New("scheduling error")).Once()

		handler.HandleScheduleEvent(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockEventService.AssertExpectations(t)
	})

	t.Run("when Successful Event Scheduling, return event", func(t *testing.T) {
		event := model.Event{
			HealthSpecialistID: "1",
			PatientID:          "2",
			EventType:          "whatever",
			EventDate:          "12/12/2012",
		}
		body, _ := json.Marshal(event)
		req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		mockEventService.On("ScheduleEvent", mock.Anything).Return(&event, nil).Once()

		handler.HandleScheduleEvent(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
		mockEventService.AssertExpectations(t)
	})
}

func TestHandleGetEventsById(t *testing.T) {
	mockService := &mockEventService{}
	handler := &EventHandler{eventService: mockService}

	t.Run("HTTP Method Not Allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/", nil)
		rr := httptest.NewRecorder()
		handler.HandleGetEvents(rr, req)
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("Missing required query parameter", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		handler.HandleGetEvents(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Handle healthSpecialistId", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/?healthSpecialistId=123", nil)
		rr := httptest.NewRecorder()

		mockService.On("GetEventsByHealthSpecialistId", "123").Return([]*model.Event{}, nil)
		handler.HandleGetEvents(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Handle PatientId", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/?patientId=123", nil)
		rr := httptest.NewRecorder()

		mockService.On("GetEventsByPatientId", "123").Return([]*model.Event{}, nil)
		handler.HandleGetEvents(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Failed to retrieve events", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/?healthSpecialistId=1", nil)
		rr := httptest.NewRecorder()

		mockService.On("GetEventsByHealthSpecialistId", "1").Return(nil, errors.New("database error")).Once()

		handler.HandleGetEvents(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Successful event retrieval", func(t *testing.T) {
		events := []*model.Event{{
			HealthSpecialistID: "1",
			PatientID:          "2",
			EventType:          "whatever",
			EventDate:          "12/12/2012",
		}}
		req, _ := http.NewRequest(http.MethodGet, "/?healthSpecialistId=1", nil)
		rr := httptest.NewRecorder()

		mockService.On("GetEventsByHealthSpecialistId", "1").Return(events, nil).Once()

		handler.HandleGetEvents(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var returnedEvents []*model.Event
		err := json.NewDecoder(rr.Body).Decode(&returnedEvents)
		assert.NoError(t, err)
		assert.Equal(t, len(returnedEvents), 1)
		mockService.AssertExpectations(t)
	})
}

func TestHandleRescheduleEvent(t *testing.T) {
	mockService := &mockEventService{}
	handler := &EventHandler{eventService: mockService}

	t.Run("HTTP Method Not Allowed", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil) // Using GET instead of PATCH
		rr := httptest.NewRecorder()
		handler.HandleRescheduleEvent(rr, req)
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/", bytes.NewBufferString("invalid json"))
		rr := httptest.NewRecorder()
		handler.HandleRescheduleEvent(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Model Validation Failure", func(t *testing.T) {
		body := bytes.NewBufferString(`{"eventID":"", "newDate":""}`)
		req, _ := http.NewRequest(http.MethodPatch, "/", body)
		rr := httptest.NewRecorder()
		handler.HandleRescheduleEvent(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Failed to reschedule event", func(t *testing.T) {
		validRequestBody := model.RescheduleRequest{EventID: "123", NewDate: "2023-01-01T15:00:00Z"} // Example of a valid request body
		requestBodyBytes, _ := json.Marshal(validRequestBody)
		req, _ := http.NewRequest(http.MethodPatch, "/", bytes.NewBuffer(requestBodyBytes))
		rr := httptest.NewRecorder()

		mockService.On("RescheduleEvent", validRequestBody).Return(&model.Event{}, errors.New("service error")).Once()

		handler.HandleRescheduleEvent(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Successful Rescheduling", func(t *testing.T) {
		validRequestBody := model.RescheduleRequest{EventID: "123", NewDate: "2023-01-01T15:00:00Z"}
		requestBodyBytes, _ := json.Marshal(validRequestBody)
		req, _ := http.NewRequest(http.MethodPatch, "/", bytes.NewBuffer(requestBodyBytes))
		rr := httptest.NewRecorder()

		mockService.On("RescheduleEvent", validRequestBody).Return(&model.Event{}, nil).Once()

		handler.HandleRescheduleEvent(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		mockService.AssertExpectations(t)
	})
}
