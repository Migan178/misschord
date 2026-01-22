package errors

import "fmt"

var (
	DuplicatedUniqueValueErr = fmt.Errorf("duplicated value")
)
