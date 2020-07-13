package models

// ConversationOrganizationMembership represents an organization's relation to a conversation.
type ConversationOrganizationMembership struct {
	Model
	Active         bool  `json:"-"`
	ConversationID int64 `json:"conversationId"`
	OrganizationID int64 `json:"organizationId"`
	Role           int   `json:"role"`
}

// ConversationOrganizationMembershipRepository provides methods for interacting with conversation memberships.
type ConversationOrganizationMembershipRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*ConversationOrganizationMembership, error)
	// FindByConversationID finds multiple entities by the user ID.
	FindByConversationID(conversationID int64) ([]ConversationOrganizationMembership, error)
	// FindByOrganizationID finds multiple entities by the organization ID.
	FindByOrganizationID(organizationID int64) ([]ConversationOrganizationMembership, error)
	// Create creates a new entity.
	Create(conversationOrganizationMembership ConversationOrganizationMembership) error
	// Update updates an entity with the ID in the provided entity.
	Update(conversationOrganizationMembership ConversationOrganizationMembership) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
