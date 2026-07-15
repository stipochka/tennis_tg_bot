package validation

import "errors"

var (
	ErrCantBook          = errors.New("too late to book this time")
	ErrInvalidDate       = errors.New("invalid date provided")
	ErrLimitReached      = errors.New("booking hours limit exceeded")
	ErrInvalidTimeRange  = errors.New("court is closed at this time")
	ErrInvalidTimeFormat = errors.New("invalid time format")
)
