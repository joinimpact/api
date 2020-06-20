package email

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Email represents a sendable email, essentially a draft.
type Email struct {
	Sender      *Sender
	Recipient   *Recipient
	Subject     string
	HTMLContent string
}

// NewEmail creates and returns a new Email with the given Sender.
func NewEmail(sender *Sender, recipient *Recipient, subject, htmlContent string) *Email {
	return &Email{
		sender,
		recipient,
		subject,
		htmlContent,
	}
}

// Send the email using the provided SendGrid client.
func (e *Email) Send(client *sendgrid.Client) error {
	// Define the from and to.
	from := mail.NewEmail(e.Sender.Name, e.Sender.Email)
	to := mail.NewEmail(e.Recipient.Name, e.Recipient.Email)
	message := mail.NewSingleEmail(from, e.Subject, to, "", e.HTMLContent)

	// Send the message.
	if _, err := client.Send(message); err != nil {
		return err
	}

	return nil
}
