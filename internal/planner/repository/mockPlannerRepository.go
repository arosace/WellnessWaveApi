package repository

import (
	"fmt"
	"sync"

	"github.com/arosace/WellnessWaveApi/internal/planner/model"
)

var mealCount, planCount, dailyPlanCount = 0, 0, 0

const layout = "2006-01-02--15:04"

// MockUserRepository is a mock implementation of UserRepository that stores user data in memory.
type MockPlannerRepository struct {
	plans          map[string]*model.Plan
	meals          map[string]*model.Meal
	dailyPlans     map[string]*model.DailyPlan
	dailyPlanMeals map[string][]string
	mux            sync.RWMutex // ensures thread-safe access
}

// NewMockUserRepository creates a new instance of MockUserRepository.
func NewMockPlannerRepository() *MockPlannerRepository {
	return &MockPlannerRepository{
		plans:          make(map[string]*model.Plan),
		meals:          make(map[string]*model.Meal),
		dailyPlans:     make(map[string]*model.DailyPlan),
		dailyPlanMeals: make(map[string][]string),
	}
}

func (r *MockPlannerRepository) AddMeal(meal *model.Meal) error {
	meal.ID = fmt.Sprintf("%d", mealCount)
	r.meals[meal.Name] = meal
	mealCount += 1
	return nil
}

func (r *MockPlannerRepository) GetMealByNameAndHealthSpecialistId(mealName string, healthSpecialistId string) (*model.Meal, error) {
	meal, ok := r.meals[mealName]
	if ok && meal.HealthSpecialistId == healthSpecialistId {
		return meal, nil
	}

	return nil, nil
}

func (r *MockPlannerRepository) GetMealById(mealId string) (*model.Meal, error) {
	for _, m := range r.meals {
		if m.ID == mealId {
			return m, nil
		}
	}
	return nil, nil
}

func (r *MockPlannerRepository) GetMealsByHealthSpecialistId(healthSpecialistId string) ([]*model.Meal, error) {
	var meals []*model.Meal
	for _, m := range r.meals {
		if m.HealthSpecialistId == healthSpecialistId {
			meals = append(meals, m)
		}
	}
	return meals, nil
}

func (r *MockPlannerRepository) AddPlan(plan *model.Plan) error {
	plan.Id = fmt.Sprintf("%d", planCount)
	r.plans[plan.Id] = plan
	planCount += 1
	return nil
}

func (r *MockPlannerRepository) GetPlanByPatientId(patientId string) (*model.Plan, error) {
	for _, p := range r.plans {
		if p.PatientId == patientId {
			return p, nil
		}
	}
	return nil, nil
}

func (r *MockPlannerRepository) GetMealPlansByHealthSpecialistId(healthSpecialistId string) ([]*model.Plan, error) {
	var plans []*model.Plan
	for _, p := range r.plans {
		if p.HealthSpecialistId == healthSpecialistId {
			plans = append(plans, p)
		}
	}
	return plans, nil
}

func (r *MockPlannerRepository) AddDailyPlan(dailyPlan *model.DailyPlan) error {
	dailyPlan.Id = fmt.Sprintf("%d", dailyPlanCount)
	r.dailyPlans[dailyPlan.Id] = dailyPlan
	dailyPlanCount += 1
	return nil
}

func (r *MockPlannerRepository) MapMealToDailyPlan(mealId string, dailyPlanId string) error {
	r.dailyPlanMeals[dailyPlanId] = append(r.dailyPlanMeals[dailyPlanId], mealId)
	return nil
}
