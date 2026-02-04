package errors

import "fmt"

var (
	ErrInternalServer = fmt.Errorf("an error occurred")
	ErrFailedToSend   = fmt.Errorf("failed to send")

	ErrorMessageChannelIDIsInvalid = "channelId is invalid"
	ErrorMessageAlreadyIdentified  = "already identified"
	ErrorMessageInvalidSyntax      = "invalid syntax"
	ErrorMessageInvalidToken       = "invalid token"
	ErrorMessageInternalDBError    = "internal database error"
	ErrorMessageConstraintErr      = "constraint violation"
	ErrorMessageUnauthorized       = "unauthorized session"
)

type ErrorCode int

const (
	ErrorCodeSyntaxError ErrorCode = 5000 + iota
	ErrorCodeAuthorizationError
	ErrorCodeNotfound
	ErrorCodeUnauthorized
	ErrorCodeInvalidValue
	ErrorCodeAuthenticationErr
)

const (
	ErrorCodeInternalError ErrorCode = 6000 + iota
)

type APIError struct {
	Code     ErrorCode         `json:"code"`
	Message  string            `json:"error,omitempty"`
	Messages map[string]string `json:"errors,omitempty"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api error %d: %s", e.Code, e.Message)
}
