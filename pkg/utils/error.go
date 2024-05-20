package utils

import (
	"errors"
	"strings"
)

func IsErrorNotFound(err error) bool {
	return strings.Contains(err.Error(), "no rows in result set")
}

func ErrorFound() error {
	return errors.New("record_found")
}

func IsErrorFound(err error) bool {
	return strings.Contains(err.Error(), "record_found")
}
