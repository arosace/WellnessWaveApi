package model

import (
	"fmt"
	"strings"
)

type Meal struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Ingredients        []string `json:"ingredients"`
	HealthSpecialistId string   `json:"health_specialist_id"`
	Cals               int      `json:"cals"`
	Description        string   `json:"description"`
	Type               string   `json:"type"`
}

func (m *Meal) ValidateModel() error {
	var errorStrings []string
	if m.Name == "" {
		errorStrings = append(errorStrings, "name")
	}

	if m.HealthSpecialistId == "" {
		errorStrings = append(errorStrings, "health_specialist_id")
	}

	if len(errorStrings) > 0 {
		return fmt.Errorf("missing_data: %s", strings.Join(errorStrings, ", "))
	}

	return nil
}
