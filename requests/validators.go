package requests

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var ValidUsername validator.Func = func(fl validator.FieldLevel) bool {
	if username, ok := fl.Field().Interface().(string); ok {
		if len(username) < 4 || len(username) > 64 {
			return false
		}

		match, _ := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9_]*$", username)
		return match
	}
	return false
}

var ValidPassword validator.Func = func(fl validator.FieldLevel) bool {
	if password, ok := fl.Field().Interface().(string); ok {
		if len(password) < 8 || len(password) > 64 {
			return false
		}

		if match, _ := regexp.MatchString("^[a-zA-Z0-9_!@#$%&*^]*$", password); !match {
			return false
		}
		if match, _ := regexp.MatchString("^.*[a-z].*$", password); !match {
			return false
		}
		if match, _ := regexp.MatchString("^.*[A-Z].*$", password); !match {
			return false
		}
		if match, _ := regexp.MatchString("^.*[0-9].*$", password); !match {
			return false
		}
		if match, _ := regexp.MatchString("^.*[_!@#$%&*^].*$", password); !match {
			return false
		}

		return true
	}
	return false
}
