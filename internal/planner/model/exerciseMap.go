package model

import (
	"fmt"
	"strings"
)

type ExerciseMap struct {
	PlanId      string `json:"plan_id"`
	DailyPlanId string `json:"daily_plan_id"`
	ExerciseId  string `json:"exercise_id"`
}

func (m *ExerciseMap) ValidateModel() error {
	var errorStrings []string

	if m.PlanId == "" {
		errorStrings = append(errorStrings, "plan_id")
	}
	if m.DailyPlanId == "" {
		errorStrings = append(errorStrings, "daily_plan_id")
	}
	if m.ExerciseId == "" {
		errorStrings = append(errorStrings, "exercise_id")
	}

	if len(errorStrings) > 0 {
		return fmt.Errorf("missing_data: %s", strings.Join(errorStrings, ", "))
	}

	return nil
}
