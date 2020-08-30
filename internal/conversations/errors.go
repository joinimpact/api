package conversations

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

// ErrConversationNotFound is thrown when the server is unable to find a Conversation.
type ErrConversationNotFound struct {
}

// NewErrConversationNotFound creates and returns a ErrConversationNotFound.
func NewErrConversationNotFound() error {
	return &ErrConversationNotFound{}
}

// Error provides a string representation of the error.
func (e *ErrConversationNotFound) Error() string {
	return "conversation not found"
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

// ErrUserNotInConversation is thrown when a user does not have a membership in a particular conversation.
type ErrUserNotInConversation struct {
}

// NewErrUserNotInConversation creates and returns a ErrUserNotInConversation.
func NewErrUserNotInConversation() error {
	return &ErrUserNotInConversation{}
}

// Error provides a string representation of the error.
func (e *ErrUserNotInConversation) Error() string {
	return "user not found in conversation"
}

// Ref provides a string representation of the error.
func (e *ErrUserNotInConversation) Ref() string {
	return "conversations.user_not_in_conversation"
}
