package model

import (
	"fmt"
	"strings"
)

type Plan struct {
	Id                 string       `json:"id"`
	DailyPlan          []*DailyPlan `json:"daily_plans"`
	HealthSpecialistId string       `json:"health_specialist_id"`
	PatientId          string       `json:"patient_id"`
}

type DailyPlan struct {
	Id                 string  `json:"id"`
	PlanId             string  `json:"plan_id"`
	DayIndex           int     `json:"day_index"`
	Meals              []*Meal `json:"meals"`
	HealthSpecialistId string  `json:"health_specialist_id"`
}

func (p *Plan) ValidateModel() error {
	var errorStrings []string
	if p.HealthSpecialistId == "" {
		errorStrings = append(errorStrings, "health_specialist_id")
	}

	if p.PatientId == "" {
		errorStrings = append(errorStrings, "patient_id")
	}

	if len(p.DailyPlan) == 0 {
		errorStrings = append(errorStrings, "daily_plan")
	} else {
		var dpErr error
		for _, dp := range p.DailyPlan {
			dpErr = dp.ValidateModel()
		}
		if dpErr != nil {
			errorStrings = append(errorStrings, dpErr.Error())
		}
	}

	if len(errorStrings) > 0 {
		return fmt.Errorf("missing_data: %s", strings.Join(errorStrings, ", "))
	}

	return nil
}

func (dp *DailyPlan) ValidateModel() error {
	if len(dp.Meals) == 0 {
		return fmt.Errorf("daily_plan_meals")
	} else {
		var mErr error
		for _, m := range dp.Meals {
			mErr = m.ValidateModel()
		}
		if mErr != nil {
			return mErr
		}
	}

	return nil
}
