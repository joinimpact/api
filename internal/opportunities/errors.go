package opportunities

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

// ErrOpportunityNotFound is thrown when an opportunity is not found.
type ErrOpportunityNotFound struct {
}

// NewErrOpportunityNotFound creates and returns a ErrOpportunityNotFound.
func NewErrOpportunityNotFound() error {
	return &ErrOpportunityNotFound{}
}

// Error provides a string representation of the error.
func (e *ErrOpportunityNotFound) Error() string {
	return "opportunity not found"
}

// ErrOpportunityNotPublishable is thrown when an opportunity is missing required fields and is unable to be published.
type ErrOpportunityNotPublishable struct {
	InvalidFields []string
}

// NewErrOpportunityNotPublishable creates and returns a ErrOpportunityNotPublishable.
func NewErrOpportunityNotPublishable(invalidFields []string) error {
	return &ErrOpportunityNotPublishable{invalidFields}
}

// Error provides a string representation of the error.
func (e *ErrOpportunityNotPublishable) Error() string {
	return "opportunity not found"
}

// ErrMembershipAlreadyRequested is thrown when a user has already requested to join an opportunity.
type ErrMembershipAlreadyRequested struct {
}

// NewErrMembershipAlreadyRequested creates and returns a ErrMembershipAlreadyRequested.
func NewErrMembershipAlreadyRequested() error {
	return &ErrMembershipAlreadyRequested{}
}

// Error provides a string representation of the error.
func (e *ErrMembershipAlreadyRequested) Error() string {
	return "opportunity membership already requested"
}
