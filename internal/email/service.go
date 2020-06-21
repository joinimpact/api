package email

import (
	"github.com/joinimpact/api/internal/config"
	"github.com/sendgrid/sendgrid-go"
)

// Service is a service for sending emails to users.
type Service interface {
	NewEmail(recipient *Recipient, subject, htmlContent string) *Email
	Send(e *Email) error
}

// service is the internal representation of the Service.
type service struct {
	config *config.Config
	sender *Sender
	client *sendgrid.Client
}

// NewService creates and returns a new Service with the given config and
// sender identity.
func NewService(config *config.Config, sender *Sender) Service {
	client := sendgrid.NewSendClient(config.SendGridAPIKey)

	return &service{
		config,
		sender,
		client,
	}
}

// NewEmail creates a new Email using the sender
func (s *service) NewEmail(recipient *Recipient, subject, htmlContent string) *Email {
	return NewEmail(s.sender, recipient, subject, htmlContent)
}

// Send sends a provided Email.
func (s *service) Send(e *Email) error {
	return e.Send(s.client)
}
