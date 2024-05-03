package util

import (
	"github.com/go-playground/validator/v10"
)

var ValidUsername validator.Func = func(fl validator.FieldLevel) bool {
	if username, ok := fl.Field().Interface().(string); ok {
		if err := validateUsername(username); err != nil {
			return false
		}
		return true
	}
	return false
}

var ValidPassword validator.Func = func(fl validator.FieldLevel) bool {
	if password, ok := fl.Field().Interface().(string); ok {
		if err := validatePassword(password); err != nil {
			return false
		}
		return true
	}
	return false
}

var ValidFullname validator.Func = func(fl validator.FieldLevel) bool {
	if fullname, ok := fl.Field().Interface().(string); ok {
		if err := validateFullname(fullname); err != nil {
			return false
		}
		return true
	}
	return false
}
