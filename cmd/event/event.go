package event

import (
	"github.com/arosace/WellnessWaveApi/internal/event/handler"
	"github.com/arosace/WellnessWaveApi/internal/event/repository"
	"github.com/arosace/WellnessWaveApi/internal/event/service"
	"github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type EventService struct {
	App            *pocketbase.PocketBase
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
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/api/events/schedule", s.ServiceHandler.HandleScheduleEvent, utils.EchoMiddleware)
		return nil
	})
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/api/events", s.ServiceHandler.HandleGetEvents, utils.EchoMiddleware)
		return nil
	})
}
