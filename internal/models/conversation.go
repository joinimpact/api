package models

import "context"

// Conversation represents a single chat conversation in the database.
type Conversation struct {
	Model
	Active                       bool                          `json:"-"`
	Name                         string                        `json:"name" gorm:"-"`
	ProfilePicture               string                        `json:"profilePicture" gorm:"-"`
	OpportunityMembershipRequest *OpportunityMembershipRequest `json:"membershipRequest,omitempty" gorm:"-"`
	CreatorID                    int64                         `json:"creatorId"`
	OrganizationID               int64                         `json:"organizationID"`
	Organization                 Organization                  `json:"-"`
	Type                         int                           `json:"type"`
}

// ConversationsResponse represents a database response with multiple Conversations.
type ConversationsResponse struct {
	Conversations []Conversation
	TotalResults  int
}

// ConversationRepository represents a repository of conversation entities.
type ConversationRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*Conversation, error)
	// FindByIDs finds multiple entities by IDs.
	FindByIDs(ctx context.Context, ids []int64) (*ConversationsResponse, error)
	// FindByCreatorID finds multiple entities by the creator ID.
	FindByCreatorID(creatorID int64) ([]Conversation, error)
	// FindByOrganizationID finds multiple entities by the organization ID.
	FindByOrganizationID(ctx context.Context, organizationID int64) (*ConversationsResponse, error)
	// Create creates a new entity.
	Create(conversation Conversation) error
	// Update updates an entity with the ID in the provided entity.
	Update(conversation Conversation) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
