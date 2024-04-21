package planner

import (
	"github.com/arosace/WellnessWaveApi/internal/planner/handler"
	"github.com/arosace/WellnessWaveApi/internal/planner/repository"
	"github.com/arosace/WellnessWaveApi/internal/planner/service"
	"github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/gorilla/mux"
)

type PlannerService struct {
	Router         *mux.Router
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
	s.Router.HandleFunc("/planner/addMeal", utils.HttpMiddleware(s.ServiceHandler.HandleAddMeal))
	s.Router.HandleFunc("/planner/addMealPlan", utils.HttpMiddleware(s.ServiceHandler.HandleAddMealPlan))
	s.Router.HandleFunc("/planner/getMeal", utils.HttpMiddleware(s.ServiceHandler.HandleGetMeal))
	s.Router.HandleFunc("/planner/getMealPlan", utils.HttpMiddleware(s.ServiceHandler.HandleGetMealPlan)) //can provide both health specialist id and patient id
	//s.Router.HandleFunc("/planner/getDailyMealPlan", utils.HttpMiddleware(s.ServiceHandler.HandleGetDailyMealPlan)) //provide planId
}
