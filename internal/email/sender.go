package email

// Sender represents a sender for an email.
type Sender struct {
	Name  string
	Email string
}

// NewSender creates and returns a new Sender with the provided name and email.
func NewSender(name, email string) *Sender {
	return &Sender{
		name,
		email,
	}
}
