package models

import "context"

// Event responses
const (
	EventResponseNull         = iota
	EventResponseCanAttend    = iota
	EventResponseCanNotAttend = iota
)

// EventResponse represents a volunteer's response to an event.
type EventResponse struct {
	Model
	UserID   int64 `json:"userId"`
	EventID  int64 `json:"eventId"`
	Response *int  `json:"response"`
}

// EventResponseRepository represents a repository of event responses.
type EventResponseRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(ctx context.Context, id int64) (*EventResponse, error)
	// FindByEventID finds multiple entities by the event ID.
	FindByEventID(ctx context.Context, eventID int64) ([]EventResponse, error)
	// FindByUserID finds multiple entities by the user's ID.
	FindByUserID(ctx context.Context, userID int64) ([]EventResponse, error)
	// FindInEventByUserID finds an entity by the entity and user's ID.
	FindInEventByUserID(ctx context.Context, eventID, userID int64) (*EventResponse, error)
	// Create creates a new entity.
	Create(ctx context.Context, eventResponse EventResponse) error
	// Update updates an entity with the ID in the provided entity.
	Update(ctx context.Context, eventResponse EventResponse) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(ctx context.Context, id int64) error
}
