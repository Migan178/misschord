package errors

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func GetErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "min":
		return fmt.Sprintf("the value must be at least %s character(s)", fe.Param())
	case "max":
		return fmt.Sprintf("the value must be %s character(s) or less", fe.Param())
	case "eqfield":
		return "the value in not equal"
	case "email":
		return "the email format is not valid"
	default:
		return "the value is not valid"
	}
}
