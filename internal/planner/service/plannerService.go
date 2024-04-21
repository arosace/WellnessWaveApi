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
	GetMealById(context.Context, string) (*model.Meal, error)
	GetMealsByHealthSpecialistId(context.Context, string) ([]*model.Meal, error)
	AddPlan(context.Context, *model.Plan) error
	GetMealPlanByPatientId(context.Context, string) (*model.Plan, error)
	GetMealPlansByHealthSpecialistId(context.Context, string) ([]string, error)
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

func (s *plannerService) GetMealById(ctx context.Context, mealId string) (*model.Meal, error) {
	m, err := s.plannerRepository.GetMealById(mealId)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *plannerService) GetMealsByHealthSpecialistId(ctx context.Context, healthSpecialistId string) ([]*model.Meal, error) {
	m, err := s.plannerRepository.GetMealsByHealthSpecialistId(healthSpecialistId)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *plannerService) AddPlan(ctx context.Context, plan *model.Plan) error {
	p, err := s.plannerRepository.GetPlanByPatientId(plan.PatientId)
	if err != nil {
		return err
	}
	if p != nil {
		return fmt.Errorf("%d - The patient has already a meal plan assigned", http.StatusFound)
	}

	//store plan in DB
	err = s.plannerRepository.AddPlan(plan)
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

func (s *plannerService) GetMealPlanByPatientId(ctx context.Context, patientId string) (*model.Plan, error) {
	plan, err := s.plannerRepository.GetPlanByPatientId(patientId)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (s *plannerService) GetMealPlansByHealthSpecialistId(ctx context.Context, healthSpecialistId string) ([]string, error) {
	plans, err := s.plannerRepository.GetMealPlansByHealthSpecialistId(healthSpecialistId)
	if err != nil {
		return nil, err
	}
	var res []string
	for _, plan := range plans {
		res = append(res, plan.Id)
	}
	return res, nil
}
