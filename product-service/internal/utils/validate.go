package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct validates a struct based on its validation tags
func ValidateStruct(s interface{}) []string {
	var validationErrors []string
	if err := validate.Struct(s); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("%s is %s", err.Field(), err.Tag()))
		}
	}
	return validationErrors
}
