package handler

import (
	"fmt"
	"net/http"

	"github.com/arosace/WellnessWaveApi/internal/planner/model"
	"github.com/arosace/WellnessWaveApi/internal/planner/service"
	"github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
)

type PlannerHandler struct {
	plannerService service.PlannerService
}

func NewPlannerHandler(plannerService service.PlannerService) *PlannerHandler {
	return &PlannerHandler{
		plannerService: plannerService,
	}
}

func (h *PlannerHandler) HandleAddMeal(ctx echo.Context) error {
	res := utils.GenericHttpResponse{}
	var meal model.Meal
	if err := ctx.Bind(&meal); err != nil {
		res.Error = "wrong_data_type"
		return apis.NewBadRequestError(res.Error, nil)
	}

	if err := meal.ValidateModel(); err != nil {
		res.Error = err.Error()
		return apis.NewBadRequestError(res.Error, nil)
	}

	record, err := h.plannerService.AddMeal(ctx, &meal)
	if err != nil {
		if utils.IsErrorFound(err) {
			res.Message = "A meal with the same name already exists"
			return apis.NewBadRequestError(res.Message, res)
		}
		return apis.NewBadRequestError(fmt.Sprintf("there was an error creating a meal: %s", err.Error()), res)
	}

	res.Data = record
	return ctx.JSON(http.StatusCreated, res)
}

func (h *PlannerHandler) HandleAddMealPlan(ctx echo.Context) error {
	res := utils.GenericHttpResponse{}
	var plan model.Plan
	if err := ctx.Bind(&plan); err != nil {
		res.Error = "wrong_data_type"
		return apis.NewBadRequestError(res.Error, nil)
	}

	if err := plan.ValidateModel(); err != nil {
		res.Error = err.Error()
		return apis.NewBadRequestError(res.Error, nil)
	}

	err := h.plannerService.AddPlan(ctx, &plan)
	if err != nil {
		res.Error = fmt.Sprintf("Failed to add plan due to: %v", err)
		return apis.NewBadRequestError(res.Error, nil)
	}

	res.Data = plan
	return ctx.JSON(http.StatusCreated, res)
}

func (h *PlannerHandler) HandleGetMeal(ctx echo.Context) error {
	res := utils.GenericHttpResponse{}
	healthSpecialistId := ctx.QueryParam("healthSpecialistId")
	mealId := ctx.QueryParam("mealId")
	if healthSpecialistId == "" && mealId == "" {
		res.Error = "query parameters missing, either provide a HealthSpecialistId or a mealId"
		return apis.NewBadRequestError(res.Error, nil)
	}

	if mealId != "" {
		meal, err := h.plannerService.GetMealById(ctx, mealId)
		if err != nil {
			res.Error = fmt.Sprintf("There was an error retrieving meal by id: %s", err.Error())
			return apis.NewBadRequestError(res.Error, nil)
		}
		if meal == nil {
			res.Message = "meal not found for given id"
			return apis.NewApiError(http.StatusNotFound, res.Message, res)
		}
		res.Data = meal
		return ctx.JSON(http.StatusOK, res)
	}

	if healthSpecialistId != "" {
		meals, err := h.plannerService.GetMealsByHealthSpecialistId(ctx, healthSpecialistId)
		if err != nil {
			res.Error = fmt.Sprintf("There was an error retrieving meals by healt specialist id: %s", err.Error())
			return apis.NewBadRequestError(res.Error, nil)
		}
		if meals == nil {
			res.Message = "No meals found for given health specialist id"
			return apis.NewApiError(http.StatusNotFound, res.Message, res)
		}
		res.Data = meals
		return ctx.JSON(http.StatusOK, res)
	}

	res.Message = fmt.Sprintf("request correclty processed but no meal was found for healthSpecialistId [%s] and mealId [%s]", healthSpecialistId, mealId)
	return ctx.JSON(http.StatusOK, res)
}

func (h *PlannerHandler) HandleGetMealPlan(ctx echo.Context) error {
	res := utils.GenericHttpResponse{}

	healthSpecialistId := ctx.QueryParam("healthSpecialistId")
	patientId := ctx.QueryParam("patientId")
	if healthSpecialistId == "" && patientId == "" {
		res.Error = "query parameters missing, either provide a HealthSpecialistId or a patientId"
		return apis.NewBadRequestError(res.Error, nil)
	}

	if patientId != "" {
		plan, err := h.plannerService.GetMealPlanByPatientId(ctx, patientId)
		if err != nil {
			res.Error = fmt.Sprintf("There was an error retrieving meal plan by patient id: %s", err.Error())
			return apis.NewBadRequestError(res.Error, nil)
		}
		if plan == nil {
			res.Message = "meal plan not found for given patient id"
			return apis.NewApiError(http.StatusNotFound, res.Message, res)
		}
		res.Data = plan
		return ctx.JSON(http.StatusOK, res)
	}

	if healthSpecialistId != "" {
		plans, err := h.plannerService.GetMealPlansByHealthSpecialistId(ctx, healthSpecialistId)
		if err != nil {
			res.Error = fmt.Sprintf("There was an error retrieving meal plan by health specialist id id: %s", err.Error())
			return apis.NewBadRequestError(res.Error, nil)
		}
		if len(plans) == 0 {
			res.Message = "meal plans not found for given health specialist id"
			return apis.NewApiError(http.StatusNotFound, res.Message, res)
		}
		res.Data = plans
		return ctx.JSON(http.StatusOK, res)
	}

	res.Message = fmt.Sprintf("request correclty processed but no meal plan was found for healthSpecialistId [%s] and patientId [%s]", healthSpecialistId, patientId)
	return ctx.JSON(http.StatusOK, res)
}
