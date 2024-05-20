package event

import (
	"fmt"

	"github.com/arosace/WellnessWaveApi/internal/event/domain"
	"github.com/arosace/WellnessWaveApi/internal/event/handler"
	"github.com/arosace/WellnessWaveApi/internal/event/repository"
	"github.com/arosace/WellnessWaveApi/internal/event/service"
	"github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type EventService struct {
	App            *pocketbase.PocketBase
	Dao            *daos.Dao
	ServiceHandler *handler.EventHandler
}

func (s EventService) Init() {
	eventRepo := repository.NewEventRepository(s.Dao)
	eventService := service.NewEventService(eventRepo)
	accountServiceHandler := handler.NewEventHandler(eventService)
	s.ServiceHandler = accountServiceHandler
	s.RegisterEndpoints()
	s.RegisterHooks()
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

func (s EventService) RegisterHooks() {
	s.App.OnModelAfterCreate(domain.TABLENAME).Add(func(e *core.ModelEvent) error {
		record := e.Model.(*models.Record)
		// TODO: Send email to patient with calendar event and call link
		fmt.Println(record)
		return nil
	})
}
