package repository

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/arosace/WellnessWaveApi/internal/planner/domain"
	"github.com/arosace/WellnessWaveApi/internal/planner/model"
	"github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type PlannerRepo struct {
	Dao *daos.Dao
}

type PlannerRepository interface {
	// Meals
	AddMeal(echo.Context, *model.Meal) (*models.Record, error)
	GetMealByNameAndHealthSpecialistId(echo.Context, string, string) (*models.Record, error)
	GetMealById(echo.Context, string) (*models.Record, error)
	GetMealsByHealthSpecialistId(echo.Context, string) ([]*models.Record, error)
	// Meal Plans
	AddPlan(echo.Context, *model.Plan) (*models.Record, error)
	GetPlanByPatientId(echo.Context, string) (*models.Record, error)
	GetMealPlansByHealthSpecialistId(echo.Context, string) ([]*models.Record, error)
	// Daily Plans
	AddDailyPlan(echo.Context, *model.DailyPlan) (*models.Record, error)
	// Meal Map
	MapMealToPlan(echo.Context, model.MealMap) (*models.Record, error)

	//Transaction Queries
	AddPlanInTransaction(echo.Context, *model.Plan) (*models.Record, error)
}

func NewPlannerRepository(dao *daos.Dao) *PlannerRepo {
	return &PlannerRepo{
		Dao: dao,
	}
}

func (r *PlannerRepo) AddMeal(ctx echo.Context, meal *model.Meal) (*models.Record, error) {
	collection, err := r.Dao.FindCollectionByNameOrId(domain.MEALS_TABLENAME)
	if err != nil {
		return nil, fmt.Errorf("there was an error retrieving meals collection: %w", err)
	}

	record := models.NewRecord(collection)
	utils.LoadFromStruct(record, &meal)
	if err := r.Dao.SaveRecord(record); err != nil {
		return nil, fmt.Errorf("Failed to save meal: %w", err)
	}

	return record, nil
}

func (r *PlannerRepo) GetMealByNameAndHealthSpecialistId(ctx echo.Context, mealName string, healthSpecialistId string) (*models.Record, error) {
	params := dbx.Params{
		"health_specialist_id": healthSpecialistId,
		"name":                 mealName,
	}
	filter := "name = {:name} && health_specialist_id = {:health_specialist_id}"

	record, err := r.Dao.FindFirstRecordByFilter(
		domain.MEALS_TABLENAME,
		filter,
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("there was an error retrieving meal [%s] for specialist [%s]: %w", mealName, healthSpecialistId, err)
	}

	return record, nil
}

func (r *PlannerRepo) GetMealById(ctx echo.Context, mealId string) (*models.Record, error) {
	record, err := r.Dao.FindRecordById(domain.MEALS_TABLENAME, mealId)
	if err != nil {
		return nil, fmt.Errorf("there was an error retrieving meal [%s]: %w", mealId, err)
	}
	return record, nil
}

func (r *PlannerRepo) GetMealsByHealthSpecialistId(ctx echo.Context, healthSpecialistId string) ([]*models.Record, error) {
	params := dbx.Params{
		"health_specialist_id": healthSpecialistId,
	}
	filter := "health_specialist_id = {:health_specialist_id}"
	records, err := r.Dao.FindRecordsByFilter(
		domain.MEALS_TABLENAME,
		filter,
		"name",
		-1,
		0,
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("there was an error retrieving meals for specialist [%s]: %w", healthSpecialistId, err)
	}

	return records, nil
}

func (r *PlannerRepo) AddPlan(ctx echo.Context, plan *model.Plan) (*models.Record, error) {
	collection, err := r.Dao.FindCollectionByNameOrId(domain.PLANS_TABLENAME)
	if err != nil {
		return nil, fmt.Errorf("there was an error retrieving meal plans collection: %w", err)
	}

	record := models.NewRecord(collection)
	utils.LoadFromStruct(record, &plan)
	if err := r.Dao.SaveRecord(record); err != nil {
		return nil, fmt.Errorf("Failed to save meal plan: %w", err)
	}

	return record, nil
}

func (r *PlannerRepo) GetPlanByPatientId(ctx echo.Context, patientId string) (*models.Record, error) {
	record, err := r.Dao.FindFirstRecordByFilter(
		domain.PLANS_TABLENAME,
		"patient_id = {:patient_id}",
		dbx.Params{
			"patient_id": patientId,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("there was an error retrieving meal plan for patient [%s]: %w", patientId, err)
	}
	return record, nil
}

func (r *PlannerRepo) GetMealPlansByHealthSpecialistId(ctx echo.Context, healthSpecialistId string) ([]*models.Record, error) {
	records, err := r.Dao.FindRecordsByFilter(
		domain.PLANS_TABLENAME,
		"health_specialist_id = {:health_specialist_id}",
		"",
		-1,
		0,
		dbx.Params{
			"health_specialist_id": healthSpecialistId,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("there was an error retrieving meal plan for specialist [%s]: %w", healthSpecialistId, err)
	}
	return records, nil
}

func (r *PlannerRepo) AddDailyPlan(ctx echo.Context, dailyPlan *model.DailyPlan) (*models.Record, error) {
	collection, err := r.Dao.FindCollectionByNameOrId(domain.DAILY_PLANS_TABLENAME)
	if err != nil {
		return nil, fmt.Errorf("there was an error retrieving meal daily plans collection: %w", err)
	}

	record := models.NewRecord(collection)
	utils.LoadFromStruct(record, &dailyPlan)
	if err := r.Dao.SaveRecord(record); err != nil {
		return nil, fmt.Errorf("Failed to save daily meal plan: %w", err)
	}

	return record, nil
}

func (r *PlannerRepo) MapMealToPlan(ctx echo.Context, model model.MealMap) (*models.Record, error) {
	collection, err := r.Dao.FindCollectionByNameOrId(domain.MEAL_MAP)
	if err != nil {
		return nil, fmt.Errorf("there was an error retrieving meal map collection: %w", err)
	}

	record := models.NewRecord(collection)
	utils.LoadFromStruct(record, &model)
	if err := r.Dao.SaveRecord(record); err != nil {
		return nil, fmt.Errorf("Failed to save daily meal plan: %w", err)
	}

	return record, nil
}

func (r *PlannerRepo) AddPlanInTransaction(ctx echo.Context, plan *model.Plan) (*models.Record, error) {
	return nil, r.Dao.RunInTransaction(func(txDao *daos.Dao) error {
		oldDao := r.Dao
		r.Dao = txDao

		// check if plan exists
		p, err := r.GetPlanByPatientId(ctx, plan.PatientId)
		if err != nil {
			if !utils.IsErrorNotFound(err) {
				return err
			}
		}
		if p != nil {
			return fmt.Errorf("the patient [%s] already has a meal plan", plan.PatientId)
		}

		// create plan record
		planRecord, err := r.AddPlan(ctx, plan)
		if err != nil {
			return err
		}

		mealMap := model.MealMap{}
		for _, dailyPlan := range plan.DailyPlan {
			dailyPlan.PlanId = planRecord.Id
			dailyPlan.HealthSpecialistId = plan.HealthSpecialistId
			dailyPlanRecord, err := r.AddDailyPlan(ctx, dailyPlan)
			if err != nil {
				return err
			}
			for _, meal := range dailyPlan.Meals {
				// store each meal in DB if it does not already exist for the specific health specialist
				meal.HealthSpecialistId = plan.HealthSpecialistId
				mealRecord, err := r.AddMeal(ctx, meal)
				if err != nil {
					if strings.Contains(err.Error(), fmt.Sprintf("%d", http.StatusFound)) {
						continue
					}
					return err
				}
				// store mapping of meal to daily plan
				mealMap.MealId = mealRecord.Id
				mealMap.DailyPlanId = dailyPlanRecord.Id
				mealMap.PlanId = planRecord.Id
				_, err = r.MapMealToPlan(ctx, mealMap)
				if err != nil {
					return err
				}
			}
		}

		r.Dao = oldDao
		return nil
	})
}
