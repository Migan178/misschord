package repository

type ErrorCode int

const (
	ErrorCodeOther ErrorCode = iota
	ErrorCodeNotFound
	ErrorCodeConstraint
	ErrorCodeAuthenticationFailed
)

type DatabaseError struct {
	RawErr error
	Code   ErrorCode
}

func (e *DatabaseError) Error() string {
	return e.RawErr.Error()
}
