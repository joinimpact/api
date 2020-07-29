package events

import (
	"time"

	"github.com/joinimpact/api/pkg/location"
)

// EventBase represents all first-level event properties.
type EventBase struct {
	ID             int64  `json:"id"`
	OpportunityID  int64  `json:"opportunityId"`
	CreatorID      int64  `json:"creatorId" scope:"manager"`
	Title          string `json:"title" validate:"min=4,max=128"`
	Description    string `json:"description"`
	Hours          int    `json:"hours" validate:"min=0,max=500"`
	HoursFrequency *uint  `json:"hoursFrequency" validate:"min=0,max=1"`
	TotalHours     int    `json:"totalHours"`
}

// EventView is the view in which event entities will be passed through the Service.
type EventView struct {
	EventBase             `validate:"dive"`
	EventSchedule         *EventSchedule         `json:"schedule"`
	EventResponsesSummary *EventResponsesSummary `json:"responses,omitempty" scope:"manager"`
	Location              *location.Location     `json:"location"`
}

// EventSchedule represents a date/time range for a single event.
type EventSchedule struct {
	SingleDate bool      `json:"singleDate"`                         // whether the event is a single date/time or a range of dates/times
	DateOnly   *bool     `json:"dateOnly"`                           // when false, show a time as well
	FromDate   time.Time `json:"from,omitempty" validate:"required"` // the starting time, if applicable
	ToDate     time.Time `json:"to,omitempty"`                       // the ending time, if applicable
}

// ModifyEventRequest represents the input to create/modify an event.
type ModifyEventRequest struct {
	EventBase
	EventSchedule *EventSchedule        `json:"schedule" validate:"dive"`
	Location      *location.Coordinates `json:"location" validate:"omitempty,dive"`
}

// EventResponsesSummary contains the number of volunteers joining, not joining, and total number of volunteers of an event.
type EventResponsesSummary struct {
	NumCanAttend    uint `json:"numCanAttend"`
	NumCanNotAttend uint `json:"numCanNotAttend"`
	TotalMembers    uint `json:"totalMembers"`
}
