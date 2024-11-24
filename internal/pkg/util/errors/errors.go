package errors

import (
	err "github.com/pkg/errors"
)

var (
	ErrUnknown       = err.New("unknown error")
	ErrInvalidParams = err.New("invalid params")
	ErrFailedOS      = err.New("failed to execute os command")
)
