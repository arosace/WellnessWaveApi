package service

import (
	"fmt"

	"github.com/arosace/WellnessWaveApi/internal/planner/model"
	"github.com/arosace/WellnessWaveApi/internal/planner/repository"
	"github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/models"
)

type PlannerService interface {
	AddMeal(echo.Context, *model.Meal) (*models.Record, error)
	GetMealById(echo.Context, string) (*models.Record, error)
	GetMealsByHealthSpecialistId(echo.Context, string) ([]*models.Record, error)
	AddPlan(echo.Context, *model.Plan) error
	GetMealPlanByPatientId(echo.Context, string) (*models.Record, error)
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

func (s *plannerService) AddMeal(ctx echo.Context, meal *model.Meal) (*models.Record, error) {
	record, err := s.plannerRepository.GetMealByNameAndHealthSpecialistId(ctx, meal.Name, meal.HealthSpecialistId)
	if err != nil {
		if !utils.IsErrorNotFound(err) {
			return nil, err
		}
	}
	if record != nil {
		return nil, utils.ErrorFound()
	}
	return s.plannerRepository.AddMeal(ctx, meal)
}

func (s *plannerService) GetMealById(ctx echo.Context, mealId string) (*models.Record, error) {
	m, err := s.plannerRepository.GetMealById(ctx, mealId)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *plannerService) GetMealsByHealthSpecialistId(ctx echo.Context, healthSpecialistId string) ([]*models.Record, error) {
	m, err := s.plannerRepository.GetMealsByHealthSpecialistId(ctx, healthSpecialistId)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *plannerService) AddPlan(ctx echo.Context, plan *model.Plan) error {
	_, err := s.plannerRepository.AddPlanInTransaction(ctx, plan)
	if err != nil {
		return fmt.Errorf("there was an error adding the meal plan in transaction: %w", err)
	}
	return nil
}

func (s *plannerService) GetMealPlanByPatientId(ctx echo.Context, patientId string) (*models.Record, error) {
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
