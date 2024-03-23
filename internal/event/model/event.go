package model

import (
	"fmt"
	"strings"
)

type Event struct {
	ID                 string `json:"id"`
	HealthSpecialistID string `json:"health_specialist_id"`
	PatientID          string `json:"patient_id"`
	EventType          string `json:"event_type"`
	EventDescription   string `json:"event_description"`
	EventDate          string `json:"event_date"`
}

// ValidateModel validates the event data.
func (e *Event) ValidateModel() error {
	var errorStrings []string
	if e.HealthSpecialistID == "" {
		errorStrings = append(errorStrings, "health_specialist_id")
	}

	if e.PatientID == "" {
		errorStrings = append(errorStrings, "patient_id")
	}

	if e.EventType == "" {
		errorStrings = append(errorStrings, "event_type")
	}

	if e.EventDate == "" {
		errorStrings = append(errorStrings, "event_date")
	}

	if len(errorStrings) > 0 {
		return fmt.Errorf("missing_data: %s", strings.Join(errorStrings, ", "))
	}

	return nil
}
