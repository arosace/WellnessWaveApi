package model

import (
	"fmt"
	"strings"
)

type Exercise struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Reps               int    `json:"reps"`
	Sets               int    `json:"sets"`
	HealthSpecialistId string `json:"health_specialist_id"`
	Description        string `json:"description"`
	Type               string `json:"type"`
}

func (e Exercise) ValidateModel() error {
	var errorStrings []string
	if e.Name == "" {
		errorStrings = append(errorStrings, "name")
	}

	if e.HealthSpecialistId == "" {
		errorStrings = append(errorStrings, "health_specialist_id")
	}

	if len(errorStrings) > 0 {
		return fmt.Errorf("missing_data: %s", strings.Join(errorStrings, ", "))
	}

	return nil
}
