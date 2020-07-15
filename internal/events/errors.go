package events

// ErrEventNotFound is thrown when the server is unable to find a Event.
type ErrEventNotFound struct {
}

// NewErrEventNotFound creates and returns a ErrEventNotFound.
func NewErrEventNotFound() error {
	return &ErrEventNotFound{}
}

// Error provides a string representation of the error.
func (e *ErrEventNotFound) Error() string {
	return "event not found"
}

// ErrServerError is thrown when the server experiences an internal error.
type ErrServerError struct {
}

// NewErrServerError creates and returns a ErrServerError.
func NewErrServerError() error {
	return &ErrServerError{}
}

// Error provides a string representation of the error.
func (e *ErrServerError) Error() string {
	return "internal error processing request, please try again"
}
