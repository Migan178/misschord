package errors

import (
	"encoding/json"
	"errors"
	"fmt"
)

func GetJSONTypeErrorMessage(field, expectedType, realType string) string {
	return fmt.Sprintf("invalid field %s type: excepted type: %s, real type: %s", field, expectedType, realType)
}

func GetJSONTypeIsNullErrorMessage(field string) string {
	return GetJSONTypeErrorMessage(field, "any (not nullable)", "null")
}

func GetUnmarshalError(err error) error {
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		return &APIError{
			Code:    ErrorCodeSyntaxError,
			Message: GetJSONTypeErrorMessage(typeErr.Field, typeErr.Type.String(), typeErr.Value),
		}
	}

	if errors.As(err, new(*json.SyntaxError)) {
		return &APIError{
			Code:    ErrorCodeSyntaxError,
			Message: ErrorMessageInvalidSyntax,
		}
	}

	return err
}
