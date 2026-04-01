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
	err := v.validator.RegisterValidation("alphanum_hyphen", func(fl go_validator.FieldLevel) bool {
		value := fl.Field().String()
		// Slug should be lowercase, alphanumeric with hyphens
		match, err := regexp.MatchString(`^[a-z0-9-]+$`, value)
		if err != nil {
			return false
		}
		return match
	})
	if err != nil {
		panic(err)
	}

	// Register slug validator (more strict)
	err = v.validator.RegisterValidation("slug", func(fl go_validator.FieldLevel) bool {
		value := fl.Field().String()
		// Slug: lowercase, starts with letter, alphanumeric with hyphens, no consecutive hyphens
		match, err := regexp.MatchString(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`, value)
		if err != nil {
			return false
		}
		return match
	})
	if err != nil {
		panic(err)
	}

	// Register date_range validator
	err = v.validator.RegisterValidation("date_range", func(_ /* fl */ go_validator.FieldLevel) bool {
		// This is a placeholder - actual date range validation should be done at service layer
		// where we have access to both start and end dates
		return true
	})
	if err != nil {
		panic(err)
	}

	// Register enum validator for specific values
	err = v.validator.RegisterValidation("role_enum", func(fl go_validator.FieldLevel) bool {
		value := fl.Field().String()
		validRoles := map[string]bool{
			"admin":  true,
			"mentor": true,
			"mentee": true,
		}
		return validRoles[value]
	})
	if err != nil {
		panic(err)
	}

	// Register difficulty enum validator
	err = v.validator.RegisterValidation("difficulty_enum", func(fl go_validator.FieldLevel) bool {
		value := fl.Field().String()
		validDifficulties := map[string]bool{
			"easy":   true,
			"medium": true,
			"hard":   true,
		}
		return validDifficulties[value]
	})
	if err != nil {
		panic(err)
	}

	// Register status enum validator
	err = v.validator.RegisterValidation("status_enum", func(fl go_validator.FieldLevel) bool {
		value := fl.Field().String()
		validStatuses := map[string]bool{
			"pending":   true,
			"attempted": true,
			"completed": true,
		}
		return validStatuses[value]
	})
	if err != nil {
		panic(err)
	}

	// Register password complexity validator
	err = v.validator.RegisterValidation("password_complexity", func(fl go_validator.FieldLevel) bool {
		value := fl.Field().String()
		hasLetter := false
		hasNumber := false

		for _, char := range value {
			if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
				hasLetter = true
			}
			if char >= '0' && char <= '9' {
				hasNumber = true
			}
			if hasLetter && hasNumber {
				return true
			}
		}

		return false
	})
	if err != nil {
		panic(err)
	}
}

// you can register your custom validation
// func (v *validator) RegisterValidation(tag string, fn go_validator.Func) error {
// 	return v.validator.RegisterValidation(tag, fn)
// }
