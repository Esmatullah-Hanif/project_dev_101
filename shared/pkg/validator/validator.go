package validator

import (
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}

func ValidatePassword(password string) bool {
	return len(password) >= 8
}

func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

func IsEmptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
