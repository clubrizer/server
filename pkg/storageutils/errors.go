package storageutils

import "errors"

var ErrNotFound = errors.New("db: row not found")
var ErrScanFailed = errors.New("db: failed to parse rows")
var ErrUnknown = errors.New("db: unknown error")

type ErrorCode int

const (
	Unknown ErrorCode = -1
)

type Error struct {
	Code ErrorCode
	Err  error
}

func (r *Error) Error() string {
	return r.Err.Error()
}
