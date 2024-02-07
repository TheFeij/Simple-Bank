package requests

import "github.com/go-playground/validator/v10"

var ValidOwner validator.Func = func(fl validator.FieldLevel) bool {
	if owner, ok := fl.Field().Interface().(string); ok {
		return len(owner) >= 1 && len(owner) <= 50
	}
	return false
}
