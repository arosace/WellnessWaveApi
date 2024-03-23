package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arosace/WellnessWaveApi/internal/event/model"
	"github.com/arosace/WellnessWaveApi/internal/event/service"
)

type EventHandler struct {
	eventService service.EventService
}

func NewEventHandler(eventService service.EventService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}

func (h *EventHandler) HandleScheduleEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var event model.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "wrong_data_type", http.StatusBadRequest)
		return
	}

	if err := event.ValidateModel(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	if _, err := h.eventService.ScheduleEvent(ctx, event); err != nil {
		http.Error(w, fmt.Sprintf("Failed to add event due to: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *EventHandler) HandleGetEventsByHealthSpecialistId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	healthSpecialistId := r.URL.Query().Get("healthSpecialistId")
	if healthSpecialistId == "" {
		http.Error(w, "healthSpecialistId is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	events, err := h.eventService.GetEventsByHealthSpecialistId(ctx, healthSpecialistId)
	if err != nil {
		http.Error(w, "Failed to retrieve events", http.StatusInternalServerError)
		return
	}

	if len(events) == 0 {
		http.Error(w, fmt.Sprintf("no events found for health specialist [%s]", healthSpecialistId), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
