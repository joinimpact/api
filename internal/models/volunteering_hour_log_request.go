package models

import "context"

// VolunteeringHourLogRequest represents a volunteer's request from an organization to log hours.
type VolunteeringHourLogRequest struct {
	Model
	Accepted       bool         `json:"accepted"`
	OpportunityID  int64        `json:"opportunityId"`
	Opportunity    Opportunity  `json:"-" gorm:"foreignkey:OpportunityID"`
	OrganizationID int64        `json:"organizationId"`
	Organization   Organization `json:"-" gorm:"foreignkey:OrganizationID"`
	EventID        int64        `json:"eventID"`
	Event          Event        `json:"-" gorm:"foreignkey:EventID"`
	VolunteerID    int64        `json:"volunteerId"`
	Volunteer      User         `json:"-" gorm:"foreignkey:VolunteerID"`
	RequestedHours float32      `json:"grantedHours"`
}

// VolunteeringHourLogRequestsResponse wraps an array of VolunteeringHourLogRequests and contains information from the database.
type VolunteeringHourLogRequestsResponse struct {
	VolunteeringHourLogRequests []VolunteeringHourLogRequest
	TotalResults                int
}

// VolunteeringHourLogRequestRepository represents a repository of VolunteeringHourLogRequest entities.
type VolunteeringHourLogRequestRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(ctx context.Context, id int64) (*VolunteeringHourLogRequest, error)
	// FindByOpportunityID finds multiple entities by the opportunity ID.
	FindByOpportunityID(ctx context.Context, opportunityID int64) (*VolunteeringHourLogRequestsResponse, error)
	// FindByOrganizationID finds multiple entities by the organization ID.
	FindByOrganizationID(ctx context.Context, organizationID int64) (*VolunteeringHourLogRequestsResponse, error)
	// FindByVolunteerID finds multiple entities by volunteer ID.
	FindByVolunteerID(ctx context.Context, volunteerID int64) (*VolunteeringHourLogRequestsResponse, error)
	// Create creates a new entity.
	Create(ctx context.Context, volunteeringHourLog VolunteeringHourLogRequest) error
	// Update updates an entity with the ID in the provided entity.
	Update(ctx context.Context, volunteeringHourLog VolunteeringHourLogRequest) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(ctx context.Context, id int64) error
}
