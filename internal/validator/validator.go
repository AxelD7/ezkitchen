package validator

import (
	"net/mail"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key string, message string) {

	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

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
