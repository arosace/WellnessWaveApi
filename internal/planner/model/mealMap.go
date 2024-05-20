package model

import (
	"fmt"
	"strings"
)

type MealMap struct {
	PlanId      string `json:"plan_id"`
	DailyPlanId string `json:"daily_plan_id"`
	MealId      string `json:"meal_id"`
}

func (m *MealMap) ValidateModel() error {
	var errorStrings []string

	if m.PlanId == "" {
		errorStrings = append(errorStrings, "plan_id")
	}
	if m.DailyPlanId == "" {
		errorStrings = append(errorStrings, "daily_plan_id")
	}
	if m.MealId == "" {
		errorStrings = append(errorStrings, "meal_id")
	}

	if len(errorStrings) > 0 {
		return fmt.Errorf("missing_data: %s", strings.Join(errorStrings, ", "))
	}

	return nil
}
