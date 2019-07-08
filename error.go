package cache

import "errors"

var (
	ErrNotFound      = errors.New("cache item not found")
	ErrZeroMaxValues = errors.New("max values cannot be zero")
)
