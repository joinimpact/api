package models

// Conversation represents a single chat conversation in the database.
type Conversation struct {
	Model
	Active         bool  `json:"-"`
	CreatorID      int64 `json:"creatorId"`
	OrganizationID int64 `json:"organizationID"`
	Organization   Organization
	Type           int `json:"type"`
}

// ConversationRepository represents a repository of conversation entities.
type ConversationRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*Conversation, error)
	// FindByCreatorID finds multiple entities by the creator ID.
	FindByCreatorID(creatorID int64) ([]Conversation, error)
	// FindByOrganizationID finds multiple entities by the organization ID.
	FindByOrganizationID(organizationID int64) ([]Conversation, error)
	// Create creates a new entity.
	Create(conversation Conversation) error
	// Update updates an entity with the ID in the provided entity.
	Update(conversation Conversation) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
