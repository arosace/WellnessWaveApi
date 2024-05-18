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
	ParentID  string `json:"parent_id"`
	Role      string `json:"role"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"account_password"`
	AuthKey   string `json:"auth_key"`
	Username  string `json:"username"`
}

type VerifyAccount struct {
	JWT string `json:"jwt"`
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
		errorStrings = append(errorStrings, "account_password")
	}

	if len(errorStrings) > 0 {
		return fmt.Errorf("missing_data: %s", strings.Join(errorStrings, ", "))
	}

	if m.Role != domain.HhealthSpecialistRole && m.Role != domain.PatientRole {
		return errors.New("invalid_role")
	}

	return nil
}

func (m *Account) ValidateModelForInfoUpdate() error {
	if m.FirstName == "" && m.LastName == "" && m.ParentID == "" {
		return errors.New("no data to update provided")
	}
	if m.ID == "" {
		return errors.New("account id is missing")
	}
	return nil
}

func (m *Account) ValidateModelForAuthUpdate() error {
	if m.Password == "" && m.Email == "" {
		return errors.New("no data to update provided")
	}
	if m.ID == "" {
		return errors.New("account id is missing")
	}
	return nil
}
