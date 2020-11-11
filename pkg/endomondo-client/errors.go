package endomondo

import "errors"

var (
	// ErrBadCredentials throws when provided user is not found
	ErrBadCredentials = errors.New("bad credentials")
	// ErrUnauthorized throws when trying to reach any endpoint without authorization
	ErrUnauthorized = errors.New("unauthorized")
	// ErrUnexpected something bad happend, but fk it
	ErrUnexpected = errors.New("unexpected error")
)
