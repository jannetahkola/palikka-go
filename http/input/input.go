package input

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	usernameLengthMin = 8
	usernameLengthMax = 24
	usernamePattern   = "^[a-zA-Z0-9-_]+$"
	passwordLengthMin = 8
	passwordLengthMax = 72 // bcrypt limit is 72 bytes
)

const (
	errEmptyUsername = "username is empty"
	errEmptyPassword = "password is empty"
)

// Sanitize calls strings.TrimSpace on the pointer's value, or returns the pointer if its value is nil.
func Sanitize(ptr *string) *string {
	if ptr != nil {
		result := strings.TrimSpace(*ptr)
		return &result
	}
	return ptr
}

// ValidateUsername checks the pointer's value with username validation rules.
// Sanitize should always be called beforehand.
func ValidateUsername(ptr *string) error {
	if ptr == nil {
		return errors.New(errEmptyUsername)
	}

	username := *ptr
	l := len(username)

	if l == 0 {
		return errors.New(errEmptyUsername)
	}

	if l < usernameLengthMin || l > usernameLengthMax {
		return errors.New(fmt.Sprintf(
			"username is invalid, must be between %d and %d characters", usernameLengthMin, usernameLengthMax))
	}

	if match, err := regexp.Match(usernamePattern, []byte(username)); err != nil || !match {
		return errors.New(fmt.Sprintf("username is invalid, must match %s", usernamePattern))
	}

	return nil
}

// ValidatePassword checks the pointer's value with password validation rules.
// Sanitize should NOT be called beforehand.
func ValidatePassword(ptr *string) error {
	if ptr == nil {
		return errors.New(errEmptyPassword)
	}

	password := *ptr
	l := len(password)

	if l == 0 {
		return errors.New(errEmptyPassword)
	}

	if l < passwordLengthMin || l > passwordLengthMax {
		return errors.New(fmt.Sprintf(
			"password is invalid, must be between %d and %d characters", passwordLengthMin, passwordLengthMax))
	}

	return nil
}
