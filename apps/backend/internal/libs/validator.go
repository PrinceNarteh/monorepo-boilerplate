// Package libs provides utility functions and types for the application.
package libs

import (
	"fmt"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

// ValidateStruct validates a struct and returns a map of field names to error messages.
// The field names are in lowercase, and the error messages are user-friendly.
// If the struct is valid, it returns nil.
func ValidateStruct(data any) map[string]string {
	validate := validator.New(validator.WithRequiredStructEnabled())
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

// getErrorMessage returns a user-friendly error message based on the validation error tag.
// It handles common validation tags and provides a meaningful message for each.
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
	case "lte":
		return fmt.Sprintf("%s must be %s or less", err.Field(), err.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", err.Field(), err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param())
	case "url":
		return fmt.Sprintf("%s is not a valid URL", err.Field())
	case "uuid":
		return fmt.Sprintf("%s is not a valid UUID", err.Field())
	case "datetime":
		return fmt.Sprintf("%s is not a valid datetime", err.Field())
	case "time":
		return fmt.Sprintf("%s is not a valid time", err.Field())
	case "ip":
		return fmt.Sprintf("%s is not a valid IP address", err.Field())
	case "cidr":
		return fmt.Sprintf("%s is not a valid CIDR notation", err.Field())
	case "json":
		return fmt.Sprintf("%s is not a valid JSON", err.Field())
	case "numeric":
		return fmt.Sprintf("%s must be a numeric value", err.Field())
	case "bool":
		return fmt.Sprintf("%s must be a boolean value", err.Field())
	case "base64":
		return fmt.Sprintf("%s is not a valid base64 string", err.Field())
	case "hexadecimal":
		return fmt.Sprintf("%s is not a valid hexadecimal string", err.Field())
	case "ascii":
		return fmt.Sprintf("%s must be a valid ASCII string", err.Field())
	case "printascii":
		return fmt.Sprintf("%s must be a printable ASCII string", err.Field())
	case "alphanum":
		return fmt.Sprintf("%s must be alphanumeric", err.Field())
	case "alpha":
		return fmt.Sprintf("%s must contain only alphabetic characters", err.Field())
	case "lowercase":
		return fmt.Sprintf("%s must be lowercase", err.Field())
	case "uppercase":
		return fmt.Sprintf("%s must be uppercase", err.Field())
	case "startswith":
		return fmt.Sprintf("%s must start with %s", err.Field(), err.Param())
	case "endswith":
		return fmt.Sprintf("%s must end with %s", err.Field(), err.Param())
	case "contains":
		return fmt.Sprintf("%s must contain %s", err.Field(), err.Param())
	case "excludes":
		return fmt.Sprintf("%s must not contain %s", err.Field(), err.Param())
	case "unique":
		return fmt.Sprintf("%s must be unique", err.Field())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", err.Field(), err.Param())
	case "eq":
		return fmt.Sprintf("%s must be equal to %s", err.Field(), err.Param())
	case "ne":
		return fmt.Sprintf("%s must not be equal to %s", err.Field(), err.Param())
	case "in":
		return fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param())
	case "notin":
		return fmt.Sprintf("%s must not be one of: %s", err.Field(), err.Param())
	case "required_if":
		return fmt.Sprintf("%s is required when %s", err.Field(), err.Param())
	case "required_unless":
		return fmt.Sprintf("%s is required unless %s", err.Field(), err.Param())
	case "required_with":
		return fmt.Sprintf("%s is required when %s is present", err.Field(), err.Param())
	case "required_without":
		return fmt.Sprintf("%s is required when %s is not present", err.Field(), err.Param())
	case "required_with_all":
		return fmt.Sprintf("%s is required when all of %s are present", err.Field(), err.Param())
	case "required_without_all":
		return fmt.Sprintf("%s is required when none of %s are present", err.Field(), err.Param())
	case "required_if_exists":
		return fmt.Sprintf("%s is required if %s exists", err.Field(), err.Param())
	case "required_if_not_exists":
		return fmt.Sprintf("%s is required if %s does not exist", err.Field(), err.Param())
	case "required_if_empty":
		return fmt.Sprintf("%s is required if %s is empty", err.Field(), err.Param())
	case "required_if_not_empty":
		return fmt.Sprintf("%s is required if %s is not empty", err.Field(), err.Param())
	case "required_if_true":
		return fmt.Sprintf("%s is required if %s is true", err.Field(), err.Param())
	case "required_if_false":
		return fmt.Sprintf("%s is required if %s is false", err.Field(), err.Param())
	case "required_if_null":
		return fmt.Sprintf("%s is required if %s is null", err.Field(), err.Param())
	case "required_if_not_null":
		return fmt.Sprintf("%s is required if %s is not null", err.Field(), err.Param())
	case "required_if_empty_string":
		return fmt.Sprintf("%s is required if %s is an empty string", err.Field(), err.Param())
	case "required_if_not_empty_string":
		return fmt.Sprintf("%s is required if %s is not an empty string", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}
