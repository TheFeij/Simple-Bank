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

		match, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
		return match
	}
	return false
}
