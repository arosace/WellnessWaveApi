package model

import (
	"fmt"
	"strings"
)

// Account represents a user in the system.
type Account struct {
	ID        string `json:"id,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Password  string `json:"password"`
	ParentID  string `json:"parent_id"`
}

func (m *Account) ValidateHealthSpecialist() error {
	var errorStrings []string

	if m.FirstName == "" {
		errorStrings = append(errorStrings, "first_name")
	}
	if m.LastName == "" {
		errorStrings = append(errorStrings, "last_name")
	}
	if m.Email == "" {
		errorStrings = append(errorStrings, "email")
	}
	if m.Role == "" {
		errorStrings = append(errorStrings, "role")
	}
	if m.Password == "" {
		errorStrings = append(errorStrings, "password")
	}

	if len(errorStrings) > 0 {
		return fmt.Errorf("missing_data: %s", strings.Join(errorStrings, ", "))
	}

	return nil
}

func (m *Account) ValidatePatient(alreadyExists bool) error {
	var errorStrings []string

	if m.FirstName == "" {
		errorStrings = append(errorStrings, "first_name")
	}
	if m.LastName == "" {
		errorStrings = append(errorStrings, "last_name")
	}
	if m.Email == "" {
		errorStrings = append(errorStrings, "email")
	}
	if m.Role == "" {
		errorStrings = append(errorStrings, "role")
	}

	if m.ParentID == "" {
		errorStrings = append(errorStrings, "parent_id")
	}

	if !alreadyExists {
		if m.Password == "" {
			errorStrings = append(errorStrings, "password")
		}
	}

	if len(errorStrings) > 0 {
		return fmt.Errorf("missing_data: %s", strings.Join(errorStrings, ", "))
	}

	return nil
}
