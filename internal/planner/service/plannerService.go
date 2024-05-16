package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/arosace/WellnessWaveApi/internal/planner/model"
	"github.com/arosace/WellnessWaveApi/internal/planner/repository"
	"github.com/labstack/echo/v5"
)

type PlannerService interface {
	AddMeal(echo.Context, *model.Meal) error
	GetMealById(echo.Context, string) (*model.Meal, error)
	GetMealsByHealthSpecialistId(echo.Context, string) ([]*model.Meal, error)
	AddPlan(echo.Context, *model.Plan) error
	GetMealPlanByPatientId(echo.Context, string) (*model.Plan, error)
	GetMealPlansByHealthSpecialistId(echo.Context, string) ([]string, error)
}

type plannerService struct {
	plannerRepository repository.PlannerRepository
}

func NewEventService(plannerRepo repository.PlannerRepository) PlannerService {
	return &plannerService{
		plannerRepository: plannerRepo,
	}
}

func (s *plannerService) AddMeal(ctx echo.Context, meal *model.Meal) error {
	m, err := s.plannerRepository.GetMealByNameAndHealthSpecialistId(ctx, meal.Name, meal.HealthSpecialistId)
	if err != nil {
		return err
	}
	if m != nil {
		return fmt.Errorf("%d", http.StatusFound)
	}
	s.plannerRepository.AddMeal(ctx, meal)
	return nil
}

func (s *plannerService) GetMealById(ctx echo.Context, mealId string) (*model.Meal, error) {
	m, err := s.plannerRepository.GetMealById(ctx, mealId)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *plannerService) GetMealsByHealthSpecialistId(ctx echo.Context, healthSpecialistId string) ([]*model.Meal, error) {
	m, err := s.plannerRepository.GetMealsByHealthSpecialistId(ctx, healthSpecialistId)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *plannerService) AddPlan(ctx echo.Context, plan *model.Plan) error {
	p, err := s.plannerRepository.GetPlanByPatientId(ctx, plan.PatientId)
	if err != nil {
		return err
	}
	if p != nil {
		return fmt.Errorf("%d - The patient has already a meal plan assigned", http.StatusFound)
	}

	//store plan in DB
	err = s.plannerRepository.AddPlan(ctx, plan)
	if err != nil {
		return err
	}

	//store each dailyPlan in DB
	for _, dailyPlan := range plan.DailyPlan {
		dailyPlan.PlanId = plan.Id
		dailyPlan.HealthSpecialistId = plan.HealthSpecialistId
		err := s.plannerRepository.AddDailyPlan(ctx, dailyPlan)
		if err != nil {
			return err
		}
		for _, meal := range dailyPlan.Meals {
			// store each meal in DB if it does not already exist for the specific health specialist
			meal.HealthSpecialistId = plan.HealthSpecialistId
			err := s.plannerRepository.AddMeal(ctx, meal)
			if err != nil {
				if strings.Contains(err.Error(), fmt.Sprintf("%d", http.StatusFound)) {
					continue
				}
				return err
			}
			// store mapping of meal to daily plan
			err = s.plannerRepository.MapMealToDailyPlan(ctx, meal.ID, dailyPlan.Id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *plannerService) GetMealPlanByPatientId(ctx echo.Context, patientId string) (*model.Plan, error) {
	plan, err := s.plannerRepository.GetPlanByPatientId(ctx, patientId)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (s *plannerService) GetMealPlansByHealthSpecialistId(ctx echo.Context, healthSpecialistId string) ([]string, error) {
	plans, err := s.plannerRepository.GetMealPlansByHealthSpecialistId(ctx, healthSpecialistId)
	if err != nil {
		return nil, err
	}
	var res []string
	for _, plan := range plans {
		res = append(res, plan.Id)
	}
	return res, nil
}
