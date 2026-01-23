package errors

import "fmt"

var (
	ErrDuplicatedUniqueValue = fmt.Errorf("duplicated value")
	ErrNoUser                = fmt.Errorf("no user")
	ErrInternalServer        = fmt.Errorf("an error occurred")
)
