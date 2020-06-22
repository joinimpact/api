package email

// Recipient represents a recipient of an email.
type Recipient struct {
	Name  string
	Email string
}

// NewRecipient creates and returns a new Recipient with the provided name and email.
func NewRecipient(name, email string) *Recipient {
	return &Recipient{
		name,
		email,
	}
}
