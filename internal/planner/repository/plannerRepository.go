package repository

import "github.com/arosace/WellnessWaveApi/internal/planner/model"

type PlannerRepository interface {
	AddMeal(*model.Meal) error
	GetMealByNameAndHealthSpecialistId(string, string) (*model.Meal, error)
	GetMealById(string) (*model.Meal, error)
	GetMealsByHealthSpecialistId(string) ([]*model.Meal, error)
	AddPlan(*model.Plan) error
	AddDailyPlan(*model.DailyPlan) error
	MapMealToDailyPlan(string, string) error
}
