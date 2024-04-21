package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/arosace/WellnessWaveApi/internal/planner/model"
	"github.com/arosace/WellnessWaveApi/internal/planner/repository"
)

type PlannerService interface {
	AddMeal(context.Context, *model.Meal) error
	AddPlan(context.Context, *model.Plan) error
}

type plannerService struct {
	plannerRepository repository.PlannerRepository
}

func NewEventService(plannerRepo repository.PlannerRepository) PlannerService {
	return &plannerService{
		plannerRepository: plannerRepo,
	}
}

func (s *plannerService) AddMeal(ctx context.Context, meal *model.Meal) error {
	m, err := s.plannerRepository.GetMealByNameAndHealthSpecialistId(meal.Name, meal.HealthSpecialistId)
	if err != nil {
		return err
	}
	if m != nil {
		return fmt.Errorf("%d", http.StatusFound)
	}
	s.plannerRepository.AddMeal(meal)
	return nil
}

func (s *plannerService) AddPlan(ctx context.Context, plan *model.Plan) error {
	//store plan in DB
	err := s.plannerRepository.AddPlan(plan)
	if err != nil {
		return err
	}

	//store each dailyPlan in DB
	for _, dailyPlan := range plan.DailyPlan {
		dailyPlan.PlanId = plan.Id
		dailyPlan.HealthSpecialistId = plan.HealthSpecialistId
		err := s.plannerRepository.AddDailyPlan(dailyPlan)
		if err != nil {
			return err
		}
		for _, meal := range dailyPlan.Meals {
			// store each meal in DB if it does not already exist for the specific health specialist
			meal.HealthSpecialistId = plan.HealthSpecialistId
			err := s.plannerRepository.AddMeal(meal)
			if err != nil {
				if strings.Contains(err.Error(), fmt.Sprintf("%d", http.StatusFound)) {
					continue
				}
				return err
			}
			// store mapping of meal to daily plan
			err = s.plannerRepository.MapMealToDailyPlan(meal.ID, dailyPlan.Id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
