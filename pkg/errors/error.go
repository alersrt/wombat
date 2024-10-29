package errors

type Error struct {
	Message string
	Cause   *error
}

func (receiver *Error) Error() string {
	return receiver.Message
}

func NewWithCauseError(message string, cause *error) *Error {
	return &Error{
		Message: message,
		Cause:   cause,
	}
}

func NewError(message string) *Error {
	return NewWithCauseError(message, nil)
}
