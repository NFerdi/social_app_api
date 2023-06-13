package util

import (
	"github.com/go-playground/validator/v10"
)

func Validate(dto interface{}) []string {
	var errors []string

	validate := validator.New()
	if err := validate.Struct(dto); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field()+" is "+err.Tag())
		}
	}

	return errors
}
