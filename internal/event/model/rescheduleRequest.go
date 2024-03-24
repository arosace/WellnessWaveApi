package model

import (
	"errors"
	"fmt"
)

type RescheduleRequest struct {
	EventID string `json:"event_id"`
	NewDate string `json:"date"`
}

func (r *RescheduleRequest) ValidateModel() error {
	var errorList []string
	if r.EventID == "" {
		errorList = append(errorList, "event_id")
	}
	if r.NewDate == "" {
		errorList = append(errorList, "date")
	}
	if len(errorList) > 0 {
		return errors.New(fmt.Sprintf("missing_parameters: %v", errorList))
	}
	return nil
}
