package validation

import (
	"fmt"
	"regexp"
	"strings"

	validator "github.com/go-playground/validator/v10"
	"github.com/terraforge-gg/terraforge/internal/models"
)

type Validator struct {
	Validator *validator.Validate
}

type ValidationError struct {
	Errors map[string]string `json:"errors"`
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

func (v *Validator) Validate(i any) error {
	if err := v.Validator.Struct(i); err != nil {
		errors := make(map[string]string)

		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			field = strings.ToLower(field[:1]) + field[1:]
			tag := err.Tag()

			switch tag {
			case "required":
				errors[field] = "Required"
			case "min":
				errors[field] = fmt.Sprintf("Must be longer than %s characters", err.Param())
			case "max":
				errors[field] = fmt.Sprintf("Must be shorter than %s characters", err.Param())
			case "url_slug":
				errors[field] = fmt.Sprintf("'%s' is not a valid slug.", err.Value())
			case "project_type":
				errors[field] = fmt.Sprintf("'%s' is not a valid project type.", err.Value())
			default:
				errors[field] = "Invalid"
			}
		}

		return &ValidationError{Errors: errors}
	}
	return nil
}

var SlugRegexValidator = regexp.MustCompile(`^[\p{L}\p{N}_-]+$`)

func ValidateUrlSlug(fl validator.FieldLevel) bool {
	return SlugRegexValidator.MatchString(fl.Field().String())
}

func ValidateProjectType(fl validator.FieldLevel) bool {
	switch models.ProjectType(fl.Field().String()) {
	case models.ProjectTypeMod:
		return true
	}
	return false
}
