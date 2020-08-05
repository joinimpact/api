package models

import "context"

// VolunteeringHourLog represents a log of a volunteer's verified hours.
type VolunteeringHourLog struct {
	Model
	OpportunityID int64       `json:"opportunityId"`
	Opportunity   Opportunity `json:"-" gorm:"foreignkey:OpportunityID"`
	VolunteerID   int64       `json:"volunteerId"`
	Volunteer     User        `json:"-" gorm:"foreignkey:VolunteerID"`
	GranterID     int64       `json:"granterId"`
	Granter       User        `json:"-" gorm:"foreignkey:GranterID"`
	GrantedHours  float32     `json:"grantedHours"`
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
