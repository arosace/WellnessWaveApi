package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/arosace/WellnessWaveApi/internal/account/domain"
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

func (m *Account) ValidateModel() error {
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

	if m.Role != domain.HhealthSpecialistRole && m.Role != domain.PatientRole {
		return errors.New("invalid_role")
	}

	return nil
}
