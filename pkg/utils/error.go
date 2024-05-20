package utils

import "strings"

func IsErrorNotFound(err error) bool {
	return strings.Contains(err.Error(), "no rows in result set")
}
