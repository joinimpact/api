package models

// ConversationOpportunityMembershipRequest represents a connection of a conversation to a single opportunity membership request.
type ConversationOpportunityMembershipRequest struct {
	Model
	ConversationID                 int64
	OpportunityMembershipRequestID int64
}

// ConversationOpportunityMembershipRequestRepository provides methods for interacting with conversation opportunity memberships.
type ConversationOpportunityMembershipRequestRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*ConversationOpportunityMembershipRequest, error)
	// FindByConversationID finds multiple entities by the conversation ID.
	FindByConversationID(conversationID int64) ([]ConversationOpportunityMembershipRequest, error)
	// Create creates a new entity.
	Create(conversationOpportunityMembershipRequest ConversationOpportunityMembershipRequest) error
	// Update updates an entity with the ID in the provided entity.
	Update(conversationOpportunityMembershipRequest ConversationOpportunityMembershipRequest) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
