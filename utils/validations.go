package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) []string {
	var validationErrors []string
	for _, fieldErr := range err.(validator.ValidationErrors) {
		switch fieldErr.Tag() {
		case "min":
			validationErrors = append(validationErrors, fmt.Sprintf("'%s' field must be at least %s characters long", fieldErr.Field(), fieldErr.Param()))
		case "max":
			validationErrors = append(validationErrors, fmt.Sprintf("'%s' field cannot exceed %s characters", fieldErr.Field(), fieldErr.Param()))
		case "required":
			validationErrors = append(validationErrors, fmt.Sprintf("'%s' field is required", fieldErr.Field()))
		case "url":
			validationErrors = append(validationErrors, fmt.Sprintf("'%s' field must be a valid URL", fieldErr.Field()))
		case "oneof":
			validationErrors = append(validationErrors, fmt.Sprintf("'%s' field must be one of %s", fieldErr.Field(), fieldErr.Param()))
		default:
			validationErrors = append(validationErrors, fmt.Sprintf("Validation failed for '%s' field with '%s' tag", fieldErr.Field(), fieldErr.Tag()))
		}
	}
	return validationErrors
}
