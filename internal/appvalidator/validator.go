package appvalidator

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InputValidator(input interface{}) map[string]string {
	validate = validator.New()
	err := validate.Struct(input)
	if err != nil {
		validationErrors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors[e.Field()] = fmt.Sprintf("%s failed validation", e.Field())
		}
		return validationErrors
	}
	return nil
}

func IsIDValid(ID string) error {
	id, err := strconv.ParseInt(ID, 10, 64)
	if err != nil || id < 1 {
		return errors.New("invalid id parameter")
	}
	return nil
}
