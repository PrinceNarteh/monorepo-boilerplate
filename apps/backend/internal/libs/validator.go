// Package libs provides utility functions and types for the application.
package libs

import (
	"fmt"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

// ValidateStruct validates a struct using the global validator.
// It returns an error if the struct is not valid.
var validate = validator.New(validator.WithRequiredStructEnabled())

// ValidateStruct validates a struct and returns a map of field names to error messages.
// The field names are in lowercase, and the error messages are user-friendly.
// If the struct is valid, it returns nil.
func ValidateStruct(data any) map[string]string {
	errors := make(map[string]string)

	err := validate.Struct(data)
	if valErrs, ok := err.(validator.ValidationErrors); ok {
		for _, v := range valErrs {
			errors[strings.ToLower(v.Field())] = getErrorMessage(v)
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s is not a valid email", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	case "gte":
		return fmt.Sprintf("%s must be %s or greater", err.Field(), err.Param())

	default:
		return ""
	}
}
