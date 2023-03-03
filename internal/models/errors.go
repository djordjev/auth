package models

import "fmt"

type ModelError struct {
	dbError error
	message string
}

func (m ModelError) Error() string {
	return fmt.Sprintf("model error: %s", m.message)
}

func (m ModelError) Unwrap() error {
	return m.dbError
}

func NewModelError(msg string, wrapping error) ModelError {
	return ModelError{
		dbError: wrapping,
		message: msg,
	}
}
