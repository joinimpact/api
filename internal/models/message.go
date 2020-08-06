package models

import (
	"context"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

// Message types.
const (
	MessageTypeStandard                   = "MESSAGE_STANDARD"
	MessageTypeVolunteerRequestProfile    = "MESSAGE_VOLUNTEER_REQUEST_PROFILE"
	MessageTypeVolunteerRequestAcceptance = "MESSAGE_VOLUNTEER_REQUEST_ACCEPTANCE"
	MessageTypeEventCreated               = "MESSAGE_EVENT_CREATED"
	MessageTypeHoursRequested             = "MESSAGE_HOURS_REQUESTED"
	MessageTypeHoursAccepted              = "MESSAGE_HOURS_ACCEPTED"
)

// Sender perspectives.
const (
	MessageSenderPerspectiveVolunteer    uint = iota
	MessageSenderPerspectiveOrganization uint = iota
)

// Message represents a single message in a conversation.
type Message struct {
	Model
	Timestamp         time.Time      `json:"timestamp"`
	ConversationID    int64          `json:"conversationId"`
	SenderID          int64          `json:"senderId"`
	SenderPerspective *uint          `json:"senderPerspective"`
	Type              string         `json:"type"`
	Body              postgres.Jsonb `json:"body"`
	Edited            bool           `json:"edited"`
	EditedTimestamp   time.Time      `json:"editedTimestamp,omitempty"`
}

// MessagesResponse represents a response from the MessageRepository with multiple messages.
type MessagesResponse struct {
	Messages     []Message
	TotalResults int
}

// MessageRepository represents a repository of message entities.
type MessageRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(ctx context.Context, id int64) (*Message, error)
	// FindByConversationID finds multiple entities by the conversation ID.
	FindByConversationID(ctx context.Context, conversationID int64) (*MessagesResponse, error)
	// FindBySenderID finds multiple entities by the sender ID.
	FindBySenderID(ctx context.Context, senderID int64) (*MessagesResponse, error)
	// FindInConversationBySenderID finds multiple entities by the sender and conversation ID.
	FindInConversationBySenderID(ctx context.Context, conversationID, senderID int64) (*MessagesResponse, error)
	// Create creates a new entity.
	Create(ctx context.Context, message Message) error
	// Update updates an entity with the ID in the provided entity.
	Update(ctx context.Context, message Message) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(ctx context.Context, id int64) error
}
