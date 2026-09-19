package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/amirjbr/shared-expense/pkg/conv"
)

func HashPassword(pass string) string {
	h := sha256.New()
	h.Write(conv.ToBytes(pass))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func CheckPassword(hashedPassword, password string) bool {
	if HashPassword(password) == hashedPassword {
		return true
	}
	return false
}
func ValidatePasswordWithFeedback(password string) error {
	tests := []struct {
		pattern string
		message string
	}{
		{".{7,}", "Password must be at least 7 characters long"},
		{"[a-z]", "Password must contain at least one lowercase letter"},
		{"[A-Z]", "Password must contain at least one uppercase letter"},
		{"[0-9]", "Password must contain at least one digit"},
	}

	var errMessages []string

	for _, test := range tests {
		match, _ := regexp.MatchString(test.pattern, password)
		if !match {
			errMessages = append(errMessages, test.message)
		}
	}

	if len(errMessages) > 0 {
		feedback := strings.Join(errMessages, "\n")
		return errors.Join(errors.New("invalid password format"), fmt.Errorf(feedback))
	}

	return nil
}
