package validations

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func UniversalValidation(body interface{}) (bool, error) {
	if err := validate.Struct(body); err != nil {
		return false, err
	}
	return true, nil
}
