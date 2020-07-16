package models

import (
	"context"
	"time"
)

// OpportunityPermissions flags
const (
	OpportunityPermissionsMember  = iota // opportunity member/volunteer
	OpportunityPermissionsAdmin   = iota
	OpportunityPermissionsCreator = iota
)

// OpportunityMembership creates a relationship between Oppourtunities and their volunteers.
type OpportunityMembership struct {
	Model
	Active          bool      `json:"-"` // controls whether or not the entity is active
	UserID          int64     `json:"-"` // the ID of the user being granted membership
	JoinedAt        time.Time `json:"joinedAt"`
	OpportunityID   int64     `json:"-"`               // the ID of the opportunity the user is being granted access to
	PermissionsFlag int       `json:"permissionsFlag"` // a flag which designates permissions the user has
	InviterID       int64     `json:"inviterId"`       // the ID of the user who invited the member to the opportunity
}

// OpportunityMembershipRepository represents a repository of opportunity memberships.
type OpportunityMembershipRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(ctx context.Context, id int64) (*OpportunityMembership, error)
	// FindByUserID finds multiple entities by the user ID.
	FindByUserID(ctx context.Context, userID int64) ([]OpportunityMembership, error)
	// FindByOpportunityID finds multiple entities by the opportunity ID.
	FindByOpportunityID(ctx context.Context, opportunityID int64) ([]OpportunityMembership, error)
	// FindUserInOpportunity finds a user's membership in a specific opportunity.
	FindUserInOpportunity(ctx context.Context, opportunityID, userID int64) (*OpportunityMembership, error)
	// Create creates a new entity.
	Create(ctx context.Context, opportunityMembership OpportunityMembership) error
	// Update updates an entity with the ID in the provided entity.
	Update(ctx context.Context, opportunityMembership OpportunityMembership) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(ctx context.Context, id int64) error
}
