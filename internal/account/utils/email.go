package utils

import (
	"errors"
	"regexp"
	"strings"
)

func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	isMatched := emailRegex.MatchString(email)
	if !isMatched {
		return errors.New("invalid email")
	}
	return nil
}
func LowerCaseEmail(email string) string {
	return strings.ToLower(email)
}
