package validator

import (
	"net/mail"
	"strings"
	"unicode/utf8"
)

// Validator stores field-specific validation errors.
type Validator struct {
	FieldErrors map[string]string
}

// Valid returns true if there are no validation errors recorded.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldError adds a new error message for the given field key
// if one does not already exist in the FieldErrors map.
func (v *Validator) AddFieldError(key string, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField adds an error message to a field if the provided condition is false.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank returns true if the provided string is not empty or whitespace.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars returns true if the string does not exceed the specified character limit.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// IsValidEmail returns true if the provided string is a valid email address.
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// GreaterThanN returns true if the value is greater than the given minimum.
// Works with int, float32, and float64 comparisons.
func GreaterThanN(value any, min any) bool {
	switch val := value.(type) {
	case int:
		if m, ok := min.(int); ok {
			return val > m
		}
	case float32:
		if m, ok := min.(float32); ok {
			return val > m
		}
	case float64:
		if m, ok := min.(float64); ok {
			return val > m
		}
	}
	return false
}

// LessThanN returns true if the value is less than the given minimum.
// Works with int, float32, and float64 comparisons.
func LessThanN(value any, min any) bool {
	switch val := value.(type) {
	case int:
		if m, ok := min.(int); ok {
			return val < m
		}
	case float32:
		if m, ok := min.(float32); ok {
			return val < m
		}
	case float64:
		if m, ok := min.(float64); ok {
			return val < m
		}
	}
	return false
}
