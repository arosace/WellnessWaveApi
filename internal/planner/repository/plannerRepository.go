package repository

import (
	"github.com/arosace/WellnessWaveApi/internal/planner/model"
	"github.com/labstack/echo/v5"
)

type PlannerRepository interface {
	AddMeal(echo.Context, *model.Meal) error
	GetMealByNameAndHealthSpecialistId(echo.Context, string, string) (*model.Meal, error)
	GetMealById(echo.Context, string) (*model.Meal, error)
	GetMealsByHealthSpecialistId(echo.Context, string) ([]*model.Meal, error)
	AddPlan(echo.Context, *model.Plan) error
	GetPlanByPatientId(echo.Context, string) (*model.Plan, error)
	GetMealPlansByHealthSpecialistId(echo.Context, string) ([]*model.Plan, error)
	AddDailyPlan(echo.Context, *model.DailyPlan) error
	MapMealToDailyPlan(echo.Context, string, string) error
}
