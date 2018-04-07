package helpers

import (
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

var validate *validator.Validate

func GetValidator() *validator.Validate {
	return validate
}

func InitValidator()  {
	validate = validator.New()

	validate.RegisterValidation("lowercase", isLowercase)
}

func isLowercase(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	return s == strings.ToLower(s)
}
