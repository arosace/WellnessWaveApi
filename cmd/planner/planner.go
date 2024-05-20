package planner

import (
	"github.com/arosace/WellnessWaveApi/internal/planner/handler"
	"github.com/arosace/WellnessWaveApi/internal/planner/repository"
	"github.com/arosace/WellnessWaveApi/internal/planner/service"
	"github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
)

type PlannerService struct {
	App            *pocketbase.PocketBase
	Dao            *daos.Dao
	ServiceHandler *handler.PlannerHandler
}

func (s PlannerService) Init() {
	plannerRepo := repository.NewPlannerRepository(s.Dao)
	plannerService := service.NewEventService(plannerRepo)
	plannerServiceHandler := handler.NewPlannerHandler(plannerService)
	s.ServiceHandler = plannerServiceHandler
	s.RegisterEndpoints()
}

func (s PlannerService) RegisterEndpoints() {
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/api/planner/addMeal", s.ServiceHandler.HandleAddMeal, utils.EchoMiddleware)
		return nil
	})
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/api/planner/addMealPlan", s.ServiceHandler.HandleAddMealPlan, utils.EchoMiddleware)
		return nil
	})
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/api/planner/getMeal", s.ServiceHandler.HandleGetMeal, utils.EchoMiddleware)
		return nil
	})
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/api/planner/getMealPlan", s.ServiceHandler.HandleGetMealPlan, utils.EchoMiddleware)
		return nil
	})
}
