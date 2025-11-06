package errors

import "errors"

var (
	// Common errors
	ErrInvalidJSON        = errors.New("invalid JSON format")
	ErrMissingParameters  = errors.New("missing required parameters")
	ErrInvalidQueryParams = errors.New("invalid query parameters")
	ErrUnsupportedMedia   = errors.New("unsupported media type")

	// Event errors
	ErrEventNotFound   = errors.New("event not found")
	ErrInvalidDate     = errors.New("invalid date format, expected YYYY-MM-DD")
	ErrEventConflict   = errors.New("event with this ID already exists")
	ErrEmptyEventID    = errors.New("event ID cannot be empty")
	ErrEmptyUserID     = errors.New("user ID cannot be empty")
	ErrEmptyTitle      = errors.New("event title cannot be empty")
	ErrEventValidation = errors.New("event validation failed")
	ErrEmptyParameters = errors.New("required parameters are empty")
)
