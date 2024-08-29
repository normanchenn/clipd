package errors

import "errors"

var (
	ErrInvalidConfiguration  = errors.New("invalid configuration")
	ErrNotFound              = errors.New("not found")
	ErrOutOfBounds           = errors.New("out of bounds")
	ErrSerializationFailed   = errors.New("failed to serialize entry")
	ErrDeserializationFailed = errors.New("failed to deserialize entry")
	ErrWriteFailed           = errors.New("failed to write data")
)
