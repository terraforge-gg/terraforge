package validation

import (
	"fmt"
	"net/url"
	"os"
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

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("url_slug", ValidateUrlSlug)
	validate.RegisterValidation("project_type", ValidateProjectType)
	validate.RegisterValidation("project_version_dependency_type", ValidateProjectDependencyType)
	validate.RegisterValidation("file_url", ValidateFileUrl)
	validate.RegisterValidation("semver", ValidateSemVer)

	return validate
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
			case "file_url":
				errors[field] = "invalid file url"
			case "project_version_dependency_type":
				errors[field] = "invalid project version dependency type"
			case "semver":
				errors[field] = "invalid semver"
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

var SemVerValidator = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-[a-zA-Z0-9.]+)?$`)

func ValidateSemVer(fl validator.FieldLevel) bool {
	return SemVerValidator.MatchString(fl.Field().String())
}

func ValidateProjectDependencyType(fl validator.FieldLevel) bool {
	switch models.ProjectReleaseDependencyType(fl.Field().String()) {
	case models.ProjectReleaseDependencyTypeRequired,
		models.ProjectReleaseDependencyTypeOptional:
		return true
	}
	return false
}

func ValidateFileUrl(fl validator.FieldLevel) bool {
	cdnUrl := os.Getenv("CDN_URL")
	parsedCdnUrl, err := url.Parse(cdnUrl)

	if err != nil {
		return false
	}

	urlStr := fl.Field().String()
	parsedURL, err := url.Parse(urlStr)

	if err != nil {
		return false
	}

	return strings.EqualFold(parsedURL.Host, parsedCdnUrl.Host)
}
