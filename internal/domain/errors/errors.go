package errors

import (
	"fmt"
)

type ErrorType int

const (
	ValidationError ErrorType = iota
	NotFoundError
	InternalError
)

type BusinessError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *BusinessError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *BusinessError) Unwrap() error {
	return e.Err
}

func (e *BusinessError) IsValidation() bool {
	return e.Type == ValidationError
}

func (e *BusinessError) IsNotFound() bool {
	return e.Type == NotFoundError
}

func NewValidationError(message string) *BusinessError {
	return &BusinessError{
		Type:    ValidationError,
		Message: message,
	}
}

func NewNotFoundError(message string) *BusinessError {
	return &BusinessError{
		Type:    NotFoundError,
		Message: message,
	}
}

func NewInternalError(message string, err error) *BusinessError {
	return &BusinessError{
		Type:    InternalError,
		Message: message,
		Err:     err,
	}
}

var (
	ErrEmptyPlaylistName = NewValidationError("playlist name cannot be empty")
	ErrTooManyTracks     = NewValidationError("playlist cannot have more than allowed tracks")
	ErrInvalidPlaylistID = NewValidationError("invalid playlist ID")
	ErrEmptyTrackTitle   = NewValidationError("track title cannot be empty")
	ErrEmptyTrackArtist  = NewValidationError("track artist cannot be empty")
	ErrPlaylistNotFound  = NewNotFoundError("playlist not found")
)
