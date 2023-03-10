package api

type Error struct {
	cause   error
	message string
}

func (e Error) Error() string {
	return e.message
}

func (e Error) Unwrap() error {
	return e.cause
}

func NewApiError(message string, cause error, code int) Error {
	return Error{
		cause:   cause,
		message: message,
	}
}
