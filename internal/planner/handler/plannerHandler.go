package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/arosace/WellnessWaveApi/internal/planner/model"
	"github.com/arosace/WellnessWaveApi/internal/planner/service"
	"github.com/arosace/WellnessWaveApi/pkg/utils"
)

type PlannerHandler struct {
	plannerService service.PlannerService
}

func NewPlannerHandler(plannerService service.PlannerService) *PlannerHandler {
	return &PlannerHandler{
		plannerService: plannerService,
	}
}

func (h *PlannerHandler) HandleAddMeal(w http.ResponseWriter, r *http.Request) {
	res := utils.GenericHttpResponse{}
	if r.Method != http.MethodPost {
		res.Error = "Method Not Allowed"
		utils.FormatResponse(w, res, http.StatusMethodNotAllowed)
		return
	}

	var meal model.Meal
	if err := json.NewDecoder(r.Body).Decode(&meal); err != nil {
		res.Error = "wrong_data_type"
		utils.FormatResponse(w, res, http.StatusBadRequest)
		return
	}

	if err := meal.ValidateModel(); err != nil {
		res.Error = err.Error()
		utils.FormatResponse(w, res, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := h.plannerService.AddMeal(ctx, &meal)
	if err != nil && strings.Contains(err.Error(), fmt.Sprintf("%d", http.StatusFound)) {
		res.Message = "A meal with the same name already exists"
		utils.FormatResponse(w, res, http.StatusFound)
		return
	}

	if err != nil {
		res.Error = fmt.Sprintf("Failed to add meal due to: %v", err)
		utils.FormatResponse(w, res, http.StatusInternalServerError)
		return
	}

	res.Data = meal
	utils.FormatResponse(w, res, http.StatusCreated)
}

func (h *PlannerHandler) HandleAddMealPlan(w http.ResponseWriter, r *http.Request) {
	res := utils.GenericHttpResponse{}
	if r.Method != http.MethodPost {
		res.Error = "Method Not Allowed"
		utils.FormatResponse(w, res, http.StatusMethodNotAllowed)
		return
	}

	var plan model.Plan
	if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
		res.Error = "wrong_data_type"
		utils.FormatResponse(w, res, http.StatusBadRequest)
		return
	}

	if err := plan.ValidateModel(); err != nil {
		res.Error = err.Error()
		utils.FormatResponse(w, res, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := h.plannerService.AddPlan(ctx, &plan)
	if err != nil {
		res.Error = fmt.Sprintf("Failed to add plan due to: %v", err)
		utils.FormatResponse(w, res, http.StatusInternalServerError)
		return
	}

	res.Data = plan
	utils.FormatResponse(w, res, http.StatusCreated)
}
