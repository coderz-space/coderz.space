package validator

import (
	"regexp"

	go_validator "github.com/go-playground/validator/v10"
)

// wrapper for go-validator
type validator struct {
	validator *go_validator.Validate
}

func NewValidator() *validator {
	v := &validator{
		validator: go_validator.New(),
	}

	// Register custom validators
	v.registerCustomValidators()

	return v
}

func (v *validator) ValidateStruct(s interface{}) error {
	return v.validator.Struct(s)
}

func (v *validator) ValidateField(field interface{}, tag string) error {
	return v.validator.Var(field, tag)
}

// registerCustomValidators registers all custom validation functions
func (v *validator) registerCustomValidators() {
	// Register alphanum_hyphen validator for slugs
	v.validator.RegisterValidation("alphanum_hyphen", func(fl go_validator.FieldLevel) bool {
		value := fl.Field().String()
		// Slug should be lowercase, alphanumeric with hyphens
		match, _ := regexp.MatchString(`^[a-z0-9-]+$`, value)
		return match
	})
}

// you can register your custom validation
// func (v *validator) RegisterValidation(tag string, fn go_validator.Func) error {
// 	return v.validator.RegisterValidation(tag, fn)
// }
