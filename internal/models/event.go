package models

import (
	"context"
	"time"
)

// HoursFrequency constants
const (
	EventHoursFrequencyOnce   uint = iota
	EventHoursFrequencyPerDay uint = iota
)

// Event represents a scheduled event under an opportunity.
type Event struct {
	Model
	Active            bool      `json:"-"`
	OpportunityID     int64     `json:"opportunityId"`
	CreatorID         int64     `json:"creatorId"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Hours             int       `json:"hours"`
	HoursFrequency    *uint     `json:"hoursFrequency"`
	DateOnly          *bool     `json:"dateOnly"`
	FromDate          time.Time `json:"from"`
	ToDate            time.Time `json:"to"`
	LocationLatitude  float64   `json:"-"` // the latitude of the events's location
	LocationLongitude float64   `json:"-"` // the longitude of the events's location
}

// EventRepository represents a repository of events.
type EventRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(ctx context.Context, id int64) (*Event, error)
	// FindByOpportunityID finds multiple entities by the opportunity ID.
	FindByOpportunityID(ctx context.Context, opportunityID int64) ([]Event, error)
	// FindByOpportunityIDs finds multiple entities by multiple opportunity IDs.
	FindByOpportunityIDs(ctx context.Context, opportunityIDs []int64) ([]Event, error)
	// FindByCreatorID finds multiple entities by the creator ID.
	FindByCreatorID(ctx context.Context, creatorID int64) ([]Event, error)
	// Create creates a new entity.
	Create(ctx context.Context, event Event) error
	// Update updates an entity with the ID in the provided entity.
	Update(ctx context.Context, event Event) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(ctx context.Context, id int64) error
}
