package model

import (
	"fmt"
	"strings"

	"github.com/arosace/WellnessWaveApi/internal/account/domain"
)

type AttachAccountBody struct {
	FirstName string
	LastName  string
	Role      string
	Email     string
	ParentID  string
}

func (m *AttachAccountBody) ValidateModel() error {
	var missingData []string

	if m.Email == "" {
		missingData = append(missingData, "email")
	}

	if m.ParentID == "" {
		missingData = append(missingData, "parent_id")
	}

	if m.LastName == "" {
		missingData = append(missingData, "last_name")
	}

	if m.FirstName == "" {
		missingData = append(missingData, "first_name")
	}

	if len(missingData) > 0 {
		return fmt.Errorf("missing_data: %s", strings.Join(missingData, ", "))
	}

	if m.Role != domain.PatientRole {
		return fmt.Errorf("invalid_role")
	}

	return nil
}
