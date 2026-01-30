package errors

import "fmt"

var (
	ErrDuplicatedUniqueValue       = fmt.Errorf("duplicated value")
	ErrNoUser                      = fmt.Errorf("no user")
	ErrNoChannel                   = fmt.Errorf("no channel")
	ErrInternalServer              = fmt.Errorf("an error occurred")
	ErrorMessageChannelIDIsInvalid = "channelId is invalid"
	ErrFailedToSend                = fmt.Errorf("failed to send")
	ErrorMessageAlreadyIdentified  = "already identified"
	ErrorMessageInvalidSyntax      = "invalid syntax"
	ErrorMessageInvalidToken       = "invalid token"
	ErrorMessageInternalDBError    = "internal database error"
	ErrorMessageConstraintErr      = "constraint violation"
	ErrInvalidData                 = fmt.Errorf("invalid data")
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
