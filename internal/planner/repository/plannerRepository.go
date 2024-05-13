package repository

import (
	"github.com/arosace/WellnessWaveApi/internal/planner/model"
	"github.com/labstack/echo/v5"
)

type PlannerRepository interface {
	AddMeal(echo.Context, *model.Meal) error
	GetMealByNameAndHealthSpecialistId(echo.Context, string, string) (*model.Meal, error)
	GetMealById(echo.Context, string) (*model.Meal, error)
	GetMealsByHealthSpecialistId(string) ([]*model.Meal, error)
	AddPlan(*model.Plan) error
	GetPlanByPatientId(string) (*model.Plan, error)
	GetMealPlansByHealthSpecialistId(string) ([]*model.Plan, error)
	AddDailyPlan(*model.DailyPlan) error
	MapMealToDailyPlan(string, string) error
}
