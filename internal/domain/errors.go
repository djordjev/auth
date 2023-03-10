package domain

import (
	"fmt"
	"github.com/djordjev/auth/internal/domain/types"
)

type Error struct {
	cause      error
	message    string
	logMessage string
	isCritical bool
	user       types.User
}

func (e Error) Error() string {
	return e.message
}

func (e Error) LogError() string {
	var errorType string
	if e.isCritical {
		errorType = "critical error"
	} else {
		errorType = "regular error"
	}

	return fmt.Sprintf("domain %s for user %+v: %s", errorType, e.user, e.message)
}

func (e Error) Unwrap() error {
	return e.cause
}

func (e Error) IsCritical() bool {
	return e.isCritical
}

func NewDomainError(message string, logMessage string, cause error, isCritical bool, user types.User) Error {
	return Error{
		cause:      cause,
		message:    message,
		logMessage: logMessage,
		isCritical: isCritical,
		user:       user,
	}
}
