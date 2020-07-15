package events

import (
	"time"

	"github.com/joinimpact/api/pkg/location"
)

// EventBase represents all first-level event properties.
type EventBase struct {
	ID            int64  `json:"id"`
	OpportunityID int64  `json:"opportunityId"`
	CreatorID     int64  `json:"creatorId" scope:"manager"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Hours         int    `json:"hours"`
}

// EventView is the view in which event entities will be passed through the Service.
type EventView struct {
	EventBase
	EventSchedule *EventSchedule     `json:"schedule"`
	Location      *location.Location `json:"location"`
}

// EventSchedule represents a date/time range for a single event.
type EventSchedule struct {
	SingleDate bool      `json:"singleDate"`     // whether the event is a single date/time or a range of dates/times
	DateOnly   *bool     `json:"dateOnly"`       // when false, show a time as well
	FromDate   time.Time `json:"from,omitempty"` // the starting time, if applicable
	ToDate     time.Time `json:"to,omitempty"`   // the ending time, if applicable
}

// ModifyEventRequest represents the input to create/modify an event.
type ModifyEventRequest struct {
	EventBase
	EventSchedule *EventSchedule        `json:"schedule"`
	Location      *location.Coordinates `json:"location"`
}
