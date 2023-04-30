package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func MsgForTag(tag string) string {

	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "regexp":
		return "Invalid Telephone"
	}

	return ""
}

func Regexp(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(fl.Param())
	return re.MatchString(fl.Field().String())
}
