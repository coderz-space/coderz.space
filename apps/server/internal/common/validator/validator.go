package validator

import go_validator "github.com/go-playground/validator/v10"

// wrapper for go-validator
type validator struct {
	validator *go_validator.Validate
}

func NewValidator() *validator {
	return &validator{
		validator: go_validator.New(),
	}
}

func (v *validator) ValidateStruct(s interface{}) error {
	return v.validator.Struct(s)
}
func (v *validator) ValidateField(field interface{}, tag string) error {
	return v.validator.Var(field, tag)
}


// you can register your custom validation
// func (v *validator) RegisterValidation(tag string, fn go_validator.Func) error {
// 	return v.validator.RegisterValidation(tag, fn)
// }
