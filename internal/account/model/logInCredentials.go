package model

import (
	"fmt"
	"strings"
)

type LogInCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *LogInCredentials) ValidateModel() error {
	var errorList []string
	if c.Email == "" {
		errorList = append(errorList, "email")
	}
	if c.Password == "" {
		errorList = append(errorList, "password")
	}
	if len(errorList) > 0 {
		return fmt.Errorf("missing_parameters:%s", strings.Join(errorList, ","))
	}
	return nil
}
