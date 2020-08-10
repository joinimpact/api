package hours

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

// ErrRequestNotFound is thrown when a volunteering hour log request is not found.
type ErrRequestNotFound struct {
}

// NewErrRequestNotFound creates and returns a ErrRequestNotFound.
func NewErrRequestNotFound() error {
	return &ErrRequestNotFound{}
}

// Error provides a string representation of the error.
func (e *ErrRequestNotFound) Error() string {
	return "hour request not found"
}

// Ref provides a representation of the error.
func (e *ErrRequestNotFound) Ref() string {
	return "hours.request_not_found"
}
