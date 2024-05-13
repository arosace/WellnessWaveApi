package planner

import (
	"github.com/arosace/WellnessWaveApi/internal/planner/handler"
	"github.com/arosace/WellnessWaveApi/internal/planner/repository"
	"github.com/arosace/WellnessWaveApi/internal/planner/service"
	"github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type PlannerService struct {
	App            *pocketbase.PocketBase
	ServiceHandler *handler.PlannerHandler
}

func (s PlannerService) Init() {
	plannerRepo := repository.NewMockPlannerRepository()
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
	//s.Router.HandleFunc("/planner/addMealPlan", utils.HttpMiddleware(s.ServiceHandler.HandleAddMealPlan))
	//s.Router.HandleFunc("/planner/getMeal", utils.HttpMiddleware(s.ServiceHandler.HandleGetMeal))
	//s.Router.HandleFunc("/planner/getMealPlan", utils.HttpMiddleware(s.ServiceHandler.HandleGetMealPlan))
}
