package models

// ConversationMembership represents a user's relation to a conversation.
type ConversationMembership struct {
	Model
	Active         bool  `json:"-"`
	ConversationID int64 `json:"conversationId"`
	UserID         int64 `json:"userId"`
	Role           int   `json:"role"`
}

// ConversationMembershipRepository provides methods for interacting with conversation memberships.
type ConversationMembershipRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*ConversationMembership, error)
	// FindByConversationID finds multiple entities by the conversation ID.
	FindByConversationID(conversationID int64) ([]ConversationMembership, error)
	// FindByUserID finds multiple entities by the user ID.
	FindByUserID(userID int64) ([]ConversationMembership, error)
	// Create creates a new entity.
	Create(conversationMembership ConversationMembership) error
	// Update updates an entity with the ID in the provided entity.
	Update(conversationMembership ConversationMembership) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
