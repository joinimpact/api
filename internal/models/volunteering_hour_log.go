package models

import (
	"context"
	"time"
)

// VolunteeringHourLog represents a log of a volunteer's verified hours.
type VolunteeringHourLog struct {
	Model
	OpportunityID  int64        `json:"opportunityId,omitempty"`
	Opportunity    Opportunity  `json:"-" gorm:"foreignkey:OpportunityID"`
	OrganizationID int64        `json:"organizationId,omitempty"`
	Organization   Organization `json:"-" gorm:"foreignkey:OrganizationID"`
	VolunteerID    int64        `json:"volunteerId"`
	Volunteer      User         `json:"-" gorm:"foreignkey:VolunteerID"`
	EventID        int64        `json:"eventID,omitempty"`
	Event          Event        `json:"-" gorm:"foreignkey:EventID"`
	GranterID      int64        `json:"granterId,omitempty"`
	Granter        User         `json:"-" gorm:"foreignkey:GranterID"`
	GrantedOn      time.Time    `json:"grantedOn"`
	GrantedHours   float32      `json:"grantedHours"`
}

// VolunteeringHourLogsResponse wraps an array of VolunteeringHourLogs and contains information from the database.
type VolunteeringHourLogsResponse struct {
	VolunteeringHourLogs []VolunteeringHourLog
	TotalResults         int
}

// VolunteeringHourLogRepository represents a repository of VolunteeringHourLog entities.
type VolunteeringHourLogRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(ctx context.Context, id int64) (*VolunteeringHourLog, error)
	// FindByOpportunityID finds multiple entities by the opportunity ID.
	FindByOpportunityID(ctx context.Context, opportunityID int64) (*VolunteeringHourLogsResponse, error)
	// FindByOrganizationID finds multiple entities by the organization ID.
	FindByOrganizationID(ctx context.Context, organizationID int64) (*VolunteeringHourLogsResponse, error)
	// FindByVolunteerID finds multiple entities by volunteer ID.
	FindByVolunteerID(ctx context.Context, volunteerID int64) (*VolunteeringHourLogsResponse, error)
	// Create creates a new entity.
	Create(ctx context.Context, volunteeringHourLog VolunteeringHourLog) error
	// Update updates an entity with the ID in the provided entity.
	Update(ctx context.Context, volunteeringHourLog VolunteeringHourLog) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(ctx context.Context, id int64) error
}
