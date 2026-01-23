package errors

import "fmt"

var (
	ErrDuplicatedUniqueValue = fmt.Errorf("duplicated value")
	ErrInternalServer        = fmt.Errorf("an error occurred")
)
