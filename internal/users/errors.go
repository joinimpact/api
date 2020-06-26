package users

// ErrUserNotFound is thrown when the server is unable to find a User.
type ErrUserNotFound struct {
}

// NewErrUserNotFound creates and returns a ErrUserNotFound.
func NewErrUserNotFound() error {
	return &ErrUserNotFound{}
}

// Error provides a string representation of the error.
func (e *ErrUserNotFound) Error() string {
	return "user not found"
}

// ErrTagNotFound is thrown when the server is unable to find a Tag.
type ErrTagNotFound struct {
}

// NewErrTagNotFound creates and returns a ErrTagNotFound.
func NewErrTagNotFound() error {
	return &ErrTagNotFound{}
}

// Error provides a string representation of the error.
func (e *ErrTagNotFound) Error() string {
	return "tag not found"
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
