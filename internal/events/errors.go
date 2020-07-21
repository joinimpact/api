package events

// ErrEventNotFound is thrown when the server is unable to find an Event.
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

// Ref provides a representation of the error.
func (e *ErrEventNotFound) Ref() string {
	return "events.event_not_found"
}

// ErrResponseNotFound is thrown when the server is unable to find an EventResponse.
type ErrResponseNotFound struct {
}

// NewErrResponseNotFound creates and returns a ErrResponseNotFound.
func NewErrResponseNotFound() error {
	return &ErrResponseNotFound{}
}

// Error provides a string representation of the error.
func (e *ErrResponseNotFound) Error() string {
	return "event response not found"
}

// Ref provides a representation of the error.
func (e *ErrResponseNotFound) Ref() string {
	return "events.response_not_found"
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

// Ref provides a representation of the error.
func (e *ErrServerError) Ref() string {
	return "generic.server_error"
}
