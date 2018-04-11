package helpers

import (
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

//Used for validating data send to the api
var validate *validator.Validate

//Initialises the validator
func InitializeValidator()  {
	validate = validator.New()
	validate.RegisterValidation("lowercase", isLowercase)
}

//Returns the initialized Validator. Do not call before calling InitializeValidator!
func GetValidator() *validator.Validate {
	return validate
}

//Checks if the specified fieldLevel is lowercase
func isLowercase(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	return s == strings.ToLower(s)
}