package event

import (
	"github.com/arosace/WellnessWaveApi/internal/event/handler"
	"github.com/arosace/WellnessWaveApi/internal/event/repository"
	"github.com/arosace/WellnessWaveApi/internal/event/service"
	"github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/gorilla/mux"
)

type EventService struct {
	Router         *mux.Router
	ServiceHandler *handler.EventHandler
}

func (s EventService) Init() {
	eventRepo := repository.NewMockEventRepository()
	eventService := service.NewEventService(eventRepo)
	accountServiceHandler := handler.NewEventHandler(eventService)
	s.ServiceHandler = accountServiceHandler
	s.RegisterEndpoints()
}

func (s EventService) RegisterEndpoints() {
	s.Router.HandleFunc("/events/schedule", utils.HttpMiddleware(s.ServiceHandler.HandleScheduleEvent))
	s.Router.HandleFunc("/events", utils.HttpMiddleware(s.ServiceHandler.HandleGetEvents))
}
