// Package storageutils provides helpers for interacting with data stores.
package storageutils

import "errors"

// ErrNotFound states that no rows could be found for a given query.
var ErrNotFound = errors.New("db: row not found")

// ErrScanFailed states that the program failed to parse the given rows into structs.
var ErrScanFailed = errors.New("db: failed to parse rows")

// ErrUnknown states that an unknown database error has happened.
var ErrUnknown = errors.New("db: unknown error")
