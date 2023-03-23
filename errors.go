package diskcache

import "errors"

var (
	ErrBadDir  = errors.New("invalid directory")
	ErrBadSize = errors.New("cache size must be greater then zero")

	ErrNotFound = errors.New("not found")

	ErrTooLarge = errors.New("file size must be less or equal cache size")
)
