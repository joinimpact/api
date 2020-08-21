package organizations

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

// ErrOrganizationNotFound is thrown when the server is unable to find a User.
type ErrOrganizationNotFound struct {
}

// NewErrOrganizationNotFound creates and returns a ErrOrganizationNotFound.
func NewErrOrganizationNotFound() error {
	return &ErrOrganizationNotFound{}
}

// Error provides a string representation of the error.
func (e *ErrOrganizationNotFound) Error() string {
	return "organization not found"
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

// ErrUserAlreadyInOrganization is thrown when a user invite is requested, but the user is already a member of the organization.
type ErrUserAlreadyInOrganization struct {
}

// NewErrUserAlreadyInOrganization creates and returns a ErrUserAlreadyInOrganization.
func NewErrUserAlreadyInOrganization() error {
	return &ErrUserAlreadyInOrganization{}
}

// Error provides a string representation of the error.
func (e *ErrUserAlreadyInOrganization) Error() string {
	return "user already a member of organization"
}

// Ref provides a string reference representation of the error.
func (e *ErrUserAlreadyInOrganization) Ref() string {
	return "organizations.user_already_in_organization"
}

// ErrUserAlreadyInvited is thrown when a user has already been invited to an opportunity.
type ErrUserAlreadyInvited struct {
}

// NewErrUserAlreadyInvited creates and returns a ErrUserAlreadyInvited.
func NewErrUserAlreadyInvited() error {
	return &ErrUserAlreadyInvited{}
}

// Error provides a string representation of the error.
func (e *ErrUserAlreadyInvited) Error() string {
	return "user already invited"
}

// Ref provides a representation of the error.
func (e *ErrUserAlreadyInvited) Ref() string {
	return "organizations.user_already_invited"
}

// ErrInviteInvalid is thrown when an invite is invalid in any way.
type ErrInviteInvalid struct {
}

// NewErrInviteInvalid creates and returns a ErrInviteInvalid.
func NewErrInviteInvalid() error {
	return &ErrInviteInvalid{}
}

// Error provides a string representation of the error.
func (e *ErrInviteInvalid) Error() string {
	return "invite invalid"
}

// Ref provides a representation of the error.
func (e *ErrInviteInvalid) Ref() string {
	return "organizations.invite_invalid"
}
